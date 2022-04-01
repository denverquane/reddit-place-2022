package pkg

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type ImgType int

const (
	Base ImgType = iota
	Diff
)

type DownloadTask struct {
	ImageType ImgType
	URL       string
}

func WebsocketWorker(c *websocket.Conn, ready chan<- struct{}, done chan struct{}, taskQueue chan<- DownloadTask) {
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
						taskQueue <- DownloadTask{
							ImageType: Base,
							URL:       rawData.Name,
						}
					} else if rawData.TypeName == "DiffFrameMessageData" {
						taskQueue <- DownloadTask{
							ImageType: Diff,
							URL:       rawData.Name,
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
