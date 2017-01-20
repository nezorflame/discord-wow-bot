package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// GetJSONResponse - function for getting the GET request response in form of JSON
func GetJSONResponse(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
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
	return body, nil
}

// PostJSONResponse - function for getting the POST request response in form of JSON
// value transmitted is a link for a Google URL Shortener
func PostJSONResponse(url, value string) ([]byte, error) {
	var jsonStr = []byte(`{"longUrl": "` + value + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if strings.Contains(r.Status, "404") {
		return nil, errors.New(r.Status)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
