package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arteev/fmttab"
	"github.com/bwmarrin/discordgo"
	"github.com/fiam/gounidecode/unidecode"
)

func (b *Bot) pingReporter(mes *discordgo.MessageCreate) {
	b.SL.Info("Sending pong to user...")
	b.SendMessage(mes.ChannelID, m.Pong)
}

func (b *Bot) helpReporter(mes *discordgo.MessageCreate) {
	b.SL.Info("Sending help to user...")
	b.SendMessage(mes.ChannelID, m.Help)
}

func (b *Bot) boobsReporter(mes *discordgo.MessageCreate) {
	b.SL.Info("Sending boobies to user...:)")
	b.SendMessage(mes.ChannelID, m.BoobiesPrefix)
	b.SendMessage(mes.ChannelID, m.Boobies[rand.Intn(len(m.Boobies))])
}

func (b *Bot) jcReporter(mes *discordgo.MessageCreate) {
	b.SL.Info("And his name is...")
	b.SendMessage(mes.ChannelID, m.JohnCena)
}

func (b *Bot) logReporter(mes *discordgo.MessageCreate) {
	b.SendMessage(mes.ChannelID, m.Logs)
}

func (b *Bot) rollReporter(mes *discordgo.MessageCreate) {
	b.SL.Info("Rolling a dice...")
	roll := rand.Intn(100) + 1
	var message string
	switch roll {
	case 1:
		message = fmt.Sprintf(m.Roll1, mes.Author.ID)
	case 100:
		message = fmt.Sprintf(m.Roll100, mes.Author.ID)
	default:
		message = fmt.Sprintf(m.RollX, mes.Author.ID, roll)
	}
	b.SendMessage(mes.ChannelID, message)
}

func (b *Bot) statusReporter(mes *discordgo.MessageCreate) {
	b.SL.Info("getting realm name string...")
	realmString := GetRealmName(mes.Content, "!status")
	b.SL.Infof("getting status of %s and sending it...", realmString)
	if realmStatus, err := GetRealmStatus(realmString); err != nil {
		b.SL.Errorf("Unable to get the realm status: %s", err)
		b.SendMessage(mes.ChannelID, m.ErrorUser)
	} else if realmStatus {
		b.SendMessage(mes.ChannelID, m.RealmOn)
	} else {
		b.SendMessage(mes.ChannelID, m.RealmOff)
	}
}

func (b *Bot) simcArmoryReporter(mes *discordgo.MessageCreate, command string, withStats, forPtr bool) {
	const timeFormat = "20060102_150405"
	var (
		simcExt = ".simc"
		htmlExt = ".html"

		argString, char, realm, region, output, profileName string

		args []string
		file *os.File
		err  error
	)
	b.SL.Info("getting simcraft sim...")
	params := strings.Split(strings.Replace(mes.Content, command+" ", "", 1), ",")
	cmdType := params[0]
	if cmdType == command || cmdType == "" {
		b.SL.Infof("Command is incorrect: %s", mes.Content)
		b.SendMessage(mes.ChannelID, m.ErrorUser)
		return
	}

	switch len(params) {
	case 1:
		region = o.GuildRegion
		if strings.Contains(params[0], "-") {
			charParams := strings.Split(params[0], "-")
			char = charParams[0]
			realm, err = GetRealmSlug(charParams[1])
			if err != nil {
				b.SL.Infof("Realm name is incorrect: %s", charParams[1])
				b.SendMessage(mes.ChannelID, m.ErrorUser)
				return
			}
		} else {
			realm = strings.Replace(o.GuildRealm, " ", "%20", -1)
			char = params[0]
		}
	case 2:
		region = o.GuildRegion
		realm, err = GetRealmSlug(params[0])
		if err != nil {
			b.SL.Infof("Realm name is incorrect: %s", params[0])
			b.SendMessage(mes.ChannelID, m.ErrorUser)
			return
		}
		char = params[1]
	case 3:
		region = params[0]
		realm, err = GetRealmSlug(params[1])
		if err != nil {
			b.SL.Infof("Realm name is incorrect: %s", params[1])
			b.SendMessage(mes.ChannelID, m.ErrorUser)
			return
		}
		char = params[2]
	default:
		b.SL.Infof("Command is incorrect: %s", mes.Content)
		b.SendMessage(mes.ChannelID, m.ErrorUser)
		return
	}

	location, _ := time.LoadLocation(o.GuildTimezone)
	now := time.Now().In(location).Format(timeFormat)
	profileName = fmt.Sprintf("%s_%s", char, now)
	profileFilePath := fmt.Sprintf("/tmp/%s%s", profileName, simcExt)
	resultsFileName := fmt.Sprintf("%s%s", profileName, htmlExt)
	resultsFilePath := "/tmp/" + resultsFileName
	argString = fmt.Sprintf(o.SimcArgsImport, region, realm, char, profileFilePath)
	args = strings.Split(argString, "|")

	b.SendMessage(mes.ChannelID, fmt.Sprintf(m.SimcArmory, char))

	if forPtr {
		command = o.SimcCmdPtr
	} else {
		command = o.SimcCmdStable
	}

	output, err = ExecuteCommand(command, o.SimcDir, args)
	// b.SL.Info(output)
	if err != nil {
		b.SL.Error(err)
		b.SendMessage(mes.ChannelID, m.SimcArmoryError)
		return
	}
	b.SL.Info("Created the user profile from Armory")
	b.SendMessage(mes.ChannelID, m.SimcImportSuccess)

	if withStats {
		argString = fmt.Sprintf(o.SimcArgsWithStats, profileFilePath, resultsFilePath)
	} else {
		argString = fmt.Sprintf(o.SimcArgsNoStats, profileFilePath, resultsFilePath)
	}
	args = strings.Split(argString, "|")

	if output, err = ExecuteCommand(command, o.SimcDir, args); err != nil {
		if strings.Contains(output, "Character not found") {
			b.SL.Error("Unable to find the character")
			b.SendMessage(mes.ChannelID, m.SimcArmoryError)
			return
		}
		b.SL.Error(err)
		b.SendMessage(mes.ChannelID, m.ErrorServer)
		return
	}
	b.SL.Info("Created the simulation")

	if file, err = os.Open(resultsFilePath); err != nil {
		b.SL.Error(err)
		b.SendMessage(mes.ChannelID, m.ErrorServer)
		return
	}

	if _, err = b.Session.ChannelFileSendWithMessage(
		mes.ChannelID,
		fmt.Sprintf("<@%s>", mes.Author.ID),
		fmt.Sprintf("%s_%s%s", unidecode.Unidecode(char), now, htmlExt),
		file,
	); err != nil {
		b.SL.Error(err)
		b.SendMessage(mes.ChannelID, m.ErrorServer)
	}
	b.SL.Info("Sent the file to the user")
}

func (b *Bot) simcProfileReporter(mes *discordgo.MessageCreate, command string, withStats, forPtr bool) {
	const timeFormat = "20060102_150405"
	var (
		simcExt = ".simc"
		htmlExt = ".html"

		argString, char, output, profileName string

		args []string
		file *os.File
		err  error
	)

	b.SL.Info("getting simcraft sim from profile...")

	char = strings.Split(mes.Content, " ")[2]
	location, _ := time.LoadLocation(o.GuildTimezone)
	now := time.Now().In(location).Format(timeFormat)
	profileName = fmt.Sprintf("%s_%s", char, now)
	profileFilePath := fmt.Sprintf("/tmp/%s%s", profileName, simcExt)
	resultsFileName := fmt.Sprintf("%s%s", profileName, htmlExt)
	resultsFilePath := "/tmp/" + resultsFileName

	b.SendMessage(mes.ChannelID, fmt.Sprintf(m.SimcProfile, char))

	if forPtr {
		command = o.SimcCmdPtr
	} else {
		command = o.SimcCmdStable
	}

	// save file
	if err = DownloadFile(profileFilePath, mes.Attachments[0].URL); err != nil {
		b.SL.Errorf("Unable to download the file: %s", err)
		b.SendMessage(mes.ChannelID, m.ErrorServer)
		return
	}

	b.SendMessage(mes.ChannelID, m.SimcImportSuccess)

	if withStats {
		argString = fmt.Sprintf(o.SimcArgsWithStats, profileFilePath, resultsFilePath)
	} else {
		argString = fmt.Sprintf(o.SimcArgsNoStats, profileFilePath, resultsFilePath)
	}
	args = strings.Split(argString, "|")

	output, err = ExecuteCommand(command, o.SimcDir, args)
	// b.SL.Info(output)
	if err != nil {
		if strings.Contains(output, "Character not found") {
			b.SL.Error("Unable to find the character")
			b.SendMessage(mes.ChannelID, m.SimcArmoryError)
			return
		}
		b.SL.Error(err)
		b.SendMessage(mes.ChannelID, m.ErrorServer)
		return
	}
	b.SL.Info("Created the simulation")

	if file, err = os.Open(resultsFilePath); err != nil {
		b.SL.Error(err)
		b.SendMessage(mes.ChannelID, m.ErrorServer)
		return
	}

	if _, err = b.Session.ChannelFileSendWithMessage(
		mes.ChannelID,
		fmt.Sprintf("<@%s>", mes.Author.ID),
		fmt.Sprintf("%s_%s%s", unidecode.Unidecode(char), now, htmlExt),
		file,
	); err != nil {
		b.SL.Error(err)
		b.SendMessage(mes.ChannelID, m.ErrorServer)
	}
	b.SL.Info("Sent the file to the user")
}

func (b *Bot) queueReporter(mes *discordgo.MessageCreate) {
	b.SL.Info("getting realm name string...")
	realmString := GetRealmName(mes.Content, "!queue")
	b.SL.Infof("getting queue status of %s and sending it...", realmString)
	if realmQueue, err := GetRealmQueueStatus(realmString); err != nil {
		b.SL.Errorf("Unable to get the realm queue status: %s", err)
		b.SendMessage(mes.ChannelID, m.ErrorUser)
	} else if realmQueue {
		b.SendMessage(mes.ChannelID, m.RealmQueue)
	} else {
		b.SendMessage(mes.ChannelID, m.RealmNoQueue)
	}
}

func (b *Bot) guildMembersReporter(mes *discordgo.MessageCreate) {
	b.SendMessage(mes.ChannelID, m.GuildMembersList)
	b.SL.Info("getting parametes string slice...")
	var parameters []string
	paramString := strings.TrimPrefix(mes.Content, "!guildmembers")
	paramString = strings.TrimPrefix(paramString, " ")
	if paramString != "" {
		parameters = strings.Split(paramString, " ")
		b.SL.Info("paramString:", paramString, "parameters len:", len(parameters))
	}
	b.SL.Info("getting guild members list and sending it...")
	tab := fmttab.New("Список согильдейцев", fmttab.BorderDouble, nil)
	tab.AddColumn("Имя", fmttab.WidthAuto, fmttab.AlignLeft).
		AddColumn("Уровень", 7, fmttab.AlignLeft).
		AddColumn("Класс", 18, fmttab.AlignLeft).
		AddColumn("Специализация", 18, fmttab.AlignLeft).
		AddColumn("iLevel", 6, fmttab.AlignLeft).
		AddColumn("Армори", 22, fmttab.AlignLeft)
	b.CharMutex.Lock()
	for _, char := range b.HighLvlCharacters {
		char.RLock()
		tab.AppendData(map[string]interface{}{
			"Имя":           char.Name,
			"Уровень":       char.Level,
			"Класс":         char.Class,
			"Специализация": char.Spec.Name,
			"iLevel":        char.Items.AvgItemLvlEq,
			"Армори":        char.Link,
		})
		char.RUnlock()
	}
	b.CharMutex.Unlock()
	b.SendMessage(mes.ChannelID, "```"+tab.String()+"```")
}

func (b *Bot) guildProfsReporter(mes *discordgo.MessageCreate) {
	b.SendMessage(mes.ChannelID, m.GuildProfsList)
	b.SL.Info("getting parametes string slice...")
	paramString := strings.TrimPrefix(mes.Content, "!guildprofs")
	paramString = strings.TrimPrefix(paramString, " ")
	b.SL.Info("getting guild profs list and sending it...")
	guildProfsInfo, err := GetGuildProfs(o.GuildRealm, o.GuildName, paramString)
	if err != nil {
		b.SL.Errorf("Unable to get guild professions: %s", err)
		return
	}
	tab := fmttab.New("Список профессий в гильдии", fmttab.BorderDouble, nil)
	tab.AddColumn("Имя", fmttab.WidthAuto, fmttab.AlignLeft).
		AddColumn("1 профа", 15, fmttab.AlignLeft).
		AddColumn("Уровень 1 профы", fmttab.WidthAuto, fmttab.AlignLeft).
		AddColumn("2 профа", 15, fmttab.AlignLeft).
		AddColumn("Уровень 2 профы", fmttab.WidthAuto, fmttab.AlignLeft)
	for _, member := range guildProfsInfo {
		tab.AppendData(map[string]interface{}{
			"Имя":             member["Name"],
			"1 профа":         member["FirstProf"],
			"Уровень 1 профы": member["FirstProfLevel"],
			"2 профа":         member["SecondProf"],
			"Уровень 2 профы": member["SecondProfLevel"],
		})
	}
	b.SendMessage(mes.ChannelID, "```"+tab.String()+"```")
}

func (b *Bot) realmInfoReporter(mes *discordgo.MessageCreate) {
	b.SL.Info("getting realm name string...")
	realmString := GetRealmName(mes.Content, "!realminfo")
	b.SL.Info(realmString)
	b.SL.Info("getting realm info and sending it...")
	realmInfo, err := GetRealmInfo(realmString)
	if err != nil {
		b.SL.Errorf("Unable to get guild info: %s", err)
		b.SendMessage(mes.ChannelID, m.ErrorUser)
	} else {
		b.SendMessage(mes.ChannelID, realmInfo)
	}
}

func (b *Bot) cleanUp(mes *discordgo.MessageCreate) {
	b.SL.Info("Removing bot messages...")
	var err error
	user := mes.Author.Username
	am := strings.Replace(mes.Message.Content, "!clean", "", 1)
	am = strings.Replace(am, " ", "", -1)
	b.SL.Infof("User %s - amount to delete: %s", user, am)
	var amount int
	switch am {
	case "all":
		amount = -1
	case "":
		amount = 1
	default:
		if amount, err = strconv.Atoi(am); err != nil {
			b.SL.Error(err)
			return
		}
	}
	if mes.ChannelID == o.GeneralChannelID && !containsUser(o.Admins, mes.Author.ID) && (amount > 3 || amount == -1) {
		b.SL.Info("User is trying to delete all bot messages from main channel! Won't work!")
		b.SendMessage(mes.ChannelID, m.Clean)
		return
	}
	lastMessageChecked := mes.ID
	chanMessages, _ := b.Session.ChannelMessages(mes.ChannelID, 100, lastMessageChecked, "", "")
	mesToDelete := make(map[string]string)
	for {
		if len(mesToDelete) == amount {
			break
		}
		for _, mes := range chanMessages {
			b.SL.Infof("%s %s %s", mes.ID, mes.Author.Username, mes.Author.ID)
			lastMessageChecked = mes.ID
			if mes.Author.ID == b.ID {
				if _, ok := mesToDelete[mes.ID]; !ok {
					mesToDelete[mes.ID] = mes.ID
				}
				if len(mesToDelete) == amount {
					break
				}
			}
		}
		chm, _ := b.Session.ChannelMessages(mes.ChannelID, 100, lastMessageChecked, "", "")
		if compareMesArrays(chm, chanMessages) {
			b.SL.Info("Reached the end, exiting loop...")
			break
		}
		chanMessages = chm
	}
	for _, mID := range mesToDelete {
		if err = b.Session.ChannelMessageDelete(mes.ChannelID, mID); err != nil {
			b.SL.Errorf("Unable to delete the message: %s", err)
		}
	}
	b.SL.Info("Deleted all messages")
	return
}
