package main

import (
	"encoding/json"
)

func (gi *GuildInfo) unmarshal(body *[]byte) error {
	return json.Unmarshal(*body, gi)
}

func (gi *GuildInfo) marshal() ([]byte, error) {
	return json.Marshal(&gi)
}

func (r *Realms) unmarshal(body *[]byte) error {
	return json.Unmarshal(*body, r)
}

func (c *Character) unmarshal(body *[]byte) error {
	return json.Unmarshal(*body, c)
}

func (i *Item) unmarshal(body *[]byte) error {
	return json.Unmarshal(*body, i)
}

func getURLFromJSON(body *[]byte) (apiResponseID string, err error) {
	apiResponse := new(googlAPIResponse)
	if err = json.Unmarshal(*body, apiResponse); err != nil {
		apiResponseID = apiResponse.ID
	}
	return
}
