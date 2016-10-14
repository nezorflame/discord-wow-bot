package wow

import (
    "sort"
    "strings"
)

func (nl *NewsList) sortGuildNewsByTimestamp() NewsList {
    logInfo("sorting guild news by timestamp...")
    gNewsTimeMap := make(map[float64]News)
    var gNews NewsList
    var keys []float64
    for _, n := range *nl {
        gNewsTimeMap[float64(n.Timestamp)] = n
        keys = append(keys, float64(n.Timestamp))
    }
    sort.Float64s(keys)
    for _, k := range keys {
        gNews = append(gNews, gNewsTimeMap[k])
    }
    return gNews
}

// SortGuildMembers - function for sorting the guild members by a slice of params
func (ml *MembersList) SortGuildMembers(params []string) MembersList {
    logInfo("sorting guild members...")
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
                    logInfo("Parameter", p, "is bad! Ignoring...")
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
                        logInfo("Unknown parameter", pName, "so skipping...")
                }
        }
    }
    if len(params) == 0 || params[0] == "" || strings.HasPrefix(params[0], "top") {
        gMembers = sortGuildMembersByInt(gMembers, "level", "desc")
        gMembers = sortGuildMembersByInt(gMembers, "ilvl", "desc")
        logInfo("No params, using default sort order...")
    }
    return gMembers[:length]
}

func sortGuildMembersByString(ml MembersList, key, order string) MembersList {
    logInfo("sorting guild members by string:", key)
    gMembersMap := make(map[string]MembersList)
    sortedMembers := new(MembersList)
    var keys []string
    ascOrder := true
    for _, m := range ml {
        switch key {
            case "name":
                members := gMembersMap[m.Member.Name]
                if !members.checkMSliceForMember(m) {
                    gMembersMap[m.Member.Name] = append(gMembersMap[m.Member.Name], m)
                    if !checkStrSliceForValue(keys, m.Member.Name) {
                        keys = append(keys, m.Member.Name)
                    }
                }
            case "class":
                members := gMembersMap[m.Member.Class]
                if !members.checkMSliceForMember(m) {
                    gMembersMap[m.Member.Class] = append(gMembersMap[m.Member.Class], m)
                    if !checkStrSliceForValue(keys, m.Member.Class) {
                        keys = append(keys, m.Member.Class)
                    }
                }
            case "spec":
                members := gMembersMap[m.Member.Spec.Name]
                if !members.checkMSliceForMember(m) {
                    gMembersMap[m.Member.Spec.Name] = append(gMembersMap[m.Member.Spec.Name], m)
                    if !checkStrSliceForValue(keys, m.Member.Spec.Name) {
                        keys = append(keys, m.Member.Spec.Name)
                    }
                }
            default:
                logInfo("Unknown key: " + key + ". Aborting...")
                return ml
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
            *sortedMembers = append(*sortedMembers, member)
        }
    }
    gMembersMap = nil
    keys = nil
    return *sortedMembers
}

func sortGuildMembersByInt(ml MembersList, key, order string) MembersList {
    logInfo("sorting guild members by int:", key)
    gMembersMap := make(map[int]MembersList)
    sortedMembers := new(MembersList)
    var keys []int
    ascOrder := true
    for _, m := range ml {
        switch key {
            case "level":
                members := gMembersMap[m.Member.Level]
                if !members.checkMSliceForMember(m) {
                    gMembersMap[m.Member.Level] = append(gMembersMap[m.Member.Level], m)
                    if !checkIntSliceForValue(keys, m.Member.Level) {
                        keys = append(keys, m.Member.Level)
                    }
                }
            case "ilvl":
                members := gMembersMap[m.Member.Items.AvgItemLvlEq]
                if !members.checkMSliceForMember(m) {
                    gMembersMap[m.Member.Items.AvgItemLvlEq] = append(gMembersMap[m.Member.Items.AvgItemLvlEq], m)
                    if !checkIntSliceForValue(keys, m.Member.Items.AvgItemLvlEq) {
                        keys = append(keys, m.Member.Items.AvgItemLvlEq)
                    }
                }
            default:
                logInfo("Unknown key: " + key + ". Aborting...")
                return ml
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
            *sortedMembers = append(*sortedMembers, member)
        }
    }
    return *sortedMembers
}

func (ml *MembersList) checkMSliceForMember(member GuildMember) bool {
    for _, m := range *ml {
        if m.Member.Name == member.Member.Name {
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

func checkStrSliceForValue(slice []string, value string) bool {
    for _, i := range slice {
        if i == value {
            return true
        }
    }
    return false
}
