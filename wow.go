package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	classes    map[int]string
	factions   map[int]string
	races      map[int]string
	genders    map[int]string
	profNames  map[int]string
	addMembers map[string]map[string]string
)

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// InitializeWoWAPI - function for initializing WoW API
func InitializeWoWAPI() {
	// TODO: Rework to WoW API structures from response
	classes = map[int]string{
		1:  "Воин",
		2:  "Паладин",
		3:  "Охотник",
		4:  "Разбойник",
		5:  "Жрец",
		6:  "Рыцарь смерти",
		7:  "Шаман",
		8:  "Маг",
		9:  "Чернокнижник",
		10: "Монах",
		11: "Друид",
		12: "Охотник на демонов",
	}
	genders = map[int]string{
		0: "Мужчина",
		1: "Женщина",
	}
	factions = map[int]string{
		0: "Альянс",
		1: "Орда",
	}
	races = map[int]string{
		1:  "Человек",
		2:  "Орк",
		3:  "Дворф",
		4:  "Ночной эльф",
		5:  "Нежить",
		6:  "Таурен",
		7:  "Гном",
		8:  "Тролль",
		9:  "Гоблин",
		10: "Эльф крови",
		11: "Дреней",
		22: "Ворген",
		24: "Пандарен",
		25: "Пандарен",
		26: "Пандарен",
	}
	profNames = map[int]string{
		171: "alchemy",
		164: "blacksmithing",
		794: "archaeology",
		185: "cooking",
		333: "enchanting",
		202: "engineering",
		129: "first-aid",
		356: "fishing",
		182: "herbalism",
		773: "inscription",
		755: "jewelcrafting",
		165: "leatherworking",
		186: "mining",
		393: "skinning",
		197: "tailoring",
	}
	addMembers = make(map[string]map[string]string)
}

// GetRealmStatus - function for receiving realm status
func GetRealmStatus(realmName string) (bool, error) {
	realms, err := getRealms()
	if err != nil {
		return false, err
	}
	for _, r := range *realms {
		if r.Name == realmName || r.Slug == realmName {
			return r.Status, nil
		}
	}
	return false, errors.New("No such realm is present")
}

// GetRealmQueueStatus - function for receiving realm queue status
func GetRealmQueueStatus(realmName string) (bool, error) {
	realms, err := getRealms()
	if err != nil {
		return false, err
	}
	for _, r := range *realms {
		if r.Name == realmName || r.Slug == realmName {
			return r.Queue, nil
		}
	}
	return false, errors.New("No such realm is present")
}

// GetRealmInfo - function for receiving realm info
func GetRealmInfo(realmName string) (string, error) {
	realm, err := getRealmByName(realmName)
	if err != nil {
		return "", err
	}
	realmInfo := "Имя сервера: %v\n"
	realmInfo += "Тип сервера: %v\n"
	realmInfo += "Населенность: %v\n"
	realmInfo += "Статус: %t\n"
	realmInfo += "Очередь на вход: %t\n"
	realmInfo += "PvP-группа: %v\n"
	realmInfo += "Язык: %v\n"
	realmInfo += "Временной пояс: %v\n"
	realmInfo += "Связанные серверы: %v"
	realmInfo = fmt.Sprintf(realmInfo, realm.Name, realm.Type, realm.Population, realm.Status,
		realm.Queue, realm.Battlegroup, realm.Locale, realm.Timezone, realm.ConnectedRealms)
	return realmInfo, nil
}

// GetGuildLegendaries - function for getting the latest guild legendaries
func GetGuildLegendaries(realmName, guildName string) ([]string, error) {
	var legendaries, params []string
	gMembers, err := getGuildMembers(realmName, guildName, params)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(len(*gMembers))
	for _, member := range *gMembers {
		go func(name string, lvl int) {
			defer wg.Done()
			if lvl < 110 {
				return
			}
			cNews, err := getCharNews(realmName, name)
			if err != nil {
				return
			}
			for _, n := range *cNews {
				if n.Type != "LOOT" && n.Type != "itemLoot" {
					continue
				}
				item := n.ItemInfo
				isLegendary := item.Quality == 5 && item.ItemLevel >= 910
				if isLegendary {
					message := fmt.Sprintf(m.Legendary, name, item.Name, item.Link)
					legendaries = append(legendaries, message)
				}
			}
		}(member.Char.Name, member.Char.Level)
	}
	wg.Wait()
	return legendaries, nil
}

// GetGuildMembers - function for receiving a list of guild members
func GetGuildMembers(realmName, guildName string, params []string) ([]map[string]string, error) {
	gMembers, err := getGuildMembers(realmName, guildName, params)
	if err != nil {
		return nil, err
	}
	var guildMembersList []map[string]string
	for _, gm := range *gMembers {
		gMember := make(map[string]string)
		gMember["Name"] = gm.Char.Name
		gMember["Guild"] = gm.Char.Guild
		gMember["Realm"] = gm.Char.GuildRealm
		gMember["Level"] = strconv.Itoa(gm.Char.Level)
		gMember["Class"] = gm.Char.Class
		if specName := gm.Char.Spec.Name; specName != "" {
			gMember["Spec"] = specName
		} else {
			gMember["Spec"] = "Нет инфы"
		}
		gMember["ItemLevel"] = strconv.Itoa(gm.Char.Items.AvgItemLvlEq)
		gMember["Link"] = gm.Char.Link
		guildMembersList = append(guildMembersList, gMember)
	}
	return guildMembersList, nil
}

// GetGuildProfs - function for receiving a list of guild professions
func GetGuildProfs(realmName, guildName string, param string) ([]map[string]string, error) {
	params := []string{"name=asc"}
	gMembers, err := getGuildMembers(realmName, guildName, params)
	if err != nil {
		return nil, err
	}
	var profName string
	if param != "" {
		s := strings.Split(param, "=")
		if len(s) < 2 {
			return nil, errors.New("Не указана желаемая профессия, повтори ввод")
		}
		profName = s[1]
	}
	var guildProfsList []map[string]string
	for _, gm := range *gMembers {
		gMember := make(map[string]string)
		gMember["Name"] = gm.Char.Name
		switch len(gm.Char.Professions.PrimaryProfs) {
		case 0:
			gMember["FirstProf"] = "Нет"
			gMember["FirstProfLevel"] = "-"
			gMember["SecondProf"] = "Нет"
			gMember["SecondProfLevel"] = "-"
		case 1:
			gMember["FirstProf"] = gm.Char.Professions.PrimaryProfs[0].Name
			gMember["FirstProfLevel"] = strconv.Itoa(gm.Char.Professions.PrimaryProfs[0].Rank) +
				" | " + gm.Char.Professions.PrimaryProfs[0].Link
			gMember["SecondProf"] = "Нет"
			gMember["SecondProfLevel"] = "-"
		case 2:
			gMember["FirstProf"] = gm.Char.Professions.PrimaryProfs[0].Name
			gMember["FirstProfLevel"] = strconv.Itoa(gm.Char.Professions.PrimaryProfs[0].Rank) +
				" | " + gm.Char.Professions.PrimaryProfs[0].Link
			gMember["SecondProf"] = gm.Char.Professions.PrimaryProfs[1].Name
			gMember["SecondProfLevel"] = strconv.Itoa(gm.Char.Professions.PrimaryProfs[1].Rank) +
				" | " + gm.Char.Professions.PrimaryProfs[1].Link
		}
		if profName == "" || gMember["FirstProf"] == profName || gMember["SecondProf"] == profName {
			guildProfsList = append(guildProfsList, gMember)
		}
	}
	if len(guildProfsList) == 0 {
		return nil, errors.New("Такой профессии ни у кого нет, или она введена неверно")
	}
	return guildProfsList, nil
}

// GetRealmName returns realm name string
func GetRealmName(message string, command string) string {
	commandString := strings.Replace(message, command, "", 1)
	if commandString == "" {
		return o.GuildRealm
	}
	return strings.TrimLeft(commandString, " ")
}

// GetRealmAndGuildNames returns realm and guild name strings
func GetRealmAndGuildNames(message string, command string) (string, string, error) {
	commandString := strings.Replace(message, command, "", 1)
	if commandString == "" {
		return o.GuildRealm, o.GuildName, nil
	}
	s := strings.Split(commandString, ", ")
	if len(s) < 2 {
		return "", "", errors.New("Команда введена неверно, попробуй еще раз")
	}
	return s[0], s[1], nil
}

// GetDefaultRealmAndGuildNames returns default realm and guild name strings
func GetDefaultRealmAndGuildNames() (string, string) {
	return o.GuildRealm, o.GuildName
}
