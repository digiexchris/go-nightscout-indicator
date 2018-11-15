package main

import (
	"encoding/json"
	"fmt"
	"github.com/digiexchris/go-nightscout-indicator/configuration"
	"io/ioutil"
	"net/http"
)

type values struct {
	Sgv float32
	Delta float32
}

func refresh() (value float32, delta float32, err error) {
	url := fmt.Sprintf("https://%s/api/v1/entries/current.json",configuration.App.NightscoutHost)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("API-SECRET", configuration.App.ApiSecret)

	response, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	} else {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, 0, err
		}

		var v []values
		err = json.Unmarshal(data,&v)
		if err != nil {
			return 0, 0, err
		}

		return v[0].Sgv, v[0].Delta, nil
	}

	return 0,0, nil
}
