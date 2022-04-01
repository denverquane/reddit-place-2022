package main

import (
	"encoding/json"
	"fmt"
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

type Ack struct {
	Type string `json:"type"`
}

type MessageData struct {
	Payload PayloadData `json:"payload"`
	ID      string      `json:"id"`
	Type    string      `json:"type"`
}

type PayloadData struct {
	Data SubscribeData `json:"data"`
}

type SubscribeData struct {
	Subscribe Data `json:"subscribe"`
}

type Data struct {
	ID       string  `json:"id"`
	Data     RawData `json:"data"`
	TypeName string  `json:"__typename"`
}

type RawData struct {
	TypeName  string  `json:"__typename"`
	Name      string  `json:"name"`
	Timestamp float64 `json:"timestamp"`
}

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

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read: ", err)
				return
			}
			var ackMsg Ack
			err = json.Unmarshal(message, &ackMsg)
			if err == nil {
				if ackMsg.Type == "connection_ack" {
					log.Println("Received connection ack")
					// send a message saying we're ready to receive data
					ready <- struct{}{}
				} else if ackMsg.Type == "data" {
					var dataMsg MessageData
					err = json.Unmarshal(message, &dataMsg)
					if err == nil {
						rawData := dataMsg.Payload.Data.Subscribe.Data
						if rawData.TypeName == "FullFrameMessageData" {
							log.Println("Full Frame Message URL: ", rawData.Name)
						} else if rawData.TypeName == "DiffFrameMessageData" {
							log.Println("Diff Frame Message URL: ", rawData.Name)
						}
					} else {
						log.Println(err)
						log.Printf("recv: %s", message)
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
