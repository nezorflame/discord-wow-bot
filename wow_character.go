package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// UpdateCharacter gets the additional info for the WoW character
func (char *Character) UpdateCharacter(sl *zap.SugaredLogger, wg *sync.WaitGroup) {
	var err error

	char.Lock()
	char.Class = o.WoWClasses[char.ClassInt]
	char.Gender = o.WoWGenders[char.GenderInt]
	char.Race = o.WoWRaces[char.RaceInt]
	char.Unlock()

	if err = char.SetRealmSlugByName(); err != nil {
		sl.Errorf("Unable to get realm slug for character %s: %s", char.Name, err)
	}

	if err = char.SetArmoryLink(sl); err != nil {
		sl.Errorf("Unable to get Armory link for character %s: %s", char.Name, err)
	}

	if err = char.SetCharacterItems(); err != nil {
		sl.Errorf("Unable to get items for character %s: %s", char.Name, err)
	}

	if err = char.SetCharacterProfessions(sl); err != nil {
		sl.Errorf("Unable to get profs for character %s: %s", char.Name, err)
	}

	wg.Done()

	return
}

// SetRealmSlugByName gets the realm slug name and sets it into the character
func (char *Character) SetRealmSlugByName() (err error) {
	var realms Realms

	char.Lock()
	defer char.Unlock()

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
func (char *Character) SetArmoryLink(sl *zap.SugaredLogger) (err error) {
	var shortLink string

	char.Lock()
	defer char.Unlock()

	longLink := fmt.Sprintf(o.ArmoryCharLink, o.GuildRegion, o.GuildLocale[:2], char.RealmSlug, char.Name)
	if shortLink, err = GetShortLink(longLink); err != nil {
		sl.Errorf("Unable to get short link for a character: %s", err)
		return
	}

	char.Link = shortLink

	return
}

// SetCharacterItems gets the character items and sets them into the character
func (char *Character) SetCharacterItems() (err error) {
	var respJSON []byte

	char.Lock()
	defer char.Unlock()

	apiLink := fmt.Sprintf(o.APICharItemsLink, o.GuildRegion, strings.Replace(char.Realm, " ", "%20", -1),
		char.Name, o.GuildLocale, o.WoWToken)
	if respJSON, err = Get(apiLink); err != nil {
		return
	}

	err = char.Unmarshal(respJSON)

	return
}

// SetCharacterProfessions gets the character professions and sets them into the character
func (char *Character) SetCharacterProfessions(sl *zap.SugaredLogger) (err error) {
	var respJSON []byte

	char.Lock()
	defer char.Unlock()

	apiLink := fmt.Sprintf(o.APICharProfsLink, o.GuildRegion, strings.Replace(char.Realm, " ", "%20", -1),
		char.Name, o.GuildLocale, o.WoWToken)
	if respJSON, err = Get(apiLink); err != nil {
		return
	}

	if err = char.Unmarshal(respJSON); err != nil {
		return
	}

	for _, p := range char.Professions.Primary {
		p.EngName = o.WoWProfessions[p.ID]
		longLink := fmt.Sprintf(o.ArmoryProfLink, o.GuildRegion, o.GuildLocale[:2], char.RealmSlug, char.Name, p.EngName)

		if shortLink, pErr := GetShortLink(longLink); pErr != nil {
			sl.Errorf("Unable to get short link for profession: %s", pErr)
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
func (char *Character) SetCharacterNewsFeed(sl *zap.SugaredLogger, mainWG *sync.WaitGroup) {
	var (
		respJSON []byte

		wg  sync.WaitGroup
		err error
	)

	defer mainWG.Done()

	char.Lock()
	defer char.Unlock()

	apiLink := fmt.Sprintf(o.APICharNewsLink, o.GuildRegion, strings.Replace(char.Realm, " ", "%20", -1),
		strings.Replace(char.Name, " ", "%20", -1), o.GuildLocale, o.WoWToken)
	if respJSON, err = Get(apiLink); err != nil {
		sl.Errorf("Unable to get JSON response: %s", err)
		return
	}

	if err = char.Unmarshal(respJSON); err != nil {
		sl.Errorf("Unable to unmarshal character from JSON: %s", err)
		return
	}

	wg.Add(len(char.NewsFeed))
	for _, n := range char.NewsFeed {
		go n.updateNews(sl, &wg)
	}
	wg.Wait()
}

// GetRecentLegendaries check the character news feed and gets legendaries from it for a time period
func (char *Character) GetRecentLegendaries() (items []*Item) {
	var item *Item

	char.RLock()

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

	char.RUnlock()

	return
}

func (n *News) updateNews(sl *zap.SugaredLogger, wg *sync.WaitGroup) {
	var (
		utc       *time.Location
		eventTime time.Time
		err       error
	)

	n.Lock()

	if utc, err = time.LoadLocation(o.GuildTimezone); err != nil {
		sl.Errorf("Unable to parse location: %s", err)
		goto out
	}

	eventTime = time.Unix(n.Timestamp/1000, 0)
	n.EventTime = eventTime.In(utc)

	if n.Type == "itemLoot" || n.Type == "LOOT" {
		item, err := getItemByID(strconv.Itoa(n.ItemID))
		if err != nil {
			sl.Errorf("Unable to get item by its ID = %d: %s", n.ItemID, err)
			goto out
		}
		n.ItemInfo = item
	}

out:
	n.Unlock()
	wg.Done()
	return
}
