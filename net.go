package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
)

// GetJSONResponse - synced function for getting the GET request response in form of JSON
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
	time.Sleep(10 * time.Millisecond)
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

// GetShortLink returns a short link to the long one from goo.gl
func GetShortLink(longLink string) (shortLink string, err error) {
	var respJSON []byte

	gURL := fmt.Sprintf(o.GoogleShortenerLink, o.GoogleToken)

	if respJSON, err = PostJSONResponse(gURL, longLink); err != nil {
		glog.Errorf("Unable to post JSON response to Google: %s", err)
		return
	}

	if shortLink, err = GetURLFromJSON(respJSON); err != nil {
		glog.Errorf("Unable to get URL from JSON: %s", err)
	}

	return
}

// DownloadFile gets the file from the URL and saves to the set path
func DownloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
