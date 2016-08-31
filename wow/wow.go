package wow

import (
    "fmt"
    "errors"
    "strings"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
)

var logger *log.Logger
var wowAPIToken string

// Realm - type for WoW server realm info
type Realm struct {
    Type            string      `json:"type"`
    Population      string      `json:"population"`
    Queue           bool        `json:"queue"`
    Status          bool        `json:"status"`
    Name            string      `json:"name"`
    Slug            string      `json:"slug"`
    Battlegroup     string      `json:"battlegroup"`
    Locale          string      `json:"locale"`
    Timezone        string      `json:"timezone"`
    ConnectedRealms []string    `json:"connected_realms"`
}

// Realms is a slice of all WoW realms
// Specializing on EU with locale "ru_RU" 
var realms []Realm

// RealmsAPIResponse is a struct for slice of Realm
type RealmsAPIResponse struct {
    RealmList []Realm `json:"realms"`
}

func logDebug(v ...interface{}) {
	logger.SetPrefix("DEBUG ")
	logger.Println(v...)
}

func logInfo(v ...interface{}) {
	logger.SetPrefix("INFO  ")
	logger.Println(v...)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

// InitializeWoWAPI - function for initializing WoW API
func InitializeWoWAPI(token *string) {
    logger = log.New(os.Stderr, "  ", log.Ldate|log.Ltime)
    wowAPIToken = *token
}

func getWoWRealms() {
    r, err := http.Get("https://eu.api.battle.net/wow/realm/status?locale=ru_RU&apikey=" + wowAPIToken)
    panicOnErr(err)
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    realmsResponse, err := getRealmsAPIResponse([]byte(body))
    panicOnErr(err)
    realms = realmsResponse.RealmList
}

func getRealmsAPIResponse(body []byte) (*RealmsAPIResponse, error) {
    var s = new(RealmsAPIResponse)
    err := json.Unmarshal(body, &s)
    panicOnErr(err)
    return s, err
}

// GetWoWRealmStatus - function for receiving realm status
func GetWoWRealmStatus(realmName string) (bool, error) {
    getWoWRealms()
    for _, r := range realms {
        if r.Name == realmName || r.Slug == realmName {
            return r.Status, nil
        }
    }
    return false, errors.New("No such realm is present!")
}

// GetWoWRealmQueueStatus - function for receiving realm queue status
func GetWoWRealmQueueStatus(realmName string) (bool, error) {
    getWoWRealms()
    for _, r := range realms {
        if r.Name == realmName || r.Slug == realmName {
            return r.Queue, nil
        }
    }
    return false, errors.New("No such realm is present!")
}

// GetWoWRealmInfo - function for receiving realm info
func GetWoWRealmInfo(realmName string) (string, error) {
    realm, err := getRealmByName(realmName)
    if err != nil {
        return "", err
    }

    realmInfo := "Имя сервера: %v\n"
    realmInfo += "Тип сервера: %v\n"
    realmInfo += "Населенность: %v\n"
    realmInfo += "Статус: %t\n"
    realmInfo += "Очередь на вход: %t\n"
    realmInfo += "PvP-группа: %v\n"
    realmInfo += "Язык: %v\n"
    realmInfo += "Временной пояс: %v\n"
    realmInfo += "Связанные серверы: %v"
    realmInfo = fmt.Sprintf(realmInfo, realm.Name, realm.Type, realm.Population, realm.Status,
        realm.Queue, realm.Battlegroup, realm.Locale, realm.Timezone, realm.ConnectedRealms)
    return realmInfo, nil
}

func getRealmByName(realmName string) (Realm, error) {
    getWoWRealms()
    logInfo("getRealmByName: " + realmName)
    var realm Realm
    for _, r := range realms {
        if strings.ToLower(r.Name) == strings.ToLower(realmName) || 
           strings.ToLower(r.Slug) == strings.ToLower(realmName) {
            return r, nil
        }
    }
    return realm, errors.New("No such realm is present!")
}

// GetRealmName returns realm name string
func GetRealmName(message string, command string) string {
    realmString := strings.Replace(message, command, "", 1)
    if realmString == "" {
        realmString = "Ревущий фьорд"
    } else {
        realmString = strings.TrimLeft(realmString, " ")
    }
    return realmString
}
