package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

// Get - function for getting the GET request's response in form of JSON
func Get(url string) (result []byte, err error) {
	for i := 0; i < o.APIMaxRetries; i++ {
		RateLimiter.Wait(1)

		if result, err = getJSONResponse(url); err == nil {
			break
		}
	}

	return
}

func getJSONResponse(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Unable to create GET request: %s", err)
	}
	req.Header.Set("Connection", "close")

	client := newClient()

	r, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer r.Body.Close()

	if r.StatusCode > 400 {
		return []byte{}, fmt.Errorf("GET request error - %s", r.Status)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Unable to read GET response: %s", err)
	}

	return body, nil
}

// PostJSONResponse - function for getting the POST request response in form of JSON
// value transmitted is a link for a Google URL Shortener
func PostJSONResponse(url, value string) ([]byte, error) {
	var jsonStr = []byte(`{"longUrl": "` + value + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return []byte{}, fmt.Errorf("Unable to create GET request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	client := newClient()

	r, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer r.Body.Close()

	if r.StatusCode == 404 {
		return []byte{}, fmt.Errorf(r.Status)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

// GetShortLink returns a short link to the long one from goo.gl
func GetShortLink(longLink string) (shortLink string, err error) {
	var respJSON []byte
	gURL := fmt.Sprintf(o.GoogleShortenerLink, o.GoogleToken)

	if respJSON, err = PostJSONResponse(gURL, longLink); err != nil {
		err = errors.Wrap(err, "Unable to post JSON response to Google")
		return
	}

	if shortLink, err = GetURLFromJSON(respJSON); err != nil {
		err = errors.Wrap(err, "Unable to get URL from JSON")
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

func newClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   0,
			KeepAlive: 0,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return &http.Client{Transport: transport}
}
