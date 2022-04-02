package pkg

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

// kinda hacky, but it works
var searchRegex = regexp.MustCompile("\"session\":{\"accessToken\":\"(?P<token>.+)\",\"expires\":\"(?P<expireTime>.+)\",\"expiresIn\":(?P<expiresIn>[0-9]+)")
var (
	oauthConfig      *oauth2.Config
	oauthStateString = "sdfkjeli239890knfvao8e"
)

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

func GetRedditOAuthToken() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	oauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("REDDIT_CLIENT_ID"),
		ClientSecret: os.Getenv("REDDIT_CLIENT_SECRET"),
		Scopes:       []string{"identity"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.reddit.com/api/v1/authorize",
			TokenURL: "https://www.reddit.com/api/v1/access_token",
		},
	}

	authUrl := oauthConfig.AuthCodeURL(oauthStateString)
	browser.OpenURL(authUrl)
}
