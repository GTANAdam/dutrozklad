package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Get(url string) io.ReadCloser {
	// defer TimeTrack(time.Now(), "GET")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil
	}

	if resp.StatusCode != 200 {
		log.Printf("Status code error: %d, response: %s\n", resp.StatusCode, resp.Status)
		return nil
	}

	return resp.Body
}

// Post ...
func Post(url string, values url.Values) (io.ReadCloser, error) {
	// defer TimeTrack(time.Now(), "POST")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.PostForm(url, values)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status code error: %d, response: %s", resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}
