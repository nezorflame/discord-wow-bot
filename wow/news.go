package wow

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nezorflame/discord-wow-bot/consts"
	"github.com/nezorflame/discord-wow-bot/net"
)

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
	gInfo := new(GuildInfo)
	err = gInfo.unmarshal(&respJSON)
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
