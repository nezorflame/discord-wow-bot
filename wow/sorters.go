package wow

import "sort"

func sortGuildNewsByTimestamp(gNews *[]News) *[]News {
    gNewsTimeMap := make(map[float64]News)
    var keys []float64
    for _, n := range *gNews {
        gNewsTimeMap[float64(n.Timestamp)] = n
        keys = append(keys, float64(n.Timestamp))
    }
    sort.Float64s(keys)
    gNews = new([]News)
    for _, k := range keys {
        *gNews = append(*gNews, gNewsTimeMap[k])
    }
    return gNews
}

func sortGuildMembersByName(gMembers *[]GuildMember) *[]GuildMember {
    gMembersNameMap := make(map[string]GuildMember)
    var keys []string
    for _, m := range *gMembers {
        gMembersNameMap[m.Member.Name] = m
        keys = append(keys, m.Member.Name)
    }
    sort.Strings(keys)
    gMembers = new([]GuildMember)
    for _, k := range keys {
        *gMembers = append(*gMembers, gMembersNameMap[k])
    }
    return gMembers
}
