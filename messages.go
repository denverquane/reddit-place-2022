package main

type ConnectionInitMessage struct {
	Type    string                       `json:"type"`
	Payload ConnectionInitMessagePayload `json:"payload"`
}

type ConnectionInitMessagePayload struct {
	Authorization string
}

type StartMessage struct {
	ID      string              `json:"id"`
	Type    string              `json:"type"`
	Payload StartMessagePayload `json:"payload"`
}

type StartMessagePayload struct {
	Extensions    struct{}                     `json:"extensions"`
	OperationName string                       `json:"operationName"`
	Query         string                       `json:"query"`
	Variables     StartMessagePayloadVariables `json:"variables"`
}

type StartMessagePayloadVariables struct {
	Input StartMessagePayloadVariablesInput `json:"input"`
}

type StartMessagePayloadVariablesInput struct {
	Channel StartMessagePayloadVariablesInputChannel `json:"channel"`
}

type StartMessagePayloadVariablesInputChannel struct {
	TeamOwner string `json:"teamOwner"`
	Category  string `json:"category"`
	Tag       string `json:"tag"`
}
