# reddit-place-2022
A Go download and processing tool for Reddit's 2022 r/place April Fools data.

https://www.reddit.com/r/place/comments/txvk2d/rplace_datasets_april_fools_2022/

# Installation
1. [Install Go 1.18](https://go.dev/dl/)
2. Clone this repo: `git clone https://github.com/denverquane/reddit-place-2022`
3. Build or run the code: `go build main.go` (builds executable) or `go run main.go`

# Usage
`go run main.go` or 

`go build main.go` followed by `./main` or `./main.exe`

Downloads and unpacks all csv files into the current directory (TODO: config download dir)

At the time of writing, this utility simply downloads and unpacks all the data files needed for further processing.
The data files have considerable overlap in timestamps, so using an external DB or storage mechanism is probably the
best course of action, as opposed to manually sorting all the events