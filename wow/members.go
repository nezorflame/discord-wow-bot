package wow

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nezorflame/discord-wow-bot/consts"
	"github.com/nezorflame/discord-wow-bot/db"
	"github.com/nezorflame/discord-wow-bot/net"
)

func getGuildMembers(realmName, guildName string, params []string) (*MembersList, error) {
	gInfo, cached, err := getGuildInfo(&realmName, &guildName)
	if err != nil {
		return nil, err
	}
	gMembers := gInfo.GuildMembersList
	if !cached {
		logInfo("Got", len(gMembers), "guild members from API. Filling the gaps...")
		gMembers = gMembers.refillMembers()
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

func (ml *MembersList) refillMembers() (guildMembers MembersList) {
	var wg sync.WaitGroup
	wg.Add(len(*ml))
	for _, m := range *ml {
		go func(m GuildMember) {
			defer wg.Done()
			gMember := updateCharacter(m)
			guildMembers = append(guildMembers, gMember)
		}(m)
	}
	logInfo("Members refilled")
	wg.Wait()
	return
}

func updateCharacter(member GuildMember) (m GuildMember) {
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
	items, err = getCharacterItems(&m.Member.Realm, &m.Member.Name)
	if err != nil {
		logInfo("updateCharacter(): unable to get items:", err)
		return member
	}
	m.Member.Items = *items
	profs, err = getCharacterProfessions(&m.Member.Realm, &m.Member.Name)
	if err != nil {
		logInfo("updateCharacter(): unable to get profs:", err)
		return member
	}
	m.Member.Professions = *profs
	return
}
