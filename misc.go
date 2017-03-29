package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

// ExecuteCommand - function which runs the shell command
func ExecuteCommand(command, dir string, args []string) (output string, err error) {
	var bOutput []byte
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	bOutput, err = cmd.Output()
	output = string(bOutput)
	return
}

// Unmarshal makes the GuildInfo from the []byte
func (gi *GuildInfo) Unmarshal(body []byte) error {
	return json.Unmarshal(body, gi)
}

// Unmarshal makes the Realms from the []byte
func (r *Realms) Unmarshal(body []byte) error {
	return json.Unmarshal(body, r)
}

// Unmarshal makes the Character from the []byte
func (c *Character) Unmarshal(body []byte) (err error) {
	return json.Unmarshal(body, c)
}

// Unmarshal makes the Item from the []byte
func (i *Item) Unmarshal(body []byte) error {
	return json.Unmarshal(body, i)
}

// GetURLFromJSON returns short goo.gl link
func GetURLFromJSON(body []byte) (apiResponseID string, err error) {
	apiResponse := new(URLShortenerAPIResponse)
	if err = json.Unmarshal(body, apiResponse); err == nil {
		apiResponseID = apiResponse.ID
	}
	return
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

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
