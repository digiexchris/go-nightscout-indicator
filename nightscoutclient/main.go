package nightscoutclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Reading struct {
	SGV       float32
	Delta     float32
	Error     error
	Direction string
}

type NightscoutClient interface {
	Get(host string, secret string) Reading
}

type Client struct {
	HttpClient HttpClient
}

//go:generate counterfeiter . HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
	client http.Client
}

func New() NightscoutClient {

	timeout := time.Duration(5 * time.Second)
	c := http.Client{
		Timeout: timeout,
	}
	httpClient := httpClient{
		client: c,
	}

	return &Client{
		HttpClient: &httpClient,
	}
}

func (hc *httpClient) Do(req *http.Request) (*http.Response, error) {
	return hc.client.Do(req)
}

func (c *Client) Get(host string, secret string) Reading {
	type values struct {
		Sgv       float32
		Delta     float32
		Direction string
	}

	url := fmt.Sprintf("https://%s/api/v1/entries/current.json", host)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("API-SECRET", secret)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return Reading{
			Error: err,
		}
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		var v []values
		err = json.Unmarshal(data, &v)
		if err != nil {
			return Reading{
				Error: err,
			}
		}

		return Reading{
			SGV:       v[0].Sgv,
			Delta:     v[0].Delta,
			Direction: v[0].Direction,
			Error:     nil,
		}
	}
}
