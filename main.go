package main

import (
	"context"
	"encoding/json"
	"image"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/denverquane/reddit-place-2022/pkg"
	"github.com/gorilla/websocket"
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

	msg, err := json.Marshal(ConnectionInitMessage{
		Type: "connection_init",
		Payload: ConnectionInitMessagePayload{
			Authorization: "Bearer " + token,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

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
			msg, err := json.Marshal(StartMessage{
				ID:   "1",
				Type: "start",
				Payload: StartMessagePayload{
					Extensions:    struct{}{},
					OperationName: "replace",
					Query: `subscription replace($input: SubscribeInput!) {
	subscribe(input: $input) {
		id
		... on BasicMessage {
			data {
				__typename
				... on FullFrameMessageData {
					__typename
					name
					timestamp
				}
				... on DiffFrameMessageData {
					__typename
					name
					currentTimestamp
					previousTimestamp
				}
			}
			__typename
		}
		__typename
	}
}`,
					Variables: StartMessagePayloadVariables{
						Input: StartMessagePayloadVariablesInput{
							Channel: StartMessagePayloadVariablesInputChannel{
								TeamOwner: "AFD2022",
								Category:  "CANVAS",
								Tag:       "0",
							},
						},
					},
				},
			})

			if err != nil {
				log.Println(err)
			}

			err = c.WriteMessage(websocket.TextMessage, msg)
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
