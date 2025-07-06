package model

type Person struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DiscordName  string `json:"discordName"`
	GoogleUserID string `json:"googleUserId"`
}
