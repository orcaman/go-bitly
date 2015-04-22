package bitly

import (
	"flag"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	token = flag.String("apiKey", "", "your bit.ly api key")
)

func TestBitly(t *testing.T) {
	flag.Parse()

	if *token == "" {
		tokenEnv := os.Getenv("BITLY_AUTH")
		if len(tokenEnv) == 0 {
			t.Error("missing required flag: token")
		} else {
			*token = tokenEnv
		}
	}

	bitlyClient := NewClient(*token)

	var testUrl = "https://github.com/streamrail/go-bitly"

	if shortUrl, err := bitlyClient.Shorten(testUrl); err != nil {
		t.Error(err.Error())
	} else {
		c := make(chan bool, 1)
		client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.String() == testUrl {
				c <- true
				return nil
			}
			for _, v := range via {
				if v.URL.String() == testUrl {
					c <- true
				}
			}

			return nil
		}}
		req, err := http.NewRequest("GET", string(shortUrl), nil)
		if err != nil {
			t.Error(err.Error())
		}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err.Error())
		}
		defer res.Body.Close()
		select {
		case <-c:
		case <-time.After(10 * time.Second):
			t.Errorf("error\n")
		}
	}
}
