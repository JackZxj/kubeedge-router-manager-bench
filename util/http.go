package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpHandler func(url string, data []byte) (string, error)

func Post(url string, data []byte) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result), nil
}
