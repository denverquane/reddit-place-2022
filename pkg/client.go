package pkg

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// kinda hacky, but it works
var searchRegex = regexp.MustCompile("\"session\":{\"accessToken\":\"(?P<token>.+)\",\"expires\":\"(?P<expireTime>.+)\",\"expiresIn\":(?P<expiresIn>[0-9]+)")

// GetRedditAuthToken returns an anonymous login Bearer Token for reddit, alongside the expiration time (unparsed),
// and the number of seconds before the token expires. Alongside any errors that make have occurred
func GetRedditAuthToken(baseUrl string) (string, string, uint64, error) {
	resp, err := http.Get(baseUrl)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", 0, err
	}

	if match := searchRegex.FindStringSubmatch(string(body)); match != nil {
		token := match[searchRegex.SubexpIndex("token")]
		expireTime := match[searchRegex.SubexpIndex("expireTime")]
		expiresIn := match[searchRegex.SubexpIndex("expiresIn")] // milliseconds
		i, err := strconv.ParseUint(expiresIn, 10, 64)
		if err != nil {
			log.Println(err)
		} else {
			i /= 1000 // seconds
		}
		return token, expireTime, i, nil
	}
	return "", "", 0, errors.New("Response did not contain a regex match for the token")
}
