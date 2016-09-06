package wow

import (
    "fmt"
    "errors"
    "strings"
    "bytes"
    "strconv"
    "sort"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
    "github.com/nezorflame/discord-wow-bot/consts"
)

var (
    logger          *log.Logger
    wowAPIToken     string
    googleAPIToken  string
    locale          string
    region          string
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
    GuildNewsList       []News          `json:"news"`
}

// GuildMember - struct for a WoW guild member
type GuildMember struct {
    Member              Character       `json:"character"`
    Rank                int             `json:"rank"`
}

// News - struct for any WoW guild news
type News struct {
    Type                string          `json:"type"`
    Character           string          `json:"character"`
    Timestamp           int             `json:"type"`
    ItemID              string          `json:"itemId"`
    Context             string          `json:"context"`
    BonusLists          []string        `json:"bonusLists"`
    Achievement         Achievement     `json:"achievement"`
}

// Achievement - struct for a WoW achievement
type Achievement struct {
    ID                  int             `json:"id"`
    Title               string          `json:"title"`
    Points              int             `json:"points"`
    Description         string          `json:"description"`
    RewardItems         []string        `json:"rewardItems"`
    Icon                string          `json:"icon"`
    Criteria            Criteria        `json:"criteria"`
    AccountWide         bool            `json:"accountWide"`
    FactionID           int             `json:"factionId"`
}

// Criteria - struct for a WoW achievement criteria
type Criteria struct {
    ID                  int             `json:"id"`
    Description         string          `json:"description"`
    OrderIndex          int             `json:"orderIndex"`
    Max                 int             `json:"max"`
}

// Character - struct for a WoW character
type Character struct {
    Name                string          `json:"name"` 
    Realm               string          `json:"realm"`
    RealmSlug           string
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
    Professions         Professions     `json:"professions"`
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

// Professions - struct for professions info for a character
type Professions struct {
    PrimaryProfs        []Profession    `json:"primary"`
    SecondaryProfs      []Profession    `json:"secondary"`
}

// Profession - struct for a profession info for a character
type Profession struct {
    ID                  int             `json:"id"`
    Name                string          `json:"name"`
    EngName             string
    Icon                string          `json:"icon"`
    Rank                int             `json:"rank"`
    Max                 int             `json:"max"`
    Recipes             []int           `json:"recipes"`
    Link                string
}

type googlAPIResponse struct {
    Kind                string          `json:"kind"`
    ID                  string          `json:"id"`
    LongURL             string          `json:"longUrl"`
}

var (
    classes             map[int]string
    factions            map[int]string
    races               map[int]string
    genders             map[int]string
    profNames           map[int]string
    addMembers          map[string]map[string]string
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
func InitializeWoWAPI(wowToken, googleToken *string) {
    logger = log.New(os.Stderr, "  ", log.Ldate|log.Ltime)
    wowAPIToken = *wowToken
    googleAPIToken = *googleToken
    locale = consts.Locale
    region = consts.Region
    // TODO: Rework to WoW API structures from response
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
    profNames = map[int]string{
        171 : "alchemy",
        164 : "blacksmithing",
        794 : "archaeology",
        185 : "cooking",
        333 : "enchanting",
        202 : "engineering",
        129 : "first-aid",
        356 : "fishing",
        182 : "herbalism",
        773 : "inscription",
        755 : "jewelcrafting",
        165 : "leatherworking",
        186 : "mining",
        393 : "skinning",
        197 : "tailoring",
    }
    addMembers = map[string]map[string]string{
        "Stormscale" : {"The Timekeepers" : "Madmaid"},
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

// GetGuildNews - function for getting the latest guild news
func GetGuildNews(realmName, guildName string) (*[]News, error) {
    gInfo, err := getGuildNews(&realmName, &guildName)
    return gInfo, err
}

// GetGuildMembers - function for receiving a list of guild members
func GetGuildMembers(realmName, guildName string) ([]map[string]string, error) {
    gMembers, err := getGuildMembers(&realmName, &guildName)
    if err != nil {
        return nil, err
    }
    gMembers, err = getAdditionalMembers(gMembers)
    if err != nil {
        return nil, err
    }
    gMembers = refillMembers(gMembers, "Items")

    var guildMembersList []map[string]string
    gMembersMap := make(map[string]map[string]string)
    var keys []string

    for _, m := range *gMembers {
        gMember := make(map[string]string)
        gMember["Name"] = m.Member.Name
        gMember["Level"] = strconv.Itoa(m.Member.Level)
        gMember["Class"] = m.Member.Class
        if specName := m.Member.Spec.Name; specName != "" {
            gMember["Spec"] = specName
        } else {
            gMember["Spec"] = "Нет инфы"
        }
        gMember["ItemLevel"] = strconv.Itoa(m.Member.Items.AvgItemLvlEq)
        gMembersMap[m.Member.Name] = gMember
        keys = append(keys, m.Member.Name)
    }

    sort.Strings(keys)
    for _, k := range keys {
        guildMembersList = append(guildMembersList, gMembersMap[k])
    }
    return guildMembersList, nil
}

// GetGuildProfs - function for receiving a list of guild professions
func GetGuildProfs(realmName string, guildName string) ([]map[string]string, error) {
    gMembers, err := getGuildMembers(&realmName, &guildName)
    if err != nil {
        return nil, err
    }
    gMembers, err = getAdditionalMembers(gMembers)
    if err != nil {
        return nil, err
    }
    gMembers = refillMembers(gMembers, "Profs")

    var guildProfsList []map[string]string
    gMembersMap := make(map[string]map[string]string)
    var keys []string

    for _, m := range *gMembers {
        gMember := make(map[string]string)
        gMember["Name"] = m.Member.Name
        switch len(m.Member.Professions.PrimaryProfs) {
            case 0:
                gMember["FirstProf"] = "Нет"
                gMember["FirstProfLevel"] = " "
                gMember["SecondProf"] = "Нет"
                gMember["SecondProfLevel"] = " "
            case 1:
                gMember["FirstProf"] = m.Member.Professions.PrimaryProfs[0].Name
                gMember["FirstProfLevel"] = strconv.Itoa(m.Member.Professions.PrimaryProfs[0].Rank) + 
                                            " | " + m.Member.Professions.PrimaryProfs[0].Link
                gMember["SecondProf"] = "Нет"
                gMember["SecondProfLevel"] = " "
            case 2:
                gMember["FirstProf"] = m.Member.Professions.PrimaryProfs[0].Name
                gMember["FirstProfLevel"] = strconv.Itoa(m.Member.Professions.PrimaryProfs[0].Rank) + 
                                            " | " + m.Member.Professions.PrimaryProfs[0].Link
                gMember["SecondProf"] = m.Member.Professions.PrimaryProfs[1].Name
                gMember["SecondProfLevel"] = strconv.Itoa(m.Member.Professions.PrimaryProfs[1].Rank) + 
                                            " | " + m.Member.Professions.PrimaryProfs[1].Link
        }
        gMembersMap[m.Member.Name] = gMember
        keys = append(keys, m.Member.Name)
    }

    sort.Strings(keys)
    for _, k := range keys {
        guildProfsList = append(guildProfsList, gMembersMap[k])
    }
    return guildProfsList, nil
}

// GetRealmName returns realm name string
func GetRealmName(message string, command string) string {
    commandString := strings.Replace(message, command, "", 1)
    if commandString == "" {
        return consts.GuildRealm
    }
    return strings.TrimLeft(commandString, " ")
}

// GetRealmAndGuildNames returns realm and guild name strings
func GetRealmAndGuildNames(message string, command string) (string, string, error) {
    commandString := strings.Replace(message, command, "", 1)
    if commandString == "" {
        return consts.GuildRealm, consts.GuildName, nil
    }
    s := strings.Split(commandString, ", ")
    if (len(s) < 2) {
        return "", "", errors.New("Команда введена неверно! Пожалуйста, попробуйте еще раз.")
    }
    return s[0], s[1], nil
}

func getRealms() (*[]Realm, error) {
    apiLink := fmt.Sprintf(consts.WoWAPIRealmsLink, region, locale, wowAPIToken)
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

func getGuildMembers(guildRealm, guildName *string) (*[]GuildMember, error) {
    logInfo("getting main guild members...")
    apiLink := fmt.Sprintf(consts.WoWAPIGuildMembersLink, region, strings.Replace(*guildRealm, " ", "%20", -1), 
        strings.Replace(*guildName, " ", "%20", -1), locale, wowAPIToken)
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
    return &members, nil
}

func getGuildNews(guildRealm, guildName *string) (*[]News, error) {
    return new([]News), nil
}

func getAdditionalMembers(guildMembers *[]GuildMember)  (*[]GuildMember, error) {
    logInfo("getting additional guild members...")
    var addGMembers []GuildMember
    for _, m := range *guildMembers {
        addGMembers = append(addGMembers, m)
    }
    for realm, m := range addMembers {
        for guild, character := range m {
            apiLink := fmt.Sprintf(consts.WoWAPIGuildMembersLink, region, strings.Replace(realm, " ", "%20", -1), 
                strings.Replace(guild, " ", "%20", -1), locale, wowAPIToken)
            r, err := http.Get(apiLink)
            panicOnErr(err)
            defer r.Body.Close()
            body, err := ioutil.ReadAll(r.Body)
            panicOnErr(err)
            addGInfo, err := getGuildMembersFromJSON([]byte(body))
            if err != nil {
                return nil, err
            }
            for _, member := range addGInfo.GuildMembersList {
                if member.Member.Name == character {
                    addGMembers = append(addGMembers, member)
                }
            }
        }
    }
    // Fill string valuables
    return &addGMembers, nil
}

func refillMembers(members *[]GuildMember, t string) *[]GuildMember {
    var guildMembers []GuildMember
    c := make(chan GuildMember)
    m := *members
    for i := range m {
        go updateCharacter(&m[i], t, c)
    }
    for i := 0; i < len(*members); i++ {
        guildMembers = append(guildMembers, <-c)
    }
    defer logInfo("Members refilled with", t)
    return &guildMembers
}

func updateCharacter(member *GuildMember, t string, c chan GuildMember) {
    var newMember = new(GuildMember)
    var items *Items
    var profs *Professions
    var err error
    m := *member
    m.Member.Class  = classes[m.Member.ClassInt]
    m.Member.Gender = genders[m.Member.GenderInt]
    m.Member.Race   = races[m.Member.RaceInt]
    switch t {
        case "Items":
            items, err = getCharacterItems(&m.Member.Realm, &m.Member.Name)
        case "Profs":
            profs, err = getCharacterProfessions(&m.Member.Realm, &m.Member.Name)
    }
    if (err != nil) {
        c <- m
        logInfo(err)
        return
    }
    if (err != nil) {
        c <- m
        logInfo(err)
        return
    }
    newMember.Member = m.Member
    switch t {
        case "Items":
            newMember.Member.Items = *items
        case "Profs":
            newMember.Member.Professions = *profs
    }
    newMember.Rank = m.Rank
    c <- *newMember
}

func getCharacterItems(characterRealm *string, characterName *string) (*Items, error) {
    apiLink := fmt.Sprintf(consts.WoWAPICharacterItemsLink, region, strings.Replace(*characterRealm, " ", "%20", -1), 
        *characterName, locale, wowAPIToken)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    if strings.Contains(r.Status, "404") {
        return nil, errors.New(r.Status)
    }
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    character, err := getCharacterFromJSON([]byte(body))
    if err != nil {
        return nil, err
    }
    character.RealmSlug, err = getRealmSlugByName(characterRealm)
    if err != nil {
        return nil, err
    }
    return &character.Items, nil
}

func getCharacterProfessions(characterRealm *string, characterName *string) (*Professions, error) {
    apiLink := fmt.Sprintf(consts.WoWAPICharacterProfsLink, region, strings.Replace(*characterRealm, " ", "%20", -1), 
        *characterName, locale, wowAPIToken)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    if strings.Contains(r.Status, "404") {
        return nil, errors.New(r.Status)
    }
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    character, err := getCharacterFromJSON([]byte(body))
    if err != nil {
        logInfo(err.Error)
        return nil, err
    }
    character.RealmSlug, err = getRealmSlugByName(characterRealm)
    if err != nil {
        return nil, err
    }
    var profs = new(Professions)
    for _, p := range character.Professions.PrimaryProfs {
        var prof = new(Profession)
        prof = &p
        prof.EngName = profNames[p.ID]
        shortLink, err := getProfShortLink(&character.RealmSlug, characterName, &p.EngName)
        if err != nil {
            logInfo(err)
            return &character.Professions, err
        }
        prof.Link = shortLink
        profs.PrimaryProfs = append(profs.PrimaryProfs, *prof)
    }
    for _, p := range character.Professions.SecondaryProfs {
        var prof = new(Profession)
        prof = &p
        prof.EngName = profNames[p.ID]
        profs.SecondaryProfs = append(profs.SecondaryProfs, *prof)
    }
    return profs, nil
}

func getProfShortLink(rSlug, cName, pName *string) (string, error) {
    link := fmt.Sprintf(consts.WoWArmoryProfLink, region, locale[:2], *rSlug, *cName, *pName)
    apiLink := fmt.Sprintf(consts.GoogleAPIShortenerLink, googleAPIToken)
    link = `{"longUrl": "` + link + `"}`
    var jsonStr = []byte(link)
    req, err := http.NewRequest("POST", apiLink, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    panicOnErr(err)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    panicOnErr(err)
    shortLink, err := getURLFromJSON([]byte(body))
    if err != nil {
        logInfo(err.Error)
        return "", err
    }

    return *shortLink, nil
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

func getCharacterFromJSON(body []byte) (*Character, error) {
    var c = new(Character)
    err := json.Unmarshal(body, &c)
    panicOnErr(err)
    return c, err
}

func getURLFromJSON(body []byte) (*string, error) {
    var apiResponse = new(googlAPIResponse)
    err := json.Unmarshal(body, &apiResponse)
    panicOnErr(err)
    return &apiResponse.ID, err
}

func getRealmByName(realmName string) (Realm, error) {
    logDebug("getRealmByName: " + realmName)
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

func getRealmSlugByName(realmName *string) (string, error) {
    realms, err := getRealms()
    if err != nil {
        return "", err
    }
    for _, r := range *realms {
        if strings.ToLower(r.Name) == strings.ToLower(*realmName) {
            return r.Slug, nil
        }
    }
    return "", errors.New("No such realm is present!")
}
