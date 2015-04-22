package bitly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	api = "https://api-ssl.bitly.com/v3/shorten"
)

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

func (c *Client) Shorten(longUrl string) (shortUrl string, err error) {
	client := &http.Client{}
	endpoint := fmt.Sprintf("%s?access_token=%s&longUrl=%s", api, c.Token, longUrl)
	fmt.Printf("GET %s", endpoint)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}

	defer res.Body.Close()

	resp, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("error: %s", string(resp))
	}

	var v map[string]interface{}
	json.Unmarshal(resp, &v)

	data := v["data"].(map[string]interface{})
	return data["url"].(string), nil
}
