package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"time"

	"github.com/golang/glog"
)

func (b *Bot) guildWatcher() {
	var (
		wg  sync.WaitGroup
		err error
	)

	for {
		b.CharMutex.Lock()

		b.HighLvlCharacters = make(map[string]*Character)

		if b.Guild, err = GetGuildInfo(); err != nil {
			glog.Errorf("Unable to get the guild info: %s", err)
			goto out
		}

		wg.Add(len(b.Guild.MembersList))
		for _, member := range b.Guild.MembersList {
			go member.Char.UpdateCharacter(&wg)
			if member.Char.Level > 100 {
				b.HighLvlCharacters[member.Char.Name] = member.Char
			}
		}
		wg.Wait()
		glog.Info("Characters imported")

	out:
		b.CharMutex.Unlock()
		time.Sleep(o.CharacterCheckPeriod)
	}
}

func (b *Bot) legendaryWatcher() {
	var (
		wg  sync.WaitGroup
		err error
	)

	// wait a bit for a guild watcher to launch
	time.Sleep(time.Second)

	for {
		b.CharMutex.Lock()

		wg.Add(len(b.HighLvlCharacters))
		for _, char := range b.HighLvlCharacters {
			go func(c *Character) {
				if err = c.SetCharacterNewsFeed(); err != nil {
					glog.Errorf("Unable to set news feed for a character %s: %s", c.Name, err)
					wg.Done()
					return
				}

				for _, l := range c.GetRecentLegendaries() {
					if !b.checkForLegendary(c.Name, l.ID) {
						b.LegendariesByChar[c.Name] = append(b.LegendariesByChar[c.Name], l)
						msg := fmt.Sprintf(m.Legendary, c.Name, l.Name, l.Link)
						b.SendMessage(o.GeneralChannelID, msg)
						// glog.Info(msg)
					}
				}

				wg.Done()
			}(char)
		}
		wg.Wait()
		glog.Info("Characters checked for legendaries")

		b.CharMutex.Unlock()
		time.Sleep(o.LegendaryCheckPeriod)
	}
}

func (b *Bot) checkForLegendary(charName string, itemID int) bool {
	for _, l := range b.LegendariesByChar[charName] {
		if l.ID == itemID {
			return true
		}
	}
	return false
}

// GetRealmStatus - function for receiving realm status
func GetRealmStatus(realmName string) (bool, error) {
	realms, err := getRealms()
	if err != nil {
		return false, err
	}

	for _, r := range realms.RealmList {
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

	for _, r := range realms.RealmList {
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

// GetGuildInfo - function for receiving the guild info
func GetGuildInfo() (gInfo *GuildInfo, err error) {
	var membersJSON []byte

	apiLink := fmt.Sprintf(o.APIGuildMembersLink, o.GuildRegion, strings.Replace(o.GuildRealm, " ", "%20", -1),
		strings.Replace(o.GuildName, " ", "%20", -1), o.GuildLocale, o.WoWToken)
	if membersJSON, err = GetJSONResponse(apiLink); err != nil {
		glog.Errorf("Unable to get guild info: %s", err)
		return
	}

	gInfo = new(GuildInfo)
	if err = gInfo.Unmarshal(membersJSON); err != nil {
		glog.Errorf("Unable to unmarshal guild info: %s", err)
	}

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
