package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"
)

// WoW item vars
var (
	WoWItemsMap map[string]*Item
	WoWItemsMtx sync.RWMutex
)

func getRealms() (realms Realms, err error) {
	var respJSON []byte

	apiLink := fmt.Sprintf(o.APIRealmsLink, o.GuildRegion, o.GuildLocale, o.WoWToken)
	if respJSON, err = Get(apiLink); err != nil {
		glog.Errorf("Unable to get JSON response: %s", err)
		return
	}

	if err = realms.Unmarshal(respJSON); err != nil {
		glog.Errorf("Unable to unmarshal realms from JSON: %s", err)
	}

	return
}

func getRealmByName(realmName string) (realm Realm, err error) {
	var realms Realms

	if !strings.Contains(realmName, " ") {
		realmName = splitStringByCase(realmName)
	}

	if realms, err = getRealms(); err != nil {
		return
	}

	for _, r := range realms.RealmList {
		if strings.ToLower(r.Name) == strings.ToLower(realmName) ||
			strings.ToLower(r.Slug) == strings.ToLower(realmName) {
			realm = r
			return
		}
	}

	err = errors.New("No such realm is present")
	return
}

func getItemByID(itemID string) (item *Item, err error) {
	var (
		respJSON []byte
		ok       bool
	)

	WoWItemsMtx.RLock()
	if item, ok = WoWItemsMap[itemID]; ok {
		WoWItemsMtx.RUnlock()
		return
	}
	WoWItemsMtx.RUnlock()

	apiLink := fmt.Sprintf(o.APIItemLink, o.GuildRegion, itemID, o.GuildLocale, o.WoWToken)
	if respJSON, err = Get(apiLink); err != nil {
		glog.Errorf("Unable to get JSON response: %s", err)
		return
	}

	item = new(Item)
	if err = item.Unmarshal(respJSON); err != nil {
		glog.Info(err)
		return
	}

	item.Link = fmt.Sprintf(o.WowheadItemLink, itemID)

	WoWItemsMtx.Lock()
	WoWItemsMap[itemID] = item
	WoWItemsMtx.Unlock()

	return
}
