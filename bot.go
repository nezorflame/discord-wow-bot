package main

import (
    "flag"
    "fmt"
    "strings"
    "log"
    "time"
    "os"
    "errors"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "github.com/bwmarrin/discordgo"
)

const (
    // global consts
    Pong        = "Pong!"
    JohnCena    = "AND HIS NAME IS JOOOOOOOOOHN CEEEEEEEEEEEENAAAAAAAA! https://youtu.be/QQUgfikLYNI"
    Relics      = "https://docs.google.com/spreadsheets/d/11RqT6EIelFWHB1b8f_scFo8sPdXGVYFii_Dr7kkOFLY/edit#gid=1060702296"
    RGB         = "https://docs.google.com/spreadsheets/d/1apphJ2vlZL4eQFZMKeUrYC34PsNt7JFeTZiqNtb0NyE/htmlview?sle=true"
    RealmOn     = "Сервер онлайн! :)"
    RealmOff    = "Сервер оффлайн! :()"
)

var (
    logger      *log.Logger
    startTime   time.Time
	// Token for bot auth
    Token       string
    // BotID for bot ID
    BotID       string
    // Realms is a structure array for WoW realms
    Realms      []Realm
)

// Realm - type for WoW server realm info
type Realm struct {
    Type            string      `json:"type"`
    Population      string      `json:"population"`
    Queue           bool        `json:"queue"`
    Status          bool        `json:"status"`
    Name            string      `json:"name"`
    Slug            string      `json:"slug"`
    Battlegroup     string      `json:"battlegroup"`
    Locale          string      `json:"locale"`
    Timezone        string      `json:"timezone"`
    ConnectedRealms []string    `json:"connected_realms"`
}
// RealmsAPIResponse is a struct for slice of Realm
type RealmsAPIResponse struct {
    RealmList []Realm `json:"realms"`
}

func init() {
    // Create initials.
	logger = log.New(os.Stderr, "  ", log.Ldate|log.Ltime)
	startTime = time.Now()
    getWoWRealms()

    // Parse command line arguments.
    flag.StringVar(&Token, "t", "", "Account Token")
    flag.Parse()
    if Token == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func logDebug(v ...interface{}) {
	logger.SetPrefix("DEBUG ")
	logger.Println(v...)
}

func logInfo(v ...interface{}) {
	logger.SetPrefix("INFO  ")
	logger.Println(v...)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

/* Tries to call a method and checking if the method returned an error, if it
did check to see if it's HTTP 502 from the Discord API and retry for
`attempts` number of times. */
func retryOnBadGateway(f func() error) {
	var err error
	for i := 0; i < 3; i++ {
		err = f()
		if err != nil {
			if strings.HasPrefix(err.Error(), "HTTP 502") {
				// If the error is Bad Gateway, try again after 1 sec.
				time.Sleep(1 * time.Second)
				continue
			} else {
				// Otherwise panic !
				panicOnErr(err)
			}
		} else {
			// In case of no error, return.
			return
		}
	}
}

func sendMessage(session *discordgo.Session, chID string, message string) error {
    logInfo("SENDING MESSAGE:", message)
	retryOnBadGateway(func() error {
		_, err := session.ChannelMessageSend(chID, message)
		return err
	})
    return nil
}

func main() {
    logInfo("Logging in...")
    session, err := discordgo.New(Token)
    logInfo("Using bot account token...")
    u, err := session.User("@me")
    panicOnErr(err)
    BotID = u.ID
    logInfo("Got BotID =", BotID)
    logInfo("session token is " + session.Token)
    setupHandlers(session)
	panicOnErr(err)
    logInfo("Opening session...")
	err = session.Open()
	panicOnErr(err)
	fmt.Println("Bot is now running.\nPress CTRL-C to exit...")
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}
    // Check the command to answer
    if strings.HasPrefix(m.Content, "!status") {
        realmString := strings.Split(m.Content, " ")
        realmStatus, err := getWoWRealmStatus(realmString[1])
        panicOnErr(err)
        switch realmStatus {
            case true:
                err = sendMessage(s, m.ChannelID, RealmOn)
                panicOnErr(err)
            default:
                err = sendMessage(s, m.ChannelID, RealmOff)
                panicOnErr(err)
        }
    }
    switch m.Content {
        case "!ping":
            err := sendMessage(s, m.ChannelID, Pong)
            panicOnErr(err)
        case "!johncena":
            err := sendMessage(s, m.ChannelID, JohnCena)
            panicOnErr(err)
        case "!status":
            getWoWRealmStatus("")
        case "!relics":
            err := sendMessage(s, m.ChannelID, Relics)
            panicOnErr(err)
        case "!godbook":
            err := sendMessage(s, m.ChannelID, RGB)
            panicOnErr(err)
        default:
            log.Println("not a command")
    }
}

func setupHandlers(session *discordgo.Session) {
	logInfo("Setting up event handlers...")
	session.AddHandler(messageCreate)
}

func containsUser(users []*discordgo.User, userID string) bool {
    for _, u := range users {
        if u.ID == userID {
            return true
        }
    }
    return false
}

func getWoWRealms() {
    r, err := http.Get("https://us.api.battle.net/wow/realm/status?locale=en_US&apikey=fdvxqkq6qkq364brvgkwuur73u5dncw8")
    panicOnErr(err)
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    panicOnErr(err)
    realmsResponse, err := getRealms([]byte(body))
    panicOnErr(err)
    Realms = realmsResponse.RealmList
}

func getRealms(body []byte) (*RealmsAPIResponse, error) {
    var s = new(RealmsAPIResponse)
    err := json.Unmarshal(body, &s)
    panicOnErr(err)
    return s, err
}

func getWoWRealmStatus(realmName string) (bool, error) {
    for _, r := range Realms {
        if r.Name == realmName || r.Slug == realmName {
            return r.Status, nil
        }
    }
    return false, errors.New("No such realm is present!")
}