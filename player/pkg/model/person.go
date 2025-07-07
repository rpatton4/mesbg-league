package model

import "github.com/rpatton4/mesbg-league/auth"

type Person struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	DiscordName   string          `json:"discordName"`
	SocialMediaID string          `json:"socialMediaUserId"`
	AuthSource    auth.AuthSource `json:"socialMediaAuthSource"`
}
