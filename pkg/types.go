package pkg

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
