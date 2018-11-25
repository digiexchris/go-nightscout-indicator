package nightscoutclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Reading struct {
	SGV        float32
	Delta      float32
	Error      error
	Direction  string
	Date       time.Time
	OldReading bool
}

type nightscoutTime struct {
	time.Time
}

//func (obj nightscoutTime) MarshalJSON() ([]byte, error) {
//	seconds := time.Time(obj).Unix()
//	return []byte(strconv.FormatInt(seconds, 10)), nil
//}

func (ct *nightscoutTime) UnmarshalJSON(b []byte) (err error) {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}

	ct.Time = time.Unix(i/1000, 1)
	return
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
		Date      nightscoutTime
	}

	url := fmt.Sprintf("https://%s/api/v1/entries/current.json", host)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("API-SECRET", secret)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return Reading{
			Error: err,
		}
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		var v []values
		err = json.Unmarshal(data, &v)
		if err != nil {
			log.Println(err)
			return Reading{
				Error: err,
			}
		}

		return Reading{
			SGV:       v[0].Sgv,
			Delta:     v[0].Delta,
			Direction: v[0].Direction,
			Date:      v[0].Date.Time,
			Error:     nil,
		}
	}
}
