package main

// Class consts
const (
	_ = iota
	Warrior
	Paladin
	Hunter
	Rogue
	Priest
	DeathKnight
	Shaman
	Mage
	Warlock
	Monk
	Druid
	DemonHunter
)

// Gender consts
const (
	Male = iota
	Female
)

// Faction consts
const (
	Alliance = iota
	Horde
)

// Race consts
const (
	_ = iota
	Human
	Orc
	Dwarf
	NElf
	Undead
	Tauren
	Gnome
	Troll
	Goblin
	BElf
	Draenei
	Worgen          = 10 + iota
	PandarenNeutral = 11 + iota
	PandarenAlliance
	PandarenHorde
)

// Profession consts
const (
	FirstAid       = 129
	Blacksmithing  = 164
	Leatherworking = 165
	Alchemy        = 171
	Herbalism      = 182
	Cooking        = 185
	Mining         = 186
	Tailoring      = 197
	engineering    = 202
	Enchanting     = 333
	Fishing        = 356
	Skinning       = 393
	Jewelcrafting  = 755
	Inscription    = 773
	Archaeology    = 794
)
