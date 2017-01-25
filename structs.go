package main

import "time"

// Realm - type for WoW server realm info
type Realm struct {
    Type                string          `json:"type"`
    Population          string          `json:"population"`
    Queue               bool            `json:"queue"`
    Status              bool            `json:"status"`
    Name                string          `json:"name"`
    Slug                string          `json:"slug"`
    Battlegroup         string          `json:"battlegroup"`
    Locale              string          `json:"locale"`
    Timezone            string          `json:"timezone"`
    ConnectedRealms     []string        `json:"connected_realms"`
}

// Realms - struct for a slice of Realm
type Realms struct {
    RealmList           []Realm        `json:"realms"`
}

// GuildInfo - struct for WoW guild information
type GuildInfo struct {
    Name                string        `json:"name"`
    Realm               string        `json:"realm"`
    BattleGroup         string        `json:"battlegroup"`
    Level               int           `json:"level"`
    SideInt             int           `json:"side"`
    Side                string
    AchievementPoints   int           `json:"achievementPoints"`
    LastModified        int64         `json:"lastModified"`
    GuildMembersList    MembersList   `json:"members"`
    GuildNewsList       NewsList      `json:"news"`
}

// MembersList - type for a slice of GuildMember
type MembersList []GuildMember

// NewsList - type for a slice of News
type NewsList    []News

// GuildMember - struct for a WoW guild member
type GuildMember struct {
    Member              Character       `json:"character"`
    Rank                int             `json:"rank"`
}

// News - struct for any WoW guild news
type News struct {
    Type                string          `json:"type"`
    Character           string          `json:"character"`
    Timestamp           int64           `json:"timestamp"`
    EventTime           time.Time
    ItemID              int             `json:"itemId"`
    ItemInfo            Item
    Context             string          `json:"context"`
    BonusLists          []int           `json:"bonusLists"`
    Achievement         Achievement     `json:"achievement"`
}

// Achievement - struct for a WoW achievement
type Achievement struct {
    ID                  int             `json:"id"`
    Title               string          `json:"title"`
    Points              int             `json:"points"`
    Description         string          `json:"description"`
    RewardItems         []Item          `json:"rewardItems"`
    Icon                string          `json:"icon"`
    Criteria            []Criteria      `json:"criteria"`
    AccountWide         bool            `json:"accountWide"`
    FactionID           int             `json:"factionId"`
}

// Criteria - struct for a WoW achievement criteria
type Criteria struct {
    ID                  int             `json:"id"`
    Description         string          `json:"description"`
    OrderIndex          int             `json:"orderIndex"`
    Max                 int             `json:"max"`
}

// Item - partly filled struct for obtaining WoW item info
// TODO: Fill it with all the fields from the server response
type Item struct {
    ID                  int             `json:"id"`
    Name                string          `json:"name"`
    Quality             int             `json:"quality"`
    ItemLevel           int             `json:"itemLevel"`
    Equippable          bool            `json:"equippable"`
    ReqLevel            int             `json:"requiredLevel"`
    Link                string
}

// Character - struct for a WoW character
type Character struct {
    Name                string          `json:"name"` 
    Realm               string          `json:"realm"`
    RealmSlug           string
    FactionInt          int             `json:"faction"`
    Faction             string
    BattleGroup         string          `json:"battlegroup"` 
    ClassInt            int             `json:"class"`
    Class               string
    RaceInt             int             `json:"race"`
    Race                string 
    GenderInt           int             `json:"gender"`
    Gender              string
    Level               int             `json:"level"` 
    AchievementPoints   int             `json:"achievementPoints"` 
    Thumbnail           string          `json:"thumbnail"` 
    Spec                Specialization  `json:"spec"` 
    Guild               string          `json:"guild"` 
    GuildRealm          string          `json:"guildRealm"` 
    LastModified        int             `json:"lastModified"`
    Items               Items           `json:"items"`
    Professions         Professions     `json:"professions"`
    Link                string
}

// Specialization - struct for a WoW character specialization
type Specialization struct {
    Name                string          `json:"name"`
    Role                string          `json:"role"`
    BackgroundImage     string          `json:"backgroundImage"`
    Icon                string          `json:"icon"`
    Description         string          `json:"description"`
    Order               int             `json:"order"`
}

// Items - struct for storing items info for a character
type Items struct {
    AvgItemLvl          int             `json:"averageItemLevel"`
    AvgItemLvlEq        int             `json:"averageItemLevelEquipped"`
}

// Professions - struct for professions info for a character
type Professions struct {
    PrimaryProfs        []Profession    `json:"primary"`
    SecondaryProfs      []Profession    `json:"secondary"`
}

// Profession - struct for a profession info for a character
type Profession struct {
    ID                  int             `json:"id"`
    Name                string          `json:"name"`
    EngName             string
    Icon                string          `json:"icon"`
    Rank                int             `json:"rank"`
    Max                 int             `json:"max"`
    Recipes             []int           `json:"recipes"`
    Link                string
}

type googlAPIResponse struct {
    Kind                string          `json:"kind"`
    ID                  string          `json:"id"`
    LongURL             string          `json:"longUrl"`
}