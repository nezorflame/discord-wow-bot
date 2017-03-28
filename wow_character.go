package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

// UpdateCharacter gets the additional info for the WoW character
func (char *Character) UpdateCharacter(wg *sync.WaitGroup) {
	var err error

	char.Class = o.WoWClasses[char.ClassInt]
	char.Gender = o.WoWGenders[char.GenderInt]
	char.Race = o.WoWRaces[char.RaceInt]

	if err = char.SetRealmSlugByName(); err != nil {
		glog.Errorf("Unable to get realm slug for character %s: %s", char.Name, err)
	}

	if err = char.SetArmoryLink(); err != nil {
		glog.Errorf("Unable to get Armory link for character %s: %s", char.Name, err)
	}

	if err = char.SetCharacterItems(); err != nil {
		glog.Errorf("Unable to get items for character %s: %s", char.Name, err)
	}

	if err = char.SetCharacterProfessions(); err != nil {
		glog.Errorf("Unable to get profs for character %s: %s", char.Name, err)
	}

	wg.Done()

	return
}

// SetRealmSlugByName gets the realm slug name and sets it into the character
func (char *Character) SetRealmSlugByName() (err error) {
	var realms Realms

	if realms, err = getRealms(); err != nil {
		return
	}

	for _, r := range realms.RealmList {
		if strings.ToLower(r.Name) == strings.ToLower(char.Realm) {
			char.RealmSlug = r.Slug
			return
		}
	}

	return errors.New("No such realm is present")
}

// SetArmoryLink gets the Armory link and sets it into the character
func (char *Character) SetArmoryLink() (err error) {
	var shortLink string

	longLink := fmt.Sprintf(o.ArmoryCharLink, o.GuildRegion, o.GuildLocale[:2], char.RealmSlug, char.Name)
	if shortLink, err = GetShortLink(longLink); err != nil {
		glog.Errorf("Unable to get short link for a character: %s", err)
		return
	}

	char.Link = shortLink

	return
}

// SetCharacterItems gets the character items and sets them into the character
func (char *Character) SetCharacterItems() (err error) {
	var (
		respJSON []byte
	)

	apiLink := fmt.Sprintf(o.APICharItemsLink, o.GuildRegion, strings.Replace(char.Realm, " ", "%20", -1),
		char.Name, o.GuildLocale, o.WoWToken)
	if respJSON, err = GetJSONResponse(apiLink); err != nil {
		glog.Errorf("Unable to get JSON response: %s", err)
		return
	}

	if err = char.Unmarshal(respJSON); err != nil {
		glog.Errorf("Unable to unmarshal character from JSON: %s", err)
	}

	return
}

// SetCharacterProfessions gets the character professions and sets them into the character
func (char *Character) SetCharacterProfessions() (err error) {
	var (
		respJSON []byte
	)

	apiLink := fmt.Sprintf(o.APICharProfsLink, o.GuildRegion, strings.Replace(char.Realm, " ", "%20", -1),
		char.Name, o.GuildLocale, o.WoWToken)
	if respJSON, err = GetJSONResponse(apiLink); err != nil {
		glog.Errorf("Unable to get JSON response: %s", err)
		return
	}

	if err = char.Unmarshal(respJSON); err != nil {
		glog.Errorf("Unable to unmarshal character from JSON: %s", err)
		return
	}

	for _, p := range char.Professions.Primary {
		p.EngName = o.WoWProfessions[p.ID]
		longLink := fmt.Sprintf(o.ArmoryProfLink, o.GuildRegion, o.GuildLocale[:2], char.RealmSlug, char.Name, p.EngName)

		if shortLink, pErr := GetShortLink(longLink); pErr != nil {
			glog.Errorf("Unable to get short link for profession: %s", pErr)
		} else {
			p.Link = shortLink
		}
	}

	for _, p := range char.Professions.Secondary {
		p.EngName = o.WoWProfessions[p.ID]
	}

	return
}

// SetCharacterNewsFeed gets the character news feed and sets it into the character
func (char *Character) SetCharacterNewsFeed() (err error) {
	var (
		respJSON []byte

		wg sync.WaitGroup
	)

	apiLink := fmt.Sprintf(o.APICharNewsLink, o.GuildRegion, strings.Replace(char.Realm, " ", "%20", -1),
		strings.Replace(char.Name, " ", "%20", -1), o.GuildLocale, o.WoWToken)
	if respJSON, err = GetJSONResponse(apiLink); err != nil {
		glog.Errorf("Unable to get JSON response: %s", err)
		return
	}

	if err = char.Unmarshal(respJSON); err != nil {
		glog.Errorf("Unable to unmarshal character from JSON: %s", err)
		return
	}

	wg.Add(len(char.NewsFeed))
	for _, n := range char.NewsFeed {
		go n.updateNews(&wg)
	}
	wg.Wait()

	return
}

// GetRecentLegendaries check the character news feed and gets legendaries from it for a time period
func (char *Character) GetRecentLegendaries() (items []*Item) {
	var item *Item

	now := time.Now()
	before := now.Add(-o.LegendaryRelevancePeriod)

	for _, n := range char.NewsFeed {
		if inTimeSpan(before, now, n.EventTime) {
			if item = n.ItemInfo; item == nil {
				continue
			}

			isLegendary := item.Quality == 5 && item.Equippable && item.ItemLevel >= 910
			if isLegendary {
				items = append(items, item)
			}
		}
	}

	return
}

func (n *News) updateNews(wg *sync.WaitGroup) {
	var (
		utc *time.Location
		err error
	)

	if utc, err = time.LoadLocation(o.GuildTimezone); err != nil {
		glog.Errorf("Unable to parse location: %s", err)
		return
	}

	eventTime := time.Unix(n.Timestamp/1000, 0)
	n.EventTime = eventTime.In(utc)

	if n.Type == "itemLoot" || n.Type == "LOOT" {
		item, err := getItemByID(strconv.Itoa(n.ItemID))
		if err != nil {
			glog.Errorf("Unable to get item by its ID = %d: %s", n.ItemID, err)
			return
		}
		n.ItemInfo = item
	}

	wg.Done()

	return
}
