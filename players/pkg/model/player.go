package model

import (
	"github.com/rpatton4/mesbg-league/pkg/auth"
	player "github.com/rpatton4/mesbg-league/players/pkg"
)

type Player struct {
	ID            player.PlayerID `json:"id"`
	Name          string          `json:"name"`
	DiscordName   string          `json:"discordName"`
	SocialMediaID string          `json:"socialMediaUserId"`
	AuthSource    auth.AuthSource `json:"socialMediaAuthSource"`
}
