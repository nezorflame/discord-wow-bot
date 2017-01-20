package main

import (
    "encoding/json"
)

func (r *Realms) getRealmsFromJSON(body *[]byte) error {
    err := json.Unmarshal(*body, r)
    panicOnErr(err)
    return err
}

func (gi *GuildInfo) unmarshal(body *[]byte) (error) {
    err := json.Unmarshal(*body, gi)
    panicOnErr(err)
    return err
}

func (gi *GuildInfo) marshal() ([]byte, error) {
    body, err := json.Marshal(&gi)
    panicOnErr(err)
    return body, err
}

func getCharacterFromJSON(body *[]byte) (*Character, error) {
    c := new(Character)
    err := json.Unmarshal(*body, c)
    panicOnErr(err)
    return c, err
}

func (i *Item) getItemFromJSON(body *[]byte) error {
    err := json.Unmarshal(*body, i)
    panicOnErr(err)
    return err
}

func getURLFromJSON(body *[]byte) (*string, error) {
    apiResponse := new(googlAPIResponse)
    err := json.Unmarshal(*body, apiResponse)
    panicOnErr(err)
    return &apiResponse.ID, err
}
