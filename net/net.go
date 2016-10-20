package net

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetJSONResponse(url string) ([]byte, error) {
	r, err := http.Get(url)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	if strings.Contains(r.Status, "404") {
		return nil, errors.New(r.Status)
	}
	if strings.Contains(r.Status, "403") {
		time.Sleep(1 * time.Second)
		return GetJSONResponse(url)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return []byte(body), nil
}

func PostJSONResponse(url, value string) ([]byte, error) {
	var jsonStr = []byte(`{"longUrl": "` + value + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(req)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	if strings.Contains(r.Status, "404") {
		return nil, errors.New(r.Status)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return []byte(body), nil
}
