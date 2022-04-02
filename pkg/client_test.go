package pkg

import (
	"log"
	"testing"
)

func TestGetRedditAuthToken(t *testing.T) {
	token, expireTime, expiresIn, err := GetRedditAuthToken("https://new.reddit.com/r/place/")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(token)
	log.Println("Expires: ", expireTime)
	log.Println("Expires in: ", expiresIn, "seconds")
}
