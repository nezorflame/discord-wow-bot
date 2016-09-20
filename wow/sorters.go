package wow

import (
    "sort"
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

func (ml *MembersList) sortGuildMembers(params []string) MembersList {
    logInfo("sorting guild members, params are:", params)
    for _, p := range params {
        // var asc bool
        switch p {

        }
    }
    gMembersNameMap := make(map[string]GuildMember)
    var gMembers MembersList
    var keys []string
    for _, m := range *ml {
        gMembersNameMap[m.Member.Name] = m
        keys = append(keys, m.Member.Name)
    }
    sort.Strings(keys)
    for _, k := range keys {
        gMembers = append(gMembers, gMembersNameMap[k])
        logInfo(gMembersNameMap[k])
    }
    return gMembers
}

func (ml *MembersList) sortGuildMembersByName() MembersList {
    logInfo("sorting guild members by name...")
    gMembersNameMap := make(map[string]GuildMember)
    var gMembers MembersList
    var keys []string
    for _, m := range *ml {
        gMembersNameMap[m.Member.Name] = m
        keys = append(keys, m.Member.Name)
    }
    sort.Strings(keys)
    for _, k := range keys {
        gMembers = append(gMembers, gMembersNameMap[k])
        logInfo(gMembersNameMap[k])
    }
    return gMembers
}
