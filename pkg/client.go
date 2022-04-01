package pkg

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

var searchRegex = regexp.MustCompile("\"session\":{\"accessToken\":\"(?P<token>.+)\",\"expires\":\"(?P<expireTime>.+)\",\"expiresIn\":(?P<expiresIn>[0-9]+)")

func GetRedditAuthToken(baseUrl string) (string, error) {
	resp, err := http.Get(baseUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if match := searchRegex.FindStringSubmatch(string(body)); match != nil {
		token := match[searchRegex.SubexpIndex("token")]
		return token, nil
		//expireTime := match[searchRegex.SubexpIndex("expireTime")]
		//expiresIn := match[searchRegex.SubexpIndex("expiresIn")]
	}
	return "", errors.New("Response did not contain a regex match for the token")
}
