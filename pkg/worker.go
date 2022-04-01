package pkg

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

func PlaceWorker(c *websocket.Conn, ready chan<- struct{}, done chan struct{}) {
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
						err = downloadFile(rawData.Name, "place.png")
						if err != nil {
							log.Println(err)
						} else {
							log.Println("Successfully downloaded place.png")
						}
					} else if rawData.TypeName == "DiffFrameMessageData" {
						err = downloadFile(rawData.Name, "diff.png")
						if err != nil {
							log.Println(err)
						} else {
							log.Println("Successfully downloaded diff.png")

							// TODO process diff and combine with place
						}
					}
				} else {
					log.Println(err)
					log.Printf("recv: %s", message)
				}
			}
		}
	}
}
