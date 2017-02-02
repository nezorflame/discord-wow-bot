package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

func getCharNews(realmName, charName string) (*NewsList, error) {
	feed, err := getCharFeed(&realmName, &charName)
	if err != nil {
		return nil, err
	}
	feed = feed.refillNews()
	sort.Sort(feed)
	// glog.Info("Got updated character news")
	return &feed, nil
}

func getCharFeed(charRealm, charName *string) (feed NewsList, err error) {
	// glog.Info("getting character feed...")
	apiLink := fmt.Sprintf(o.APICharNewsLink, o.GuildRegion, strings.Replace(*charRealm, " ", "%20", -1),
		strings.Replace(*charName, " ", "%20", -1), o.GuildLocale, o.WoWToken)
	respJSON, err := GetJSONResponse(apiLink)
	if err != nil {
		return
	}
	char := new(Character)
	err = char.unmarshal(&respJSON)
	if err != nil {
		return
	}
	now := time.Now()
	before := now.Add(time.Duration(-6 * time.Minute))
	// Fill string valuables
	char.Faction = factions[char.FactionInt]
	for _, n := range char.Feed {
		var (
			utc *time.Location
			err error
		)
		eventTime := time.Unix(n.Timestamp/1000, 0)
		if utc, err = time.LoadLocation(o.GuildTimezone); err != nil {
			glog.Error(err)
		}
		n.EventTime = eventTime.In(utc)
		if inTimeSpan(before, now, eventTime) {
			feed = append(feed, n)
		}
	}
	return
}

func getGuildNews(realmName, guildName string) (*NewsList, error) {
	gNews, err := getGuildNewsList(&realmName, &guildName)
	if err != nil {
		return nil, err
	}
	gNews = gNews.refillNews()
	sort.Sort(gNews)
	glog.Info("Got updated guild news")
	return &gNews, nil
}

func getGuildNewsList(guildRealm, guildName *string) (gNews NewsList, err error) {
	glog.Info("getting guild news...")
	apiLink := fmt.Sprintf(o.APIGuildNewsLink, o.GuildRegion, strings.Replace(*guildRealm, " ", "%20", -1),
		strings.Replace(*guildName, " ", "%20", -1), o.GuildLocale, o.WoWToken)
	respJSON, err := GetJSONResponse(apiLink)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	gInfo := new(GuildInfo)
	err = gInfo.unmarshal(&respJSON)
	if err != nil {
		return
	}
	now := time.Now()
	before := now.Add(time.Duration(-6 * 240 * time.Minute))
	// Fill string valuables
	gInfo.Side = factions[gInfo.SideInt]
	for _, n := range gInfo.GuildNewsList {
		var (
			utc *time.Location
			err error
		)
		eventTime := time.Unix(n.Timestamp/1000, 0)
		if utc, err = time.LoadLocation(o.GuildTimezone); err != nil {
			glog.Error(err)
		}
		n.EventTime = eventTime.In(utc)
		if inTimeSpan(before, now, eventTime) {
			gNews = append(gNews, n)
		}
	}
	return
}

func (nl *NewsList) refillNews() (refilledNews NewsList) {
	var wg sync.WaitGroup
	wg.Add(len(*nl))
	for _, newsrecord := range *nl {
		go func(n News) {
			defer wg.Done()
			news := updateNews(&n)
			refilledNews = append(refilledNews, news)
		}(newsrecord)
	}
	wg.Wait()
	// glog.Info("News refilled")
	return
}

func updateNews(newsrecord *News) (n News) {
	n = *newsrecord
	if n.Type == "itemLoot" || n.Type == "LOOT" {
		item, err := getItemByIDFromAPI(n.ItemID)
		if err != nil {
			glog.Errorf("updateCharacter(): unable to get item by its ID = %d: %s", n.ItemID, err)
			return
		}
		n.ItemInfo = *item
	}
	return
}
