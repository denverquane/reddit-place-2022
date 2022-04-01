package main

import (
	"context"
	"fmt"
	"github.com/denverquane/reddit-place-2022/pkg"
	"github.com/gorilla/websocket"
	"image"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	baseUrl = "https://new.reddit.com/r/place/"
	addr    = "gql-realtime-2.reddit.com"
	origin  = "https://hot-potato.reddit.com"
)

var GlobalImage image.Image
var GlobalImageLock sync.RWMutex

func main() {
	token, err := pkg.GetRedditAuthToken(baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: addr, Path: "/query"}
	log.Printf("connecting to %s", u.String())

	headers := http.Header{}
	headers.Add("Sec-WebSocket-Protocol", "graphql-ws")
	headers.Add("Origin", origin)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("dial: ", err)
	}
	defer c.Close()

	done := make(chan struct{})
	ready := make(chan struct{})

	taskQueue := make(chan pkg.DownloadTask)
	// start the worker to process messages as we receive them over the websocket
	go pkg.WebsocketWorker(c, ready, done, taskQueue)

	go func() {
		for {
			select {
			case <-done:
				return

			case task := <-taskQueue:
				if task.ImageType == pkg.Base {
					GlobalImageLock.Lock()
					GlobalImage, err = pkg.DownloadImage(task.URL)
					GlobalImageLock.Unlock()

					if err != nil {
						log.Println(err)
					} else {
						log.Println("Downloaded base image", task.URL)
					}
				} else {
					startTime := time.Now()
					diffImage, err := pkg.DownloadImage(task.URL)
					if err != nil {
						log.Println(err)
					} else {

						GlobalImageLock.Lock()
						GlobalImage, err = pkg.CombineDiffToBase(GlobalImage, diffImage)
						GlobalImageLock.Unlock()

						if err != nil {
							log.Println(err)
						} else {
							log.Println("Downloaded and combined diff img in", time.Now().Sub(startTime).String())
						}
					}
				}
			}
		}
	}()

	// TODO properly jsonify
	msg := fmt.Sprintf("{\"type\":\"connection_init\",\"payload\":{\"Authorization\":\"Bearer %s\"}}", token)
	err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write error: ", err)
		return
	}

	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", pkg.HandleRequestWrapper(&GlobalImage, &GlobalImageLock))
	go func() {
		srv.ListenAndServe()
	}()

	for {
		select {
		case <-done:
			return

		case <-ready:
			// TODO obv this shouldn't be an explicit string, should be a struct that gets jsonified
			err = c.WriteMessage(websocket.TextMessage, []byte("{\"id\":\"2\",\"type\":\"start\",\"payload\":{\"variables\":{\"input\":{\"channel\":{\"teamOwner\":\"AFD2022\",\"category\":\"CANVAS\",\"tag\":\"0\"}}},\"extensions\":{},\"operationName\":\"replace\",\"query\":\"subscription replace($input: SubscribeInput!) {\\n  subscribe(input: $input) {\\n    id\\n    ... on BasicMessage {\\n      data {\\n        __typename\\n        ... on FullFrameMessageData {\\n          __typename\\n          name\\n          timestamp\\n        }\\n        ... on DiffFrameMessageData {\\n          __typename\\n          name\\n          currentTimestamp\\n          previousTimestamp\\n        }\\n      }\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\"}}{\"type\":\"ka\"}"))
			if err != nil {
				log.Println(err)
			}
			log.Println("Sent Ready message to Reddit")
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			srv.Shutdown(context.TODO())
			return
		}
	}
}
