package pkg

import (
	"log"
	"testing"
)

func TestGetRedditAuthToken(t *testing.T) {
	token, err := GetRedditAuthToken("https://new.reddit.com/r/place/")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(token)
}
