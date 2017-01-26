package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

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
	glog.Info("News refilled")
	wg.Wait()
	return
}

func updateNews(newsrecord News) (news News) {
	if newsrecord.Type == "itemLoot" {
		item, err := getItemByIDFromAPI(newsrecord.ItemID)
		if err != nil {
			glog.Info("updateCharacter(): unable to get item by its ID =", newsrecord.ItemID, ":", err)
			return newsrecord
		}
		newsrecord.ItemInfo = *item
	}
	return newsrecord
}
