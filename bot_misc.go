package main

import "github.com/bwmarrin/discordgo"

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
