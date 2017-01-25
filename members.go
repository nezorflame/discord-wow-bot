package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"
)

func getGuildMembers(realmName, guildName string, params []string) (*MembersList, error) {
	gInfo, cached, err := getGuildInfo(&realmName, &guildName)
	if err != nil {
		return nil, err
	}
	gMembers := gInfo.GuildMembersList
	if !cached {
		glog.Info("Got", len(gMembers), "guild members from API. Filling the gaps...")
		gMembers = gMembers.refillMembers()
		glog.Info("Saving guild members into cache...")
		gInfo.GuildMembersList = gMembers
		giJSON, err := gInfo.marshal()
		err = Put("Main", o.Bucket, giJSON)
		glog.Error(err)
	} else {
		glog.Info("Got", len(gMembers), "guild members from cache")
	}
	gMembers = gMembers.SortGuildMembers(params)
	glog.Info("Sorted guild members")
	return &gMembers, nil
}

func getGuildInfo(guildRealm, guildName *string) (gInfo GuildInfo, cached bool, err error) {
	glog.Info("getting main guild members...")
	cached = true
	membersJSON := Get("Main", o.Bucket)
	if membersJSON == nil {
		glog.Info("No cache is present, getting from API...")
		cached = false
		apiLink := fmt.Sprintf(o.APIGuildMembersLink, o.GuildRegion, strings.Replace(*guildRealm, " ", "%20", -1),
			strings.Replace(*guildName, " ", "%20", -1), o.GuildLocale, o.WoWToken)
		membersJSON, err = GetJSONResponse(apiLink)
		if err != nil {
			glog.Error(err)
			return
		}
	}
	gi := new(GuildInfo)
	err = gi.unmarshal(&membersJSON)
	if err != nil {
		return
	}
	err = gi.GuildMembersList.getAdditionalMembers()
	gInfo = *gi
	return
}

func (ml *MembersList) getAdditionalMembers() error {
	glog.Info("getting additional guild members...")
	for realm, m := range addMembers {
		for guild, character := range m {
			apiLink := fmt.Sprintf(o.APIGuildMembersLink, o.GuildRegion, strings.Replace(realm, " ", "%20", -1),
				strings.Replace(guild, " ", "%20", -1), o.GuildLocale, o.WoWToken)
			respJSON, err := GetJSONResponse(apiLink)
			if err != nil {
				glog.Error(err)
				return err
			}
			gi := new(GuildInfo)
			err = gi.unmarshal(&respJSON)
			if err != nil {
				return err
			}
			for _, member := range gi.GuildMembersList {
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
	glog.Info("Members refilled")
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
		glog.Info("updateCharacter(): unable to get realm slug:", err)
		return member
	}
	shortLink, err := getArmoryLink(m.Member.RealmSlug, m.Member.Name)
	if err != nil {
		glog.Info("updateCharacter(): unable to get Armory link:", err)
		return member
	}
	m.Member.Link = shortLink
	items, err = getCharacterItems(&m.Member.Realm, &m.Member.Name)
	if err != nil {
		glog.Info("updateCharacter(): unable to get items:", err)
		return member
	}
	m.Member.Items = *items
	profs, err = getCharacterProfessions(&m.Member.Realm, &m.Member.Name)
	if err != nil {
		glog.Info("updateCharacter(): unable to get profs:", err)
		return member
	}
	m.Member.Professions = *profs
	return
}
