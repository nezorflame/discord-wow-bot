package wow

import (
    "fmt"
    "bytes"
    "io/ioutil"
    "net/http"
    "errors"
    "strings"
    "time"
    "github.com/nezorflame/discord-wow-bot/consts"
)

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

func getGuildNews(guildRealm, guildName *string) (gNews NewsList, err error) {
    logInfo("getting guild news...")
    apiLink := fmt.Sprintf(consts.WoWAPIGuildNewsLink, region, strings.Replace(*guildRealm, " ", "%20", -1), 
        strings.Replace(*guildName, " ", "%20", -1), locale, wowAPIToken)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    gInfo, err := getGuildInfoFromJSON([]byte(body))
    if err != nil {
        return
    }
    now := time.Now()
    before := now.Add(-24 * time.Hour)
    // Fill string valuables
    gInfo.Side = factions[gInfo.SideInt]
    for _, n := range gInfo.GuildNewsList {
        eventTime := time.Unix(n.Timestamp / 1000, 0)
        utc, err := time.LoadLocation(consts.Timezone)
        panicOnErr(err)
        n.EventTime = eventTime.In(utc)
        if inTimeSpan(before, now, n.EventTime) {
            gNews = append(gNews, n)
        }
    }
    return
}

func getUpdatedGuildNews(realmName, guildName string) (*NewsList, error) {
    var gNews NewsList
    gNews, err := getGuildNews(&realmName, &guildName)
    if err != nil {
        return nil, err
    }
    done := make(chan NewsList, 1)
    go gNews.refillNews(done)
    gNews = <-done
    gNews = gNews.sortGuildNewsByTimestamp()
    logInfo("Got updated guild news")
    return &gNews, nil
}

func getGuildMembers(guildRealm, guildName *string) (ml MembersList, err error) {
    logInfo("getting main guild members...")
    apiLink := fmt.Sprintf(consts.WoWAPIGuildMembersLink, region, strings.Replace(*guildRealm, " ", "%20", -1), 
        strings.Replace(*guildName, " ", "%20", -1), locale, wowAPIToken)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    gInfo, err := getGuildInfoFromJSON([]byte(body))
    if err != nil {
        return
    }
    // Fill string valuables
    gInfo.Side = factions[gInfo.SideInt]
    ml = gInfo.GuildMembersList
    return
}

func (ml *MembersList) getAdditionalMembers() error {
    logInfo("getting additional guild members...")
    for realm, m := range addMembers {
        for guild, character := range m {
            apiLink := fmt.Sprintf(consts.WoWAPIGuildMembersLink, region, strings.Replace(realm, " ", "%20", -1), 
                strings.Replace(guild, " ", "%20", -1), locale, wowAPIToken)
            r, err := http.Get(apiLink)
            panicOnErr(err)
            defer r.Body.Close()
            body, err := ioutil.ReadAll(r.Body)
            panicOnErr(err)
            addGInfo, err := getGuildInfoFromJSON([]byte(body))
            if err != nil {
                return err
            }
            for _, member := range addGInfo.GuildMembersList {
                if member.Member.Name == character {
                    *ml = append(*ml, member)
                }
            }
        }
    }
    return nil
}

func (ml *MembersList) refillMembers(t string, done chan MembersList) {
    var guildMembers MembersList
    c := make(chan GuildMember, len(*ml))
    for _, m := range *ml {
        go updateCharacter(&m, t, c)
    }
    for i := 0; i < len(*ml); i++ {
        guildMembers = append(guildMembers, <-c)
    }
    logInfo("Members refilled with", t)
    done <- guildMembers
}

func (nl *NewsList) refillNews(done chan NewsList){
    var guildNews NewsList
    for _, n := range *nl {
        if n.Type == "itemLoot" {
            guildNews = append(guildNews, n)
        }
    }
    l := len(guildNews)
    c := make(chan News, l)
    for _, n := range guildNews {
        go updateNews(n, c)
    }
    guildNews = make(NewsList, 0)
    for i := 0; i < l; i++ {
        guildNews = append(guildNews, <-c)
    }
    logInfo("News refilled")
    done <- guildNews
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
    m.Member.RealmSlug, err = getRealmSlugByName(&m.Member.Realm)
    if err != nil {
        c <- m
        logInfo(err)
        return
    }
    shortLink, err := getArmoryLink(&m.Member.RealmSlug, &m.Member.Name)
    if err != nil {
        c <- m
        logInfo(err)
        return
    }
    m.Member.Link = shortLink
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

func updateNews(newsrecord News, c chan News) {
    if newsrecord.Type == "itemLoot" {
        item, err := getItemByID(&newsrecord.ItemID)
        if (err != nil) {
            c <- newsrecord
            logInfo(err)
            return
        }
        newsrecord.ItemInfo = *item
    }
    c <- newsrecord
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
        logInfo(err)
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

func getArmoryLink(rSlug, cName *string) (string, error) {
    link := fmt.Sprintf(consts.WoWArmoryLink, region, locale[:2], *rSlug, *cName)
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
        logInfo(err)
        return "", err
    }

    return *shortLink, nil
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
        logInfo(err)
        return "", err
    }

    return *shortLink, nil
}

func getItemByID(itemID *int) (*Item, error) {
    apiLink := fmt.Sprintf(consts.WoWAPIItemLink, region, *itemID, locale, wowAPIToken)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    if strings.Contains(r.Status, "404") {
        return new(Item), errors.New(r.Status)
    }
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    item, err := getItemFromJSON([]byte(body))
    if err != nil {
        logInfo(err)
        return new(Item), err
    }
    item.Link = apiLink
    return item, nil
}

func getItemQualityByID(itemID *int) (int, error) {
    apiLink := fmt.Sprintf(consts.WoWAPIItemLink, region, *itemID, locale, wowAPIToken)
    r, err := http.Get(apiLink)
    panicOnErr(err)
    if strings.Contains(r.Status, "404") {
        return -1, errors.New(r.Status)
    }
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    item, err := getItemFromJSON([]byte(body))
    if err != nil {
        logInfo(err)
        return -1, err
    }
    return item.Quality, nil
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
