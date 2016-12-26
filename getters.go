package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func getRealms() (*[]Realm, error) {
	apiLink := fmt.Sprintf(WoWAPIRealmsLink, region, locale, wowAPIToken)
	respJSON, err := GetJSONResponse(apiLink)
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
	apiLink := fmt.Sprintf(WoWAPICharacterItemsLink, region, strings.Replace(*characterRealm, " ", "%20", -1),
		*characterName, locale, wowAPIToken)
	respJSON, err := GetJSONResponse(apiLink)
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
	apiLink := fmt.Sprintf(WoWAPICharacterProfsLink, region, strings.Replace(*characterRealm, " ", "%20", -1),
		*characterName, locale, wowAPIToken)
	respJSON, err := GetJSONResponse(apiLink)
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
	gAPILink := fmt.Sprintf(GoogleAPIShortenerLink, googleAPIToken)
	link := fmt.Sprintf(WoWArmoryLink, region, locale[:2], *rSlug, *cName)
	respJSON, err := PostJSONResponse(gAPILink, link)
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
	gAPILink := fmt.Sprintf(GoogleAPIShortenerLink, googleAPIToken)
	link := fmt.Sprintf(WoWArmoryProfLink, region, locale[:2], *rSlug, *cName, *pName)
	respJSON, err := PostJSONResponse(gAPILink, link)
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
	itemJSON := Get("Items", itemID)
	item = new(Item)
	cached := true
	if itemJSON == nil {
		apiLink := fmt.Sprintf(WoWAPIItemLink, region, itemID, locale, wowAPIToken)
		itemJSON, err = GetJSONResponse(apiLink)
		if err != nil {
			logInfo(err)
			return
		}
		err = Put("Items", itemID, itemJSON)
		if err != nil {
			logInfo(err)
			return
		}
		cached = false
	}
	if itemJSON == nil {
		err = errors.New("Null JSON! itemID = " + itemID + ", cached = " + strconv.FormatBool(cached))
		return
	}
	err = item.getItemFromJSON(&itemJSON)
	if err != nil {
		logInfo(err)
		return
	}
	item.Link = fmt.Sprintf(WowheadItemLink, itemID)
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
