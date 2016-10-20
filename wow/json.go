package wow

import "encoding/json"

func getRealmsFromJSON(body []byte) (*Realms, error) {
    var r = new(Realms)
    err := json.Unmarshal(body, &r)
    panicOnErr(err)
    return r, err
}

func getGuildInfoFromJSON(body []byte) (*GuildInfo, error) {
    var gi = new(GuildInfo)
    err := json.Unmarshal(body, &gi)
    panicOnErr(err)
    return gi, err
}

func getJSONFromGuildInfo(gi *GuildInfo) ([]byte, error) {
    body, err := json.Marshal(&gi)
    panicOnErr(err)
    return body, err
}

func getCharacterFromJSON(body []byte) (*Character, error) {
    var c = new(Character)
    err := json.Unmarshal(body, &c)
    panicOnErr(err)
    return c, err
}

func getItemFromJSON(body []byte) (*Item, error) {
    var i = new(Item)
    err := json.Unmarshal(body, &i)
    panicOnErr(err)
    return i, err
}

func getURLFromJSON(body []byte) (*string, error) {
    var apiResponse = new(googlAPIResponse)
    err := json.Unmarshal(body, &apiResponse)
    panicOnErr(err)
    return &apiResponse.ID, err
}