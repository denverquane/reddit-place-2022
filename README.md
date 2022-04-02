# reddit-place-2022
Testing interaction with Reddit's r/place for 2022

Currently this tool is solely for processing/proxying the current PNG canvas of reddit/r/place in purely a read-only fashion.

# Help Wanted
- Submitting or engaging with the Websocket/place in a more substantive way, such as for marking pixels automatically

Here is the main .js file used to reconstruct the GraphQL query that subscribes for incoming canvas data:
https://www.redditstatic.com/mona-lisa/en-US/index-3cc1ba23.js
Search for "setPixel" to see info for how to mark a pixel, but note, this will require a valid Reddit login,
not the anonymous Bearer token login that is currently in place

# Installation
Install Go 1.18, clone this repo, and run `go build main.go` (builds executable) or `go run main.go`

Create a `.env` file with the following values, which can be created here: https://www.reddit.com/prefs/apps
```bash
REDDIT_CLIENT_ID=
REDDIT_CLIENT_SECRET=
```

# Usage
`go run main.go` or 

`go build main.go` followed by `./main` or `./main.exe`

There is also a Docker image provided at `denverquane/reddit-place-2022`

You can then visit `localhost:8080` to view the current place image