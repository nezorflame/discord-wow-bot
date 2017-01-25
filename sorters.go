package main

import (
	"sort"
	"strings"

	"github.com/golang/glog"
)

// SortGuildNews - function for sorting the guild news by timestamp
func (nl *NewsList) SortGuildNews() NewsList {
	glog.Info("sorting guild news by timestamp...")
	gNewsTimeMap := make(map[float64]NewsList)
	sortedNews := new(NewsList)
	var keys []float64
	for _, n := range *nl {
		k := float64(n.Timestamp)
		news := gNewsTimeMap[k]
		if !news.checkNSliceForNews(n) {
			gNewsTimeMap[k] = append(gNewsTimeMap[k], n)
			if !checkFloatSliceForValue(keys, k) {
				keys = append(keys, k)
			}
		}
	}
	sort.Float64s(keys)
	for _, k := range keys {
		for _, n := range gNewsTimeMap[k] {
			*sortedNews = append(*sortedNews, n)
			glog.Infof("%s %s %s %d", n.Character, n.ItemInfo.Name, n.ItemInfo.Link, n.ItemInfo.Quality)
		}
	}
	return *sortedNews
}

// SortGuildMembers - function for sorting the guild members by a slice of params
func (ml *MembersList) SortGuildMembers(params []string) MembersList {
	glog.Info("sorting guild members, count =", len(*ml))
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
				glog.Info("Parameter", p, "is bad! Ignoring...")
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
				glog.Info("Unknown parameter", pName, "so skipping...")
			}
		}
	}
	if len(params) == 0 || params[0] == "" || strings.HasPrefix(params[0], "top") {
		gMembers = sortGuildMembersByInt(gMembers, "level", "desc")
		gMembers = sortGuildMembersByInt(gMembers, "ilvl", "desc")
		glog.Info("No sorting params, used only default sort order...")
	}
	return gMembers[:length]
}

func sortGuildMembersByString(ml MembersList, key, order string) MembersList {
	glog.Info("sorting guild members by string:", key, "and order:", order)
	gMembersMap := make(map[string]MembersList)
	var sortedMembers MembersList
	var keys []string
	ascOrder := true
	for _, m := range ml {
		var mKey string
		switch key {
		case "name":
			mKey = m.Member.Name
		case "class":
			mKey = m.Member.Class
		case "spec":
			mKey = m.Member.Spec.Name
		default:
			glog.Info("Unknown key: " + key + ". Aborting...")
			return ml
		}
		members := gMembersMap[mKey]
		if !members.checkMSliceForMember(m) {
			gMembersMap[mKey] = append(gMembersMap[mKey], m)
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
	glog.Info("sorting guild members by int:", key, "and order:", order)
	gMembersMap := make(map[int]MembersList)
	var sortedMembers MembersList
	var keys []int
	ascOrder := true
	for _, m := range ml {
		var k int
		switch key {
		case "level":
			k = m.Member.Level
		case "ilvl":
			k = m.Member.Items.AvgItemLvlEq
		default:
			glog.Info("Unknown key: " + key + ". Aborting...")
			return ml
		}
		members := gMembersMap[k]
		if !members.checkMSliceForMember(m) {
			gMembersMap[k] = append(gMembersMap[k], m)
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
	for _, m := range *ml {
		if m.Member.Name == member.Member.Name {
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