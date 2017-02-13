package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"
)

var refillMtx sync.Mutex

func getGuildMembers(realmName, guildName string, params []string) (*MembersList, error) {
	gInfo, cached, err := getGuildInfo(&realmName, &guildName)
	if err != nil {
		return nil, err
	}
	gMembers := gInfo.GuildMembersList
	if !cached {
		glog.Infof("Got %d guild members from API. Filling the gaps...", len(gMembers))
		gMembers = gMembers.refillMembers()
		// glog.Info("Saving guild members into cache...")
		gInfo.GuildMembersList = gMembers
		giJSON, err := gInfo.marshal()
		if err != nil {
			glog.Errorf("Unable to marshal guild info: %s", err)
			return nil, err
		}
		err = Put("Main", o.Bucket, giJSON)
		if err != nil {
			glog.Errorf("Unable to save guild info in DB: %s", err)
			return nil, err
		}
	} else {
		glog.Infof("Got %d guild members from cache", len(gMembers))
	}
	gMembers = gMembers.SortGuildMembers(params)
	// glog.Info("Sorted guild members")
	return &gMembers, nil
}

func getGuildInfo(guildRealm, guildName *string) (gInfo GuildInfo, cached bool, err error) {
	// glog.Info("getting main guild members...")
	cached = true
	membersJSON := Get("Main", o.Bucket)
	if membersJSON == nil {
		// glog.Info("No cache is present, getting from API...")
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
	// glog.Info("getting additional guild members...")
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
				if member.Char.Name == character {
					*ml = append(*ml, member)
				}
			}
		}
	}
	return nil
}

func (ml *MembersList) refillMembers() (refilledMembers MembersList) {
	var wg sync.WaitGroup
	wg.Add(len(*ml))
	for _, member := range *ml {
		go func(gm GuildMember) {
			defer wg.Done()
			char := updateCharacter(&gm)
			appendMember(&char, &refilledMembers)
		}(member)
	}
	wg.Wait()
	// glog.Info("Members refilled")
	return
}

func updateCharacter(member *GuildMember) (gm GuildMember) {
	var (
		items *Items
		profs *Professions
		err   error
	)
	gm = *member
	gm.Char.Class = classes[gm.Char.ClassInt]
	gm.Char.Gender = genders[gm.Char.GenderInt]
	gm.Char.Race = races[gm.Char.RaceInt]
	gm.Char.RealmSlug, err = getRealmSlugByName(gm.Char.Realm)
	if err != nil {
		glog.Errorf("updateCharacter(): unable to get realm slug: %s", err)
		return
	}
	shortLink, err := getArmoryLink(gm.Char.RealmSlug, gm.Char.Name)
	if err != nil {
		glog.Errorf("updateCharacter(): unable to get Armory link: %s", err)
		return
	}
	gm.Char.Link = shortLink
	items, err = getCharacterItems(&gm.Char.Realm, &gm.Char.Name)
	if err != nil {
		glog.Errorf("updateCharacter(): unable to get items: %s", err)
		return
	}
	gm.Char.Items = *items
	profs, err = getCharacterProfessions(&gm.Char.Realm, &gm.Char.Name)
	if err != nil {
		glog.Errorf("updateCharacter(): unable to get profs: %s", err)
		return
	}
	gm.Char.Professions = *profs
	return
}

func appendMember(gm *GuildMember, ml *MembersList) {
	refillMtx.Lock()
	*ml = append(*ml, *gm)
	refillMtx.Unlock()
}
