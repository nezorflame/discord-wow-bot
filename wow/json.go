package wow

import (
    "encoding/json"
)

func (r *Realms) getRealmsFromJSON(body *[]byte) error {
    err := json.Unmarshal(*body, r)
    panicOnErr(err)
    return err
}

func getGuildInfoFromJSON(body *[]byte) (*GuildInfo, error) {
    var gi = new(GuildInfo)
    err := json.Unmarshal(*body, gi)
    panicOnErr(err)
    return gi, err
}

func getCharacterFromJSON(body *[]byte) (*Character, error) {
    var c = new(Character)
    err := json.Unmarshal(*body, c)
    panicOnErr(err)
    return c, err
}

func (i *Item) getItemFromJSON(body *[]byte) error {
    err := json.Unmarshal(*body, i)
    panicOnErr(err)
    return err
}

func getJSONFromGuildInfo(gi *GuildInfo) ([]byte, error) {
    body, err := json.Marshal(&gi)
    panicOnErr(err)
    return body, err
}

func getURLFromJSON(body *[]byte) (*string, error) {
    var apiResponse = new(googlAPIResponse)
    err := json.Unmarshal(*body, apiResponse)
    panicOnErr(err)
    return &apiResponse.ID, err
}
