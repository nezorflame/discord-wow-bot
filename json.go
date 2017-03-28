package main

import (
	"encoding/json"
)

// Unmarshal makes the GuildInfo from the []byte
func (gi *GuildInfo) Unmarshal(body []byte) error {
	return json.Unmarshal(body, gi)
}

// Unmarshal makes the Realms from the []byte
func (r *Realms) Unmarshal(body []byte) error {
	return json.Unmarshal(body, r)
}

// Unmarshal makes the Character from the []byte
func (c *Character) Unmarshal(body []byte) error {
	return json.Unmarshal(body, c)
}

// Unmarshal makes the Item from the []byte
func (i *Item) Unmarshal(body []byte) error {
	return json.Unmarshal(body, i)
}

// GetURLFromJSON returns short goo.gl link
func GetURLFromJSON(body []byte) (apiResponseID string, err error) {
	apiResponse := new(URLShortenerAPIResponse)
	if err = json.Unmarshal(body, apiResponse); err == nil {
		apiResponseID = apiResponse.ID
	}
	return
}
