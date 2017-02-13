package main

// Public maps
var ()

// Options struct holds all the options
type Options struct {
	DiscordToken string
	WoWToken     string
	GoogleToken  string

	SSHAddress string
	SSHUser    string

	SimcImport    string
	SimcNoStats   string
	SimcWithStats string

	Bucket string

	Admins []string

	GeneralChannelID string

	GuildRegion   string
	GuildLocale   string
	GuildTimezone string
	GuildName     string
	GuildRealm    string

	WoWClasses     map[int64]string
	WoWGenders     map[int64]string
	WoWFactions    map[int64]string
	WoWRaces       map[int64]string
	WoWProfessions map[int64]string

	WowheadItemLink     string
	GoogleShortenerLink string

	APIRealmsLink       string
	APIGuildMembersLink string
	APIGuildNewsLink    string
	APICharItemsLink    string
	APICharNewsLink     string
	APICharProfsLink    string
	APIItemLink         string

	WoWDBItemLink string

	ArmoryCharLink string
	ArmoryProfLink string
}

// Messages struct holds all the bot message strings
type Messages struct {
	Help        string
	Pong        string
	Clean       string
	ErrorUser   string
	ErrorServer string

	Legendary string

	Boobies  string
	JohnCena string
	Relics   string
	Godbook  string
	Logs     string

	Roll1   string
	RollX   string
	Roll100 string

	RealmOn      string
	RealmOff     string
	RealmQueue   string
	RealmNoQueue string

	GuildMembersList string
	GuildProfsList   string

	SimcArmory        string
	SimcArmoryError   string
	SimcImport        string
	SimcImportSuccess string
	SimcProfile       string
}

// WoWConfig struct for WoW config slice
type WoWConfig struct {
	ID   int
	Name string
}
