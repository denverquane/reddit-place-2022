package pkg

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

var searchRegex = regexp.MustCompile("\"session\":{\"accessToken\":\"(?P<token>.+)\",\"expires\":\"(?P<expireTime>.+)\",\"expiresIn\":(?P<expiresIn>[0-9]+)")
var (
	oauthConfig      *oauth2.Config
	oauthStateString = "sdfkjeli239890knfvao8e"
)

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
