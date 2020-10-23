// Package query ..
package query

import (
	"dutrozkladbot/config"
	"dutrozkladbot/model"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func get(url string) (io.ReadCloser, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if config.Bot.Debug {
		fmt.Println(url)
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %v, response: %s", resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}

func process(response *io.ReadCloser) *model.Response {
	body, err := ioutil.ReadAll(*response)
	if err != nil {
		log.Panicln(err)
		return nil
	}

	defer (*response).Close()

	var result model.Response
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(err)
	}

	return &result
}

func query(str string) (*model.Response, error) {
	response, err := get(str)
	if err != nil {
		return nil, err
	}

	return process(&response), nil
}
