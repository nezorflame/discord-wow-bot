package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"time"
)

func (b *Bot) guildWatcher() {
	var (
		wg  sync.WaitGroup
		err error
	)

	for {
		b.CharMutex.Lock()

		b.HighLvlCharacters = make(map[string]*Character)

		if err = b.GetGuildInfo(); err != nil {
			b.SL.Errorf("Unable to get the guild info: %s", err)
			goto out
		}

		wg.Add(len(b.Guild.MembersList))
		for _, member := range b.Guild.MembersList {
			member.Char.RLock()

			if member.Char.Name == "" {
				b.SL.Errorf("Faulty character: %v", member.Char)
				wg.Done()
				continue
			}

			if member.Char.Level > 100 {
				b.HighLvlCharacters[member.Char.Name] = member.Char
			}

			member.Char.RUnlock()

			go member.Char.UpdateCharacter(b.SL, &wg)
		}
		wg.Wait()
		b.SL.Info("Characters imported")

	out:
		b.CharMutex.Unlock()
		time.Sleep(o.CharacterCheckPeriod)
	}
}

func (b *Bot) legendaryWatcher() {
	var wg sync.WaitGroup

	for {
		b.CharMutex.Lock()

		wg.Add(len(b.HighLvlCharacters))
		for _, char := range b.HighLvlCharacters {
			go char.SetCharacterNewsFeed(b.SL, &wg)
		}
		wg.Wait()

		b.SL.Info("Characters updated with the latest news")

		for _, char := range b.HighLvlCharacters {
			legendaries := char.GetRecentLegendaries()

			char.RLock()
			for _, l := range legendaries {
				if !b.checkForLegendary(char.Name, l.ID) {
					b.LegendariesByChar[char.Name] = append(b.LegendariesByChar[char.Name], l)
					msg := fmt.Sprintf(m.Legendary, char.Name, l.Name, l.Link)
					b.SendMessage(o.GeneralChannelID, msg)
					// b.SL.Info(msg)
				}
			}
			char.RUnlock()
		}

		b.SL.Infof("Characters checked for legendaries, current item count: %d", len(WoWItemsMap))

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
func (b *Bot) GetGuildInfo() (err error) {
	var (
		membersJSON []byte
		membersList []*GuildMember
	)

	apiLink := fmt.Sprintf(o.APIGuildMembersLink, o.GuildRegion, strings.Replace(o.GuildRealm, " ", "%20", -1),
		strings.Replace(o.GuildName, " ", "%20", -1), o.GuildLocale, o.WoWToken)
	if membersJSON, err = Get(apiLink); err != nil {
		err = errors.Wrap(err, "Unable to get guild info")
		return
	}

	gInfo := new(GuildInfo)
	if err = gInfo.Unmarshal(membersJSON); err != nil {
		err = errors.Wrap(err, "Unable to unmarshal guild info")
	}

	for _, member := range gInfo.MembersList {
		if member.Char.Name != "" {
			membersList = append(membersList, member)
		}
	}
	gInfo.MembersList = membersList

	b.Guild = gInfo
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
