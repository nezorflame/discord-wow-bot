package main

import (
	"sort"
	"strings"

	"github.com/golang/glog"
)

// Len - function which returns length
func (nl NewsList) Len() int {
	return len(nl)
}

// Less - function which returns if element i is less than j
func (nl NewsList) Less(i, j int) bool {
	return nl[i].Timestamp < nl[j].Timestamp
}

// Swap - function which swaps element i with j
func (nl NewsList) Swap(i, j int) {
	nl[i], nl[j] = nl[j], nl[i]
}

// SortGuildMembers - function for sorting the guild members by a slice of params
func (ml *MembersList) SortGuildMembers(params []string) MembersList {
	// glog.Infof("sorting guild members, count = %d", len(*ml))
	var gMembers MembersList
	gMembers = sortGuildMembersByString(*ml, "name", "asc")
	length := len(gMembers)
	for _, p := range params {
		switch p {
		case "top5":
			length = 5
			continue
		case "top10":
			length = 10
			continue
		default:
			s := strings.Split(p, "=")
			if len(s) < 2 {
				// glog.Infof("Parameter '%s' is bad! Ignoring...", p)
				continue
			}
			pName := s[0]
			sOrder := s[1]
			switch pName {
			case "name", "class", "spec":
				gMembers = sortGuildMembersByString(gMembers, pName, sOrder)
			case "level", "ilvl":
				gMembers = sortGuildMembersByInt(gMembers, pName, sOrder)
			default:
				// glog.Infof("Unknown parameter '%s', so skipping...", pName)
			}
		}
	}
	if len(params) == 0 || params[0] == "" || strings.HasPrefix(params[0], "top") {
		gMembers = sortGuildMembersByInt(gMembers, "level", "desc")
		gMembers = sortGuildMembersByInt(gMembers, "ilvl", "desc")
		// glog.Info("No sorting params, used only default sort order...")
	}
	return gMembers[:length]
}

func sortGuildMembersByString(ml MembersList, key, order string) MembersList {
	// glog.Infof("sorting guild members by string '%s' and order '%s'", key, order)
	gMembersMap := make(map[string]MembersList)
	var sortedMembers MembersList
	var keys []string
	ascOrder := true
	for _, gm := range ml {
		var mKey string
		switch key {
		case "name":
			mKey = gm.Char.Name
		case "class":
			mKey = gm.Char.Class
		case "spec":
			mKey = gm.Char.Spec.Name
		default:
			glog.Infof("Unknown key '%s'. Aborting...", key)
			return ml
		}
		members := gMembersMap[mKey]
		if !members.checkMSliceForMember(gm) {
			gMembersMap[mKey] = append(gMembersMap[mKey], gm)
			if !checkStrSliceForValue(keys, mKey) {
				keys = append(keys, mKey)
			}
		}
	}
	if order == "desc" {
		ascOrder = false
	}
	if ascOrder {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}
	for _, k := range keys {
		for _, member := range gMembersMap[k] {
			sortedMembers = append(sortedMembers, member)
		}
	}
	return sortedMembers
}

func sortGuildMembersByInt(ml MembersList, key, order string) MembersList {
	// glog.Infof("sorting guild members by int '%s' and order '%s'", key, order)
	gMembersMap := make(map[int]MembersList)
	var sortedMembers MembersList
	var keys []int
	ascOrder := true
	for _, gm := range ml {
		var k int
		switch key {
		case "level":
			k = gm.Char.Level
		case "ilvl":
			k = gm.Char.Items.AvgItemLvlEq
		default:
			glog.Infof("Unknown key '%s'. Aborting...", key)
			return ml
		}
		members := gMembersMap[k]
		if !members.checkMSliceForMember(gm) {
			gMembersMap[k] = append(gMembersMap[k], gm)
			if !checkIntSliceForValue(keys, k) {
				keys = append(keys, k)
			}
		}
	}
	if order == "desc" {
		ascOrder = false
	}
	if ascOrder {
		sort.Ints(keys)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	}
	for _, k := range keys {
		for _, member := range gMembersMap[k] {
			sortedMembers = append(sortedMembers, member)
		}
	}
	return sortedMembers
}

func (ml *MembersList) checkMSliceForMember(member GuildMember) bool {
	for _, gm := range *ml {
		if gm.Char.Name == member.Char.Name {
			return true
		}
	}
	return false
}

func (nl *NewsList) checkNSliceForNews(news News) bool {
	for _, n := range *nl {
		if n.Timestamp == news.Timestamp && n.Character == news.Character {
			return true
		}
	}
	return false
}

func checkIntSliceForValue(slice []int, value int) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}

func checkFloatSliceForValue(slice []float64, value float64) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}

func checkStrSliceForValue(slice []string, value string) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}
