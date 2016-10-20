package wow

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
    "strconv"
	"time"

	"github.com/nezorflame/discord-wow-bot/consts"
	"github.com/nezorflame/discord-wow-bot/db"
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

func getGuildNews(realmName, guildName string) (*NewsList, error) {
	gNews, err := getGuildNewsList(&realmName, &guildName)
	if err != nil {
		return nil, err
	}
	done := make(chan NewsList, 1)
	go gNews.refillNews(done)
	gNews = <-done
	gNews = gNews.SortGuildNews()
	logInfo("Got updated guild news")
	return &gNews, nil
}

func getGuildNewsList(guildRealm, guildName *string) (gNews NewsList, err error) {
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
	before := now.AddDate(0, 0, -3)
	// Fill string valuables
	gInfo.Side = factions[gInfo.SideInt]
	for _, n := range gInfo.GuildNewsList {
		eventTime := time.Unix(n.Timestamp/1000, 0)
		utc, err := time.LoadLocation(consts.Timezone)
		panicOnErr(err)
		n.EventTime = eventTime.In(utc)
		if inTimeSpan(before, now, eventTime) {
			gNews = append(gNews, n)
		}
	}
	return
}

func getGuildMembers(realmName, guildName, option string, params []string) (*MembersList, error) {
	gMembers, err := getGuildMembersList(&realmName, &guildName)
	if err != nil {
		return nil, err
	}
	logInfo("Got", len(gMembers), "guild members in total. Filling the gaps...")
	done := make(chan MembersList, 1)
	go gMembers.refillMembers(option, done)
	gMembers = <-done
	gMembers = gMembers.SortGuildMembers(params)
	logInfo("Got sorted guild members")
	return &gMembers, nil
}

func getGuildMembersList(guildRealm, guildName *string) (ml MembersList, err error) {
	logInfo("getting main guild members...")
	membersJSON := db.Get("Main", consts.GuildMembersBucketKey)
	if membersJSON == nil {
		logInfo("No cache is present, getting from API...")
		apiLink := fmt.Sprintf(consts.WoWAPIGuildMembersLink, region, strings.Replace(*guildRealm, " ", "%20", -1),
			strings.Replace(*guildName, " ", "%20", -1), locale, wowAPIToken)
		r, err := http.Get(apiLink)
		panicOnErr(err)
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		panicOnErr(err)
		membersJSON = []byte(body)
        err = db.Put("Main", consts.GuildMembersBucketKey, membersJSON)
        logOnErr(err)
	}
	gInfo, err := getGuildInfoFromJSON(membersJSON)
	if err != nil {
		return
	}
	ml = gInfo.GuildMembersList
	err = ml.getAdditionalMembers()
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
	// c := make(chan GuildMember, len(*ml))
	// for _, m := range *ml {
	// 	go updateCharacter(&m, t, c)
	// }
	c := fillMembers(t, *ml)
	for m := range c {
		guildMembers = append(guildMembers, m)
	}
	logInfo("Members refilled with", t)
	done <- guildMembers
}

func fillMembers(t string, ml MembersList) <-chan GuildMember {
	out := make(chan GuildMember)
	var count int
	for _, m := range ml {
    	go func() {
            out <- updateCharacter(&m, t)
			count++
    	}()
	}
	for {
		if count == len(ml) {
			close(out)
			break
		}
	}
    return out
}

func updateCharacter(member *GuildMember, t string) GuildMember {
	var newMember = new(GuildMember)
	var items *Items
	var profs *Professions
	var err error
	m := *member
	m.Member.Class = classes[m.Member.ClassInt]
	m.Member.Gender = genders[m.Member.GenderInt]
	m.Member.Race = races[m.Member.RaceInt]
	m.Member.RealmSlug, err = getRealmSlugByName(&m.Member.Realm)
	if err != nil {
		logInfo(err)
		return m
	}
	shortLink, err := getArmoryLink(&m.Member.RealmSlug, &m.Member.Name)
	if err != nil {
		logInfo(err)
		return m
	}
	m.Member.Link = shortLink
	switch t {
	case "Items":
		items, err = getCharacterItems(&m.Member.Realm, &m.Member.Name)
	case "Profs":
		profs, err = getCharacterProfessions(&m.Member.Realm, &m.Member.Name)
	}
	if err != nil {
		logInfo(err)
		return m
	}
	if err != nil {
		logInfo(err)
		return m
	}
	newMember.Member = m.Member
	switch t {
	case "Items":
		newMember.Member.Items = *items
	case "Profs":
		newMember.Member.Professions = *profs
	}
	newMember.Rank = m.Rank
	return *newMember
}

func (nl *NewsList) refillNews(done chan NewsList) {
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

func fillNews(nl NewsList) <-chan News {
	out := make(chan News)
    go func() {
        for _, n := range nl {
            out <- n
        }
        close(out)
    }()
    return out
}

func updateNews(newsrecord News, c chan News) {
	if newsrecord.Type == "itemLoot" {
		item, err := getItemByID(strconv.Itoa(newsrecord.ItemID))
		if err != nil {
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

func getItemByID(itemID string) (item *Item, err error) {
	itemJSON := db.Get("Items", itemID)
    if itemJSON == nil {
        apiLink := fmt.Sprintf(consts.WoWAPIItemLink, region, itemID, locale, wowAPIToken)
        r, err := http.Get(apiLink)
        panicOnErr(err)
        if strings.Contains(r.Status, "404") {
            return new(Item), errors.New(r.Status)
        }
        defer r.Body.Close()
        body, err := ioutil.ReadAll(r.Body)
        panicOnErr(err)
		itemJSON = []byte(body)
        err = db.Put("Items", itemID, itemJSON)
        panicOnErr(err)
    }
	item, err = getItemFromJSON(itemJSON)
	if err != nil {
		logInfo(err)
		return new(Item), err
	}
	item.Link = fmt.Sprintf(consts.WowheadItemLink, itemID)
	return item, nil
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
