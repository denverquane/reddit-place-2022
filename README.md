# reddit-place-2022
A Go download and processing tool for Reddit's 2022 r/place April Fools data.

https://www.reddit.com/r/place/comments/txvk2d/rplace_datasets_april_fools_2022/

# Installation
1. [Install Go 1.18](https://go.dev/dl/)
2. Clone this repo: `git clone https://github.com/denverquane/reddit-place-2022`
3. Startup a Postgres DB. For example, `docker run -e POSTGRES_PASSWORD=<something_secure> -p 5432:5432 -d postgres`
4. Build or run the code: `go build main.go` (builds executable) or `go run main.go`

# Usage
First set the following environment variables relevant to your Postgres DB from Installation step 3:
```
POSTGRES_URL= # for example, localhost:5432
POSTGRES_USER=postgres
POSTGRES_PASS=<something_secure>
```

Then run `go run main.go` or 

`go build main.go` followed by `./main` or `./main.exe`

Downloads and unpacks all csv files into the current directory (TODO: config download dir)

At the time of writing, this utility downloads and unpacks all the data files needed for further processing, and inserts
the events into a Postgres Database. See [postgres.sql](internal/postgres.sql) for more information about the schema