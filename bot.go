package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"os"

	"github.com/arteev/fmttab"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/glog"
)

// Users - map of guild members
var (
	BotID string
	Users map[string]string
)

// Start - function to start Discord bot
func Start() {
	var (
		err     error
		u       *discordgo.User
		session *discordgo.Session
	)
	InitializeWoWAPI()
	glog.Info("Logging in...")
	if session, err = discordgo.New(o.DiscordToken); err != nil {
		glog.Fatalf("Unable to connect to Discord: %s", err)
	}
	glog.Info("Using bot account token...")
	if u, err = session.User("@me"); err != nil {
		glog.Fatalf("Unable to get @me: %s", err)
	} else {
		BotID = u.ID
		glog.Infof("Got BotID = %s", BotID)
	}
	glog.Info("Adding handlers...")
	setup(session)
	glog.Info("Opening session...")
	if err = session.Open(); err != nil {
		glog.Fatalf("Unable to open the session: %s", err)
	}
	glog.Info("Starting guild watcher and spammer...")
	go RunGuildWatcher(session)
	glog.Info("Bot started")
}

/* Tries to call a method and checking if the method returned an error, if it
did check to see if it's HTTP 502 from the Discord API and retry for
`attempts` number of times. */
func retryOnBadGateway(f func() error) {
	var err error
	for i := 0; i < 3; i++ {
		if err = f(); err != nil {
			if strings.HasPrefix(err.Error(), "HTTP 502") {
				// If the error is Bad Gateway, try again after 1 sec.
				time.Sleep(1 * time.Second)
				continue
			} else {
				// Otherwise panic !
				glog.Fatal(err)
			}
		} else {
			// In case of no error, return.
			return
		}
	}
}

func sendMessage(session *discordgo.Session, chID string, message string) (err error) {
	glog.Info("SENDING MESSAGE:", message)
	retryOnBadGateway(func() error {
		return sendFormattedMessage(session, chID, message)
	})
	return
}

func sendFormattedMessage(session *discordgo.Session, chID string, message string) (err error) {
	i := len(message)
	if len(message) > 1999 {
		for i > 1999 {
			messageSlice := strings.Split(message, "\n")
			mes := messageSlice[0]
			l := len(messageSlice)
			if l == 2 {
				_, err = session.ChannelMessageSend(chID, mes+"\n")
				if err != nil {
					return err
				}
				_, err = session.ChannelMessageSend(chID, messageSlice[1])
				if err != nil {
					return err
				}
				return nil
			}
		Loop:
			for j := 1; j < l-1; j++ {
				mes += "\n" + messageSlice[j]
				if len(mes+"\n"+messageSlice[j+1]) > 1999 {
					if strings.HasPrefix(mes, "```") {
						_, err = session.ChannelMessageSend(chID, mes+"```")
					} else {
						_, err = session.ChannelMessageSend(chID, mes)
					}
					if err != nil {
						return err
					}
					message = strings.Replace(message, mes, "", 1)
					if strings.HasPrefix(mes, "```") {
						message = "```" + message
					}
					break Loop
				}
			}
			i = len(message)
		}
		_, err = session.ChannelMessageSend(chID, message)
	} else {
		_, err = session.ChannelMessageSend(chID, message)
	}
	return
}

func logPinnedMessages(s *discordgo.Session) {
	var (
		err    error
		pinned []*discordgo.Message
	)
	glog.Info("getPinnedMessages called")
	if pinned, err = s.ChannelMessagesPinned(o.GeneralChannelID); err != nil {
		glog.Errorf("Unable to get pinned messages: %s", err)
		return
	}
	glog.Info(len(pinned), "messages are pinned:")
	for _, message := range pinned {
		glog.Info("["+message.ID+"]", message.Content)
	}
}

func printMessageByID(s *discordgo.Session, chID string, mesID string) {
	glog.Info("printMessageByID called")
	message, err := s.ChannelMessage(o.GeneralChannelID, mesID)
	if err != nil {
		glog.Errorf("Unable to get the message: %s", err)
		return
	}
	if err = sendMessage(s, chID, message.Content); err != nil {
		glog.Errorf("Unable to send the message: %s", err)
	}
}

func setup(session *discordgo.Session) {
	glog.Info("Setting up event handlers...")
	session.AddHandler(messageCreate)
}

// RunGuildWatcher - function for starting the guild news watcher
// TODO: Very dirty, need to rewrite
func RunGuildWatcher(s *discordgo.Session) {
	var (
		err         error
		messages    []string
		legendaries = make(map[string]bool)
	)

	for {
		if messages, err = GetGuildLegendaries(o.GuildRealm, o.GuildName); err != nil {
			glog.Errorf("Unable to get guild legendaries: %s", err)
			goto Sleep
		}
		for _, m := range messages {
			if _, ok := legendaries[m]; !ok {
				if err = sendMessage(s, o.GeneralChannelID, m); err != nil {
					glog.Errorf("Unable to send the message: %s", err)
				}
				glog.Info(m)
				legendaries[m] = true
			}
		}
	Sleep:
		time.Sleep(5 * time.Minute)
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, mes *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if mes.Author.ID == BotID {
		return
	}
	// Check the command to react and answer
	if strings.HasPrefix(mes.Content, "!status") {
		statusReporter(s, mes)
	}
	if strings.HasPrefix(mes.Content, "!simc") {
		simcReporter(s, mes)
	}
	if strings.HasPrefix(mes.Content, "!queue") {
		queueReporter(s, mes)
	}
	if strings.HasPrefix(mes.Content, "!realminfo") {
		realmInfoReporter(s, mes)
	}
	if strings.HasPrefix(mes.Content, "!guildmembers") {
		if err := sendMessage(s, mes.ChannelID, m.GuildMembersList); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
		guildMembersReporter(s, mes)
	}
	if strings.HasPrefix(mes.Content, "!guildprofs") {
		if err := sendMessage(s, mes.ChannelID, m.GuildProfsList); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
		guildProfsReporter(s, mes)
	}
	if strings.HasPrefix(mes.Content, "!clean") {
		cleanUp(s, mes)
	}
	if strings.HasPrefix(mes.Content, "!announce") {
		message := strings.TrimPrefix(mes.Message.Content, "!announce")
		glog.Info(mes.Author.Username, "is announcing a message:", message)
		if err := sendMessage(s, o.GeneralChannelID, message); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
	}
	switch mes.Content {
	case "!ping":
		if err := sendMessage(s, mes.ChannelID, m.Pong); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
	case "!roll":
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
		if err := sendMessage(s, mes.ChannelID, message); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
	case "!johncena":
		if err := sendMessage(s, mes.ChannelID, m.JohnCena); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
	case "!relics":
		if err := sendMessage(s, mes.ChannelID, m.Relics); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
	case "!godbook":
		if err := sendMessage(s, mes.ChannelID, m.Godbook); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
	case "!logs":
		if err := sendMessage(s, mes.ChannelID, m.Logs); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
	case "!help", "!помощь":
		helpReporter(s, mes)
	case "!boobs":
		boobsReporter(s, mes)
	case "!!printpinned":
		logPinnedMessages(s)
	case "!!terminate":
		panic("Terminating ..")
	}
}

func cleanUp(s *discordgo.Session, mes *discordgo.MessageCreate) {
	glog.Info("Removing bot messages...")
	var err error
	user := mes.Author.Username
	am := strings.Replace(mes.Message.Content, "!clean", "", 1)
	am = strings.Replace(am, " ", "", -1)
	glog.Infof("User %s - amount to delete: %s", user, am)
	var amount int
	switch am {
	case "all":
		amount = -1
	case "":
		amount = 1
	default:
		if amount, err = strconv.Atoi(am); err != nil {
			glog.Error(err)
			return
		}
	}
	if mes.ChannelID == o.GeneralChannelID && !containsUser(o.Admins, mes.Author.ID) && (amount > 3 || amount == -1) {
		glog.Info("User is trying to delete all bot messages from main channel! Won't work!")
		if err = sendMessage(s, mes.ChannelID, m.Clean); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
		return
	}
	lastMessageChecked := mes.ID
	chanMessages, _ := s.ChannelMessages(mes.ChannelID, 100, lastMessageChecked, "")
	mesToDelete := make(map[string]string)
	for {
		if len(mesToDelete) == amount {
			break
		}
		for _, mes := range chanMessages {
			glog.Infof("%s %s %s", mes.ID, mes.Author.Username, mes.Author.ID)
			lastMessageChecked = mes.ID
			if mes.Author.ID == BotID {
				if _, ok := mesToDelete[mes.ID]; !ok {
					mesToDelete[mes.ID] = mes.ID
				}
				if len(mesToDelete) == amount {
					break
				}
			}
		}
		chm, _ := s.ChannelMessages(mes.ChannelID, 100, lastMessageChecked, "")
		if compareMesArrays(chm, chanMessages) {
			glog.Info("Reached the end, exiting loop...")
			break
		}
		chanMessages = chm
	}
	for _, mID := range mesToDelete {
		if err = s.ChannelMessageDelete(mes.ChannelID, mID); err != nil {
			glog.Errorf("Unable to delete the message: %s", err)
		}
	}
	glog.Info("Deleted all messages")
	return
}

func helpReporter(s *discordgo.Session, mes *discordgo.MessageCreate) {
	glog.Info("Sending help to user...")
	if err := sendMessage(s, mes.ChannelID, m.Help); err != nil {
		glog.Errorf("Unable to send the message: %s", err)
	}
}

func boobsReporter(s *discordgo.Session, mes *discordgo.MessageCreate) {
	glog.Info("Sending boobies to user...:)")
	if err := sendMessage(s, mes.ChannelID, m.Boobies); err != nil {
		glog.Errorf("Unable to send the message: %s", err)
	}
}

func statusReporter(s *discordgo.Session, mes *discordgo.MessageCreate) {
	glog.Info("getting realm name string...")
	realmString := GetRealmName(mes.Content, "!status")
	glog.Infof("getting status of %s and sending it...", realmString)
	if realmStatus, err := GetRealmStatus(realmString); err != nil {
		glog.Errorf("Unable to get the realm status: %s", err)
	} else if realmStatus {
		if sErr := sendMessage(s, mes.ChannelID, m.RealmOn); sErr != nil {
			glog.Error(sErr)
		}
	} else {
		if sErr := sendMessage(s, mes.ChannelID, m.RealmOff); sErr != nil {
			glog.Error(sErr)
		}
	}
}

func simcReporter(s *discordgo.Session, mes *discordgo.MessageCreate) {
	var (
		simcExt = ".simc"
		htmlExt = ".html"

		command, output, profileName string

		file *os.File

		err error
	)
	glog.Info("getting simcraft sim...")
	params := strings.Split(strings.Replace(mes.Content, "!simc ", "", 1), " ")
	char := params[0]
	if len(params) == 0 || char == "!simc" || char == "" {
		glog.Infof("Command is incorrect: %s", mes.Content)
		if err = sendMessage(s, mes.ChannelID, m.ErrorUser); err != nil {
			glog.Errorf("Unable to send the message: %s", err)
		}
		return
	}

	profileName = fmt.Sprintf("%s_%d", mes.Author.Username, time.Now().Unix())
	profileFilePath := fmt.Sprintf("/tmp/%s%s", profileName, simcExt)
	resultsFileName := fmt.Sprintf("%s%s", profileName, htmlExt)
	resultsFilePath := "/tmp/" + resultsFileName
	realm := strings.Replace(o.GuildRealm, " ", "%20", -1)
	command = fmt.Sprintf(o.SimcImport, realm, char, profileFilePath)
	// for _, p := range params {
	// 	args := strings.Split(p, "=")
	// 	if len(args) != 2 {
	// 		continue
	// 	}
	// 	switch args[0] {
	// 	case "armory":
	// 		if args[1] == "no" {
	// 			isImported = true
	// 		} else {
	// 			strings.Replace(args[1], "_", "%20", -1)
	// 			profile = args[0] + "=" + args[1]
	// 		}
	// 	default:
	// 		command += " " + p
	// 	}
	// }
	glog.Info(command)

	if err = sendMessage(s, mes.ChannelID, fmt.Sprintf(m.SimcArmory, char)); err != nil {
		glog.Errorf("Unable to send the message: %s", err)
	}

	if err = ConnectToServer(o.SSHUser, o.SSHAddress); err != nil {
		glog.Errorf("Unable to connect to SSH: %s", err)
		return
	}
	defer SSHConn.Close()
	defer SFTPConn.Close()

	output, err = ExecuteCommand(command)
	glog.Info(output)
	if err != nil {
		glog.Error(err)
		if sErr := sendMessage(s, mes.ChannelID, m.ErrorUser); sErr != nil {
			glog.Errorf("Unable to send the message: %s", sErr)
		}
		return
	}
	glog.Info("Created the user profile from Armory")
	if err = sendMessage(s, mes.ChannelID, m.SimcImportSuccess); err != nil {
		glog.Errorf("Unable to send the message: %s", err)
	}

	command = fmt.Sprintf(o.SimcWithStats, profileFilePath, resultsFilePath)

	output, err = ExecuteCommand(command)
	glog.Info(output)
	if err != nil {
		glog.Error(err)
		if sErr := sendMessage(s, mes.ChannelID, m.ErrorServer); sErr != nil {
			glog.Errorf("Unable to send the message: %s", sErr)
		}
		return
	}
	glog.Info("Created the simulation")

	if err = DownloadFile(resultsFilePath, resultsFilePath); err != nil {
		glog.Error(err)
		if sErr := sendMessage(s, mes.ChannelID, m.ErrorServer); sErr != nil {
			glog.Errorf("Unable to send the message: %s", sErr)
		}
		return
	}
	glog.Info("Downloaded the results")

	if file, err = os.Open(resultsFilePath); err != nil {
		glog.Error(err)
		if sErr := sendMessage(s, mes.ChannelID, m.ErrorServer); sErr != nil {
			glog.Errorf("Unable to send the message: %s", sErr)
		}
		return
	}

	if _, err = s.ChannelFileSendWithMessage(
		mes.ChannelID,
		fmt.Sprintf("<@%s>", mes.Author.ID),
		mes.Author.ID+htmlExt,
		file,
	); err != nil {
		glog.Error(err)
		if sErr := sendMessage(s, mes.ChannelID, m.ErrorServer); sErr != nil {
			glog.Errorf("Unable to send the message: %s", sErr)
		}
	}
	glog.Info("Sent the file to the user")
}

func queueReporter(s *discordgo.Session, mes *discordgo.MessageCreate) {
	glog.Info("getting realm name string...")
	realmString := GetRealmName(mes.Content, "!queue")
	glog.Infof("getting queue status of %s and sending it...", realmString)
	if realmQueue, err := GetRealmQueueStatus(realmString); err != nil {
		glog.Errorf("Unable to get the realm queue status: %s", err)
	} else if realmQueue {
		if sErr := sendMessage(s, mes.ChannelID, m.RealmQueue); sErr != nil {
			glog.Error(sErr)
		}
	} else {
		if sErr := sendMessage(s, mes.ChannelID, m.RealmNoQueue); sErr != nil {
			glog.Error(sErr)
		}
	}
}

func guildMembersReporter(s *discordgo.Session, mes *discordgo.MessageCreate) {
	glog.Info("getting parametes string slice...")
	var parameters []string
	paramString := strings.TrimPrefix(mes.Content, "!guildmembers")
	paramString = strings.TrimPrefix(paramString, " ")
	if paramString != "" {
		parameters = strings.Split(paramString, " ")
		glog.Info("paramString:", paramString, "parameters len:", len(parameters))
	}
	glog.Info("getting guild members list and sending it...")
	guildMembersInfo, err := GetGuildMembers(o.GuildRealm, o.GuildName, parameters)
	if err != nil {
		glog.Errorf("Unable to get guild members: %s", err)
		return
	}
	tab := fmttab.New("Список согильдейцев", fmttab.BorderDouble, nil)
	tab.AddColumn("Имя", fmttab.WidthAuto, fmttab.AlignLeft).
		AddColumn("Уровень", 7, fmttab.AlignLeft).
		AddColumn("Класс", 18, fmttab.AlignLeft).
		AddColumn("Специализация", 18, fmttab.AlignLeft).
		AddColumn("iLevel", 6, fmttab.AlignLeft).
		AddColumn("Армори", 22, fmttab.AlignLeft)
	for _, member := range guildMembersInfo {
		tab.AppendData(map[string]interface{}{
			"Имя":           member["Name"],
			"Уровень":       member["Level"],
			"Класс":         member["Class"],
			"Специализация": member["Spec"],
			"iLevel":        member["ItemLevel"],
			"Армори":        member["Link"],
		})
	}
	if err = sendMessage(s, mes.ChannelID, "```"+tab.String()+"```"); err != nil {
		glog.Errorf("Unable to send the message: %s", err)
	}
}

func guildProfsReporter(s *discordgo.Session, mes *discordgo.MessageCreate) {
	glog.Info("getting parametes string slice...")
	paramString := strings.TrimPrefix(mes.Content, "!guildprofs")
	paramString = strings.TrimPrefix(paramString, " ")
	glog.Info("getting guild profs list and sending it...")
	guildProfsInfo, err := GetGuildProfs(o.GuildRealm, o.GuildName, paramString)
	if err != nil {
		glog.Errorf("Unable to get guild professions: %s", err)
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
	if err = sendMessage(s, mes.ChannelID, "```"+tab.String()+"```"); err != nil {
		glog.Errorf("Unable to send the message: %s", err)
	}
}

func realmInfoReporter(s *discordgo.Session, mes *discordgo.MessageCreate) {
	glog.Info("getting realm name string...")
	realmString := GetRealmName(mes.Content, "!realminfo")
	glog.Info(realmString)
	glog.Info("getting realm info and sending it...")
	realmInfo, err := GetRealmInfo(realmString)
	if err != nil {
		glog.Errorf("Unable to get guild info: %s", err)
	} else {
		if sErr := sendMessage(s, mes.ChannelID, realmInfo); sErr != nil {
			glog.Errorf("Unable to send the message: %s", sErr)
		}
	}
}

func containsUser(users []string, userID string) bool {
	for _, u := range users {
		if u == userID {
			return true
		}
	}
	return false
}

func compareMesArrays(a, b []*discordgo.Message) bool {
	for i := range a {
		if a[i].ID != b[i].ID {
			return false
		}
	}
	return true
}

func timeIsAllowed() bool {
	location, _ := time.LoadLocation(o.GuildTimezone)
	now := time.Now().In(location)
	hour := now.Hour()
	weekday := now.Weekday()
	switch weekday {
	// saturday has raids and is a holiday
	case time.Saturday:
		if !(hour >= 2 && hour <= 10 || hour >= 20 && hour <= 23) {
			glog.Info("Saturday spam :) time now:", now.String())
			return true
		}
	// sunday
	case time.Sunday:
		if !(hour >= 2 && hour <= 8) {
			glog.Info("Sunday spam :) time now:", now.String())
			return true
		}
	// work days
	default:
		if !(hour >= 2 && hour <= 8 || hour >= 20 && hour <= 23) {
			glog.Info("Workday spam :) time now:", now.String())
			return true
		}
	}
	return false
}
