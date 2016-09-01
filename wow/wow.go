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

var (
    logger          *log.Logger
    apiToken        string
    locale          string
    region          string
)

const (
    apiRealmsLink       = "https://%v.api.battle.net/wow/realm/status?locale=%v&apikey=%v"
    apiGuildMembersLink = "https://%v.api.battle.net/wow/guild/%v/%v?fields=members&locale=%v&apikey=%v"
    apiCharacterItems   = "https://%v.api.battle.net/wow/character/%v/%v?fields=items&locale=%v&apikey=%v"
)

// Realm - type for WoW server realm info
type Realm struct {
    Type                string          `json:"type"`
    Population          string          `json:"population"`
    Queue               bool            `json:"queue"`
    Status              bool            `json:"status"`
    Name                string          `json:"name"`
    Slug                string          `json:"slug"`
    Battlegroup         string          `json:"battlegroup"`
    Locale              string          `json:"locale"`
    Timezone            string          `json:"timezone"`
    ConnectedRealms     []string        `json:"connected_realms"`
}

// Realms - struct for a slice of Realm
type Realms struct {
    RealmList           []Realm        `json:"realms"`
}

// GuildInfo - struct for WoW guild information
type GuildInfo struct {
    Name                string          `json:"name"`
    Realm               string          `json:"realm"`
    BattleGroup         string          `json:"battlegroup"`
    Level               int             `json:"level"`
    SideInt             int             `json:"side"`
    Side                string
    AchievementPoints   int             `json:"achievementPoints"`
    LastModified        int             `json:"lastModified"`
    GuildMembersList    []GuildMember   `json:"members"`
}

// GuildMember - struct for a WoW guild member
type GuildMember struct {
    Member              Character       `json:"character"`
    Rank                int             `json:"rank"`
}

// Character - struct for a WoW character
type Character struct {
    Name                string          `json:"name"` 
    Realm               string          `json:"realm"`
    FactionInt          int             `json:"faction"`
    Faction             string
    BattleGroup         string          `json:"battlegroup"` 
    ClassInt            int             `json:"class"`
    Class               string
    RaceInt             int             `json:"race"`
    Race                string 
    GenderInt           int             `json:"gender"`
    Gender              string
    Level               int             `json:"level"` 
    AchievementPoints   int             `json:"achievementPoints"` 
    Thumbnail           string          `json:"thumbnail"` 
    Spec                Specialization  `json:"spec"` 
    Guild               string          `json:"guild"` 
    GuildRealm          string          `json:"guildRealm"` 
    LastModified        int             `json:"lastModified"`
    Items               Items           `json:"items"`
}

// Specialization - struct for a WoW character specialization
type Specialization struct {
    Name                string          `json:"name"`
    Role                string          `json:"role"`
    BackgroundImage     string          `json:"backgroundImage"`
    Icon                string          `json:"icon"`
    Description         string          `json:"description"`
    Order               int             `json:"order"`
}

// Items - struct for storing items info for a character
type Items struct {
    AvgItemLvl          int             `json:"averageItemLevel"`
    AvgItemLvlEq        int             `json:"averageItemLevelEquipped"`
}

var (
    classes             map[int]string
    factions            map[int]string
    races               map[int]string
    genders             map[int]string
)

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
    apiToken = *token
    // TODO: Rework to config
    locale = "ru_RU"
    region = "eu"
    // TODO: Rework to WoW API
    classes = map[int]string{
        1  : "Воин",
        2  : "Паладин",
        3  : "Охотник",
        4  : "Разбойник",
        5  : "Жрец",
        6  : "Рыцарь смерти",
        7  : "Шаман",
        8  : "Маг",
        9  : "Чернокнижник",
        10 : "Монах",
        11 : "Друид",
        12 : "Охотник на демонов",
    }
    genders = map[int]string{
        0  : "Мужчина",
        1  : "Женщина",
    }
    factions   = map[int]string{
        0  : "Альянс",
        1  : "Орда",
    }
    races   = map[int]string{
        1  : "Человек",
        2  : "Орк",
        3  : "Дворф",
        4  : "Ночной эльф",
        5  : "Нежить",
        6  : "Таурен",
        7  : "Гном",
        8  : "Тролль",
        9  : "Гоблин",
        10 : "Эльф крови",
        11 : "Дреней",
        22 : "Ворген",
        24 : "Пандарен",
        25 : "Пандарен",
        26 : "Пандарен",
    } 
}

// GetRealmStatus - function for receiving realm status
func GetRealmStatus(realmName string) (bool, error) {
    realms, err := getRealms()
    if err != nil {
        return false, err
    }
    for _, r := range *realms {
        if r.Name == realmName || r.Slug == realmName {
            return r.Status, nil
        }
    }
    return false, errors.New("No such realm is present!")
}

// GetRealmQueueStatus - function for receiving realm queue status
func GetRealmQueueStatus(realmName string) (bool, error) {
    realms, err := getRealms()
    if err != nil {
        return false, err
    }
    for _, r := range *realms {
        if r.Name == realmName || r.Slug == realmName {
            return r.Queue, nil
        }
    }
    return false, errors.New("No such realm is present!")
}

// GetRealmInfo - function for receiving realm info
func GetRealmInfo(realmName string) (string, error) {
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

// GetGuildMembers - function for receiving a list of guild members
func GetGuildMembers(realmName string, guildName string) (string, error) {
    guildMembersList := "__**Имя         Уровень, раса, класс и текущий спек    Уровень вещей**__\n\n"
    characterStringPreset := "**%v**    __%d %v %v (%v)__    **%d**\n"

    gMembers, err := getGuildMembers(&realmName, &guildName)
    if err != nil {
        return "", err
    }

    for _, m := range *gMembers {
        specName := m.Member.Spec.Name
        if specName == "" {
            specName = "Нет инфы"
        }
        charString := fmt.Sprintf(characterStringPreset,
                                  m.Member.Name,
                                  m.Member.Level,
                                  m.Member.Race,
                                  m.Member.Class,
                                  specName,
                                  m.Member.Items.AvgItemLvlEq)
        guildMembersList += charString
    }
    return guildMembersList, nil
}

// GetRealmName returns realm name string
func GetRealmName(message string, command string) string {
    commandString := strings.Replace(message, command, "", 1)
    if commandString == "" {
        return "Ревущий фьорд"
    }
    return strings.TrimLeft(commandString, " ")
}

// GetRealmAndGuildNames returns realm and guild name strings
func GetRealmAndGuildNames(message string, command string) (string, string, error) {
    commandString := strings.Replace(message, command, "", 1)
    if commandString == "" {
        return "Ревущий фьорд", "Аэтернум", nil
    }
    s := strings.Split(commandString, ", ")
    if (len(s) < 2) {
        return "", "", errors.New("Команда введена неверно! Пожалуйста, попробуйте еще раз.")
    }
    return s[0], s[1], nil
}

func getRealms() (*[]Realm, error) {
    apiLink := fmt.Sprintf(apiRealmsLink, region, locale, apiToken)
    logInfo(apiLink)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    realms, err := getRealmsFromJSON([]byte(body))
    if err != nil {
        return nil, err
    }
    return &realms.RealmList, nil
}

func getGuildMembers(guildRealm *string, guildName *string) (*[]GuildMember, error) {
    apiLink := fmt.Sprintf(apiGuildMembersLink, region, strings.Replace(*guildRealm, " ", "%20", -1), 
        *guildName, locale, apiToken)
    logInfo(apiLink)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    gInfo, err := getGuildMembersFromJSON([]byte(body))
    if err != nil {
        return nil, err
    }
    members := gInfo.GuildMembersList
    // Fill string valuables
    gInfo.Side = factions[gInfo.SideInt]
    return refillMembers(&members), nil
}

func refillMembers(members *[]GuildMember) *[]GuildMember {
    var guildMembers []GuildMember
    c := make(chan GuildMember)
    m := *members
    for i := range m {
        go updateCharacter(&m[i], c)
    }
    for i := 0; i < len(*members); i++ {
        guildMembers = append(guildMembers, <-c)
    }
    defer logInfo("Members refilled")
    return &guildMembers
}

func updateCharacter(member *GuildMember, c chan GuildMember) {
    var newMember = new(GuildMember)
    m := *member
    m.Member.Class  = classes[m.Member.ClassInt]
    m.Member.Gender = genders[m.Member.GenderInt]
    m.Member.Race   = races[m.Member.RaceInt]
    items, err := getCharacterItems(&m.Member.Realm, &m.Member.Name)
    if (err != nil) {
        c <- m
        return
    }
    newMember.Member        = m.Member
    newMember.Member.Items  = *items
    newMember.Rank          = m.Rank
    c <- *newMember
}

func getCharacterItems(characterRealm *string, characterName *string) (*Items, error) {
    apiLink := fmt.Sprintf(apiCharacterItems, region, strings.Replace(*characterRealm, " ", "%20", -1), 
        *characterName, locale, apiToken)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    if strings.Contains(r.Status, "404") {
        return nil, errors.New(r.Status)
    }
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    character, err := getCharacterWithItemsFromJSON([]byte(body))
    if err != nil {
        return nil, err
    }
    return &character.Items, nil
}

func getRealmsFromJSON(body []byte) (*Realms, error) {
    var r = new(Realms)
    err := json.Unmarshal(body, &r)
    panicOnErr(err)
    return r, err
}

func getGuildMembersFromJSON(body []byte) (*GuildInfo, error) {
    var gi = new(GuildInfo)
    err := json.Unmarshal(body, &gi)
    panicOnErr(err)
    return gi, err
}

func getCharacterWithItemsFromJSON(body []byte) (*Character, error) {
    var c = new(Character)
    err := json.Unmarshal(body, &c)
    panicOnErr(err)
    return c, err
}

func getRealmByName(realmName string) (Realm, error) {
    logInfo("getRealmByName: " + realmName)
    realms, err := getRealms()
    if err != nil {
        return *new(Realm), err
    }
    for _, r := range *realms {
        if strings.ToLower(r.Name) == strings.ToLower(realmName) || 
           strings.ToLower(r.Slug) == strings.ToLower(realmName) {
            return r, nil
        }
    }
    return *new(Realm), errors.New("No such realm is present!")
}
