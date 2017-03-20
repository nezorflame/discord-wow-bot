package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/golang/glog"
)

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
	realmInfo := fmt.Sprintf(m.RealmInfo, realm.Name, realm.Type, realm.Population, realm.Status,
		realm.Queue, realm.Battlegroup, realm.Locale, realm.Timezone, realm.ConnectedRealms)
	return realmInfo, nil
}

// GetRealmSlug - function for receiving realm slug
func GetRealmSlug(realmName string) (string, error) {
	realm, err := getRealmByName(realmName)
	if err != nil {
		return "", err
	}
	return realm.Slug, nil
}

// GetGuildMembers - function for receiving a list of guild members
func GetGuildMembers(realmName, guildName string, params []string) (members []map[string]string, err error) {
	return
}

// GetGuildProfs - function for receiving a list of guild professions
func GetGuildProfs(realmName, guildName string, param string) (profs []map[string]string, err error) {
	return
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

func getRealms() (*[]Realm, error) {
	apiLink := fmt.Sprintf(o.APIRealmsLink, o.GuildRegion, o.GuildLocale, o.WoWToken)
	respJSON, err := GetJSONResponse(apiLink)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	var realms Realms
	if err = realms.unmarshal(&respJSON); err != nil {
		return nil, err
	}
	return &realms.RealmList, nil
}

func getRealmByName(realmName string) (Realm, error) {
	if !strings.Contains(realmName, " ") {
		realmName = splitStringByCase(realmName)
	}
	glog.Infof("getRealmByName: %s", realmName)
	realms, err := getRealms()
	if err != nil {
		return *new(Realm), err
	}
	for _, r := range *realms {
		if strings.ToLower(r.Name) == strings.ToLower(realmName) ||
			strings.ToLower(r.Slug) == strings.ToLower(realmName) {
			return r, nil
		}
	}
	return *new(Realm), errors.New("No such realm is present")
}

func splitStringByCase(splitString string) (result string) {
	l := 0
	for s := splitString; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		if result == "" {
			result = s[:l]
		} else {
			result += " " + s[:l]
		}
	}
	return
}
