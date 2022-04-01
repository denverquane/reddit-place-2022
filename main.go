package main

import (
	"fmt"
	"github.com/denverquane/reddit-place-2022/pkg"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	addr   = "gql-realtime-2.reddit.com"
	origin = "https://hot-potato.reddit.com"
)

func main() {
	var err error
	token := os.Getenv("REDDIT_BEARER_TOKEN")
	if token == "" {
		log.Fatal("Please supply your REDDIT_BEARER_TOKEN in env")
	} else {
		log.Println("Loaded REDDIT_BEARER_TOKEN from env")
	}
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.Replace(token, "Bearer ", "", 1)
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

	// start the worker to process messages as we receive them over the websocket
	go pkg.PlaceWorker(c, ready, done)

	// TODO properly jsonify
	msg := fmt.Sprintf("{\"type\":\"connection_init\",\"payload\":{\"Authorization\":\"Bearer %s\"}}", token)
	err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write error: ", err)
		return
	}

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
			return
		}
	}
}
