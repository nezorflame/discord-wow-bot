package wow

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nezorflame/discord-wow-bot/consts"
	"github.com/nezorflame/discord-wow-bot/db"
	"github.com/nezorflame/discord-wow-bot/net"
)

func getRealms() (*[]Realm, error) {
	apiLink := fmt.Sprintf(consts.WoWAPIRealmsLink, region, locale, wowAPIToken)
	respJSON, err := net.GetJSONResponse(apiLink)
	if err != nil {
		logOnErr(err)
		return nil, err
	}
	realms, err := getRealmsFromJSON(respJSON)
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
	gNews = gNews.refillNews()
	gNews = gNews.SortGuildNews()
	logInfo("Got updated guild news")
	return &gNews, nil
}

func getGuildNewsList(guildRealm, guildName *string) (gNews NewsList, err error) {
	logInfo("getting guild news...")
	apiLink := fmt.Sprintf(consts.WoWAPIGuildNewsLink, region, strings.Replace(*guildRealm, " ", "%20", -1),
		strings.Replace(*guildName, " ", "%20", -1), locale, wowAPIToken)
	respJSON, err := net.GetJSONResponse(apiLink)
	if err != nil {
		logOnErr(err)
		return nil, err
	}
	gInfo, err := getGuildInfoFromJSON(respJSON)
	if err != nil {
		return
	}
	now := time.Now()
	before := now.Add(time.Duration(-5 * time.Minute))
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
	gInfo, cached, err := getGuildInfo(&realmName, &guildName)
	if err != nil {
		return nil, err
	}
	gMembers := gInfo.GuildMembersList
	if !cached {
		logInfo("Got", len(gMembers), "guild members from API. Filling the gaps...")
		gMembers = gMembers.refillMembers(option)
		logInfo("Saving guild members into cache...")
		gInfo.GuildMembersList = gMembers
		giJSON, err := getJSONFromGuildInfo(&gInfo)
		err = db.Put("Main", consts.GuildMembersBucketKey, giJSON)
		logOnErr(err)
	} else {
		logInfo("Got", len(gMembers), "guild members from cache")
	}
	gMembers = gMembers.SortGuildMembers(params)
	logInfo("Sorted guild members")
	return &gMembers, nil
}

func getGuildInfo(guildRealm, guildName *string) (gInfo GuildInfo, cached bool, err error) {
	logInfo("getting main guild members...")
	cached = true
	membersJSON := db.Get("Main", consts.GuildMembersBucketKey)
	if membersJSON == nil {
		logInfo("No cache is present, getting from API...")
		cached = false
		apiLink := fmt.Sprintf(consts.WoWAPIGuildMembersLink, region, strings.Replace(*guildRealm, " ", "%20", -1),
			strings.Replace(*guildName, " ", "%20", -1), locale, wowAPIToken)
		membersJSON, err = net.GetJSONResponse(apiLink)
		if err != nil {
			logOnErr(err)
			return
		}
	}
	gi, err := getGuildInfoFromJSON(membersJSON)
	if err != nil {
		return
	}
	err = gi.GuildMembersList.getAdditionalMembers()
	gInfo = *gi
	return
}

func (ml *MembersList) getAdditionalMembers() error {
	logInfo("getting additional guild members...")
	for realm, m := range addMembers {
		for guild, character := range m {
			apiLink := fmt.Sprintf(consts.WoWAPIGuildMembersLink, region, strings.Replace(realm, " ", "%20", -1),
				strings.Replace(guild, " ", "%20", -1), locale, wowAPIToken)
			respJSON, err := net.GetJSONResponse(apiLink)
			if err != nil {
				logOnErr(err)	
				return err
			}
			addGInfo, err := getGuildInfoFromJSON(respJSON)
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

func (ml *MembersList) refillMembers(t string) (guildMembers MembersList) {
	var wg sync.WaitGroup
	wg.Add(len(*ml))
	for _, m := range *ml {
		go func(m GuildMember) {
			defer wg.Done()
			gMember := updateCharacter(m, t)
			guildMembers = append(guildMembers, gMember)
		}(m)
	}
	logInfo("Members refilled with", t)
	wg.Wait()
	return
}

func updateCharacter(member GuildMember, t string) (m GuildMember) {
	var items *Items
	var profs *Professions
	var err error
	m.Member = member.Member
	m.Rank = member.Rank
	m.Member.Class = classes[m.Member.ClassInt]
	m.Member.Gender = genders[m.Member.GenderInt]
	m.Member.Race = races[m.Member.RaceInt]
	m.Member.RealmSlug, err = getRealmSlugByName(&m.Member.Realm)
	if err != nil {
		logInfo("updateCharacter(): unable to get realm slug:", err)
		return member
	}
	shortLink, err := getArmoryLink(&m.Member.RealmSlug, &m.Member.Name)
	if err != nil {
		logInfo("updateCharacter(): unable to get Armory link:", err)
		return member
	}
	m.Member.Link = shortLink
	switch t {
	case "Items":
		items, err = getCharacterItems(&m.Member.Realm, &m.Member.Name)
	case "Profs":
		profs, err = getCharacterProfessions(&m.Member.Realm, &m.Member.Name)
	}
	if err != nil {
		logInfo("updateCharacter(): unable to get", t+":", err)
		return member
	}
	switch t {
	case "Items":
		m.Member.Items = *items
	case "Profs":
		m.Member.Professions = *profs
	}
	return
}

func (nl *NewsList) refillNews() (guildNews NewsList) {
	var wg sync.WaitGroup
	wg.Add(len(*nl))
	for _, n := range *nl {
		go func(n News) {
			defer wg.Done()
			news := updateNews(n)
			guildNews = append(guildNews, news)
		}(n)
	}
	logInfo("News refilled")
	wg.Wait()
	return
}

func updateNews(newsrecord News) (news News) {
	if newsrecord.Type == "itemLoot" {
		item, err := getItemByID(strconv.Itoa(newsrecord.ItemID))
		if err != nil {
			logInfo("updateCharacter(): unable to get item by its ID =", newsrecord.ItemID, ":", err)
			return newsrecord
		}
		newsrecord.ItemInfo = *item
	}
	return newsrecord
}

func getCharacterItems(characterRealm *string, characterName *string) (*Items, error) {
	apiLink := fmt.Sprintf(consts.WoWAPICharacterItemsLink, region, strings.Replace(*characterRealm, " ", "%20", -1),
		*characterName, locale, wowAPIToken)
	respJSON, err := net.GetJSONResponse(apiLink)
	if err != nil {
		return nil, err
	}
	character, err := getCharacterFromJSON(respJSON)
	if err != nil {
		return nil, err
	}
	return &character.Items, nil
}

func getCharacterProfessions(characterRealm *string, characterName *string) (*Professions, error) {
	apiLink := fmt.Sprintf(consts.WoWAPICharacterProfsLink, region, strings.Replace(*characterRealm, " ", "%20", -1),
		*characterName, locale, wowAPIToken)
	respJSON, err := net.GetJSONResponse(apiLink)
	if err != nil {
		logInfo(err)
		return nil, err
	}
	character, err := getCharacterFromJSON(respJSON)
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
	gAPILink := fmt.Sprintf(consts.GoogleAPIShortenerLink, googleAPIToken)
	link := fmt.Sprintf(consts.WoWArmoryLink, region, locale[:2], *rSlug, *cName)
	respJSON, err := net.PostJSONResponse(gAPILink, link)
	if err != nil {
		logInfo(err)
		return "", err
	}
	shortLink, err := getURLFromJSON(respJSON)
	if err != nil {
		logInfo(err)
		return "", err
	}
	return *shortLink, nil
}

func getProfShortLink(rSlug, cName, pName *string) (string, error) {
	gAPILink := fmt.Sprintf(consts.GoogleAPIShortenerLink, googleAPIToken)
	link := fmt.Sprintf(consts.WoWArmoryProfLink, region, locale[:2], *rSlug, *cName, *pName)
	respJSON, err := net.PostJSONResponse(gAPILink, link)
	if err != nil {
		logInfo(err)
		return "", err
	}
	shortLink, err := getURLFromJSON(respJSON)
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
		itemJSON, err := net.GetJSONResponse(apiLink)
		if err != nil {
			logInfo(err)
			return nil, err
		}
		err = db.Put("Items", itemID, itemJSON)
		if err != nil {
			return nil, err
		}
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
