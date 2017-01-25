package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

func getRealms() (*[]Realm, error) {
	apiLink := fmt.Sprintf(o.APIRealmsLink, o.GuildRegion, o.GuildLocale, o.WoWToken)
	respJSON, err := GetJSONResponse(apiLink)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	var realms Realms
	if err = realms.unmarshal(&respJSON); err != nil {
		return nil, err
	}
	return &realms.RealmList, nil
}

func getCharacterItems(characterRealm *string, characterName *string) (*Items, error) {
	apiLink := fmt.Sprintf(o.APICharItemsLink, o.GuildRegion, strings.Replace(*characterRealm, " ", "%20", -1),
		*characterName, o.GuildLocale, o.WoWToken)
	respJSON, err := GetJSONResponse(apiLink)
	if err != nil {
		return nil, err
	}
	var character Character
	if err = character.unmarshal(&respJSON); err != nil {
		return nil, err
	}
	return &character.Items, nil
}

func getCharacterProfessions(characterRealm *string, characterName *string) (*Professions, error) {
	apiLink := fmt.Sprintf(o.APICharProfsLink, o.GuildRegion, strings.Replace(*characterRealm, " ", "%20", -1),
		*characterName, o.GuildLocale, o.WoWToken)
	respJSON, err := GetJSONResponse(apiLink)
	if err != nil {
		glog.Info(err)
		return nil, err
	}
	var character Character
	if err := character.unmarshal(&respJSON); err != nil {
		glog.Info(err)
		return nil, err
	}
	var profs = new(Professions)
	for _, p := range character.Professions.PrimaryProfs {
		var prof = new(Profession)
		prof = &p
		prof.EngName = profNames[p.ID]
		shortLink, err := getProfShortLink(character.RealmSlug, *characterName, p.EngName)
		if err != nil {
			glog.Info(err)
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

func getArmoryLink(rSlug, cName string) (string, error) {
	gAPILink := fmt.Sprintf(o.GoogleShortenerLink, o.GoogleToken)
	link := fmt.Sprintf(o.ArmoryCharLink, o.GuildRegion, o.GuildLocale[:2], rSlug, cName)
	respJSON, err := PostJSONResponse(gAPILink, link)
	if err != nil {
		glog.Info(err)
		return "", err
	}
	shortLink, err := getURLFromJSON(&respJSON)
	if err != nil {
		glog.Info(err)
		return "", err
	}
	return shortLink, nil
}

func getProfShortLink(rSlug, cName, pName string) (string, error) {
	gAPILink := fmt.Sprintf(o.GoogleShortenerLink, o.GoogleToken)
	link := fmt.Sprintf(o.ArmoryProfLink, o.GuildRegion, o.GuildLocale[:2], rSlug, cName, pName)
	respJSON, err := PostJSONResponse(gAPILink, link)
	if err != nil {
		glog.Info(err)
		return "", err
	}
	shortLink, err := getURLFromJSON(&respJSON)
	if err != nil {
		glog.Info(err)
		return "", err
	}
	return shortLink, nil
}

func getItemByID(itemID string) (item *Item, err error) {
	itemJSON := Get("Items", itemID)
	item = new(Item)
	cached := true
	if itemJSON == nil {
		apiLink := fmt.Sprintf(o.APIItemLink, o.GuildRegion, itemID, o.GuildLocale, o.WoWToken)
		itemJSON, err = GetJSONResponse(apiLink)
		if err != nil {
			glog.Info(err)
			return
		}
		err = Put("Items", itemID, itemJSON)
		if err != nil {
			glog.Info(err)
			return
		}
		cached = false
	}
	if itemJSON == nil {
		err = errors.New("Null JSON! itemID = " + itemID + ", cached = " + strconv.FormatBool(cached))
		return
	}
	if err = item.unmarshal(&itemJSON); err != nil {
		glog.Info(err)
		return
	}
	item.Link = fmt.Sprintf(o.WowheadItemLink, itemID)
	return
}

func getRealmByName(realmName string) (Realm, error) {
	glog.Infof("getRealmByName: %s", realmName)
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
	return *new(Realm), errors.New("No such realm is present")
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
	return "", errors.New("No such realm is present")
}
