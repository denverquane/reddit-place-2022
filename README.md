# reddit-place-2022
A Go download and processing tool for Reddit's 2022 r/place April Fools data.

Makes use of Devin Smith's incredible Parquet-packaged dataset, more details here:
https://deephaven.io/blog/2022/04/08/place-csv-to-parquet/

This eliminates the need for over 22Gb of CSV data (or a ~90Gb Postgres DB).

# Installation
1. [Install Go 1.18](https://go.dev/dl/)
2. Clone this repo: `git clone https://github.com/denverquane/reddit-place-2022`
4. Build or run the code: `go build main.go` (builds executable) or `go run main.go`

# Usage
Run `go run main.go` or 

`go build main.go` followed by `./main` or `./main.exe`

The application will download the parquet file if it doesn't exist, and, at the time of writing,
iterate over all the events and generate a number of output png images to `images/`

These include:
```
// the final image of place (includes white-out pixel events)
place.png 

// snapshots for every 5% of total pixels drawn to place (todo this should be time-based)
place_<percent>.png 

// snapshots of regions before they were blanked-out by moderator intervention
place_mod_[<start>]-[<end>].png
```