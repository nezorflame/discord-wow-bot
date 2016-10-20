package wow

import (
	"errors"
	"fmt"
	"strings"
	"strconv"

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
	var realms Realms
	err = realms.getRealmsFromJSON(&respJSON)
	if err != nil {
		return nil, err
	}
	return &realms.RealmList, nil
}

func getCharacterItems(characterRealm *string, characterName *string) (*Items, error) {
	apiLink := fmt.Sprintf(consts.WoWAPICharacterItemsLink, region, strings.Replace(*characterRealm, " ", "%20", -1),
		*characterName, locale, wowAPIToken)
	respJSON, err := net.GetJSONResponse(apiLink)
	if err != nil {
		return nil, err
	}
	character, err := getCharacterFromJSON(&respJSON)
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
	character, err := getCharacterFromJSON(&respJSON)
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
	shortLink, err := getURLFromJSON(&respJSON)
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
	shortLink, err := getURLFromJSON(&respJSON)
	if err != nil {
		logInfo(err)
		return "", err
	}
	return *shortLink, nil
}

func getItemByID(itemID string) (item *Item, err error) {
	itemJSON := db.Get("Items", itemID)
	item = new(Item)
	flag := false
	if itemJSON == nil {
		apiLink := fmt.Sprintf(consts.WoWAPIItemLink, region, itemID, locale, wowAPIToken)
		itemJSON, err = net.GetJSONResponse(apiLink)
		if err != nil {
			logInfo(err)
			return
		}
		err = db.Put("Items", itemID, itemJSON)
		if err != nil {
			logInfo(err)
			return
		}
		flag = true
	}
	logInfo(itemID, flag)
	if itemJSON == nil {
		err = errors.New("Null JSON! itemID = " + itemID + ", flag = " + strconv.FormatBool(flag))
		return
	}
	err = item.getItemFromJSON(&itemJSON)
	if err != nil {
		logInfo(err)
		return
	}
	item.Link = fmt.Sprintf(consts.WowheadItemLink, itemID)
	return
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
