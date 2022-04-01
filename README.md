# reddit-place-2022
Testing interaction with Reddit's r/place for 2022

Currently this tool is solely for processing/proxying the current PNG canvas of reddit/r/place in purely a read-only fashion.

# Help Wanted
- Submitting or engaging with the Websocket/place in a more substantive way, such as for marking pixels automatically

# Installation
Install Go 1.18, clone this repo, and run `go build main.go` (builds executable) or `go run main.go`

# Usage
`go run main.go` or 

`go build main.go` followed by `./main` or `./main.exe`

There is also a Docker image provided at `denverquane/reddit-place-2022`

You can then visit `localhost:8080` to view the current place image