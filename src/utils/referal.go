package utils

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ReferalResult struct {
	Type string   `json:"type"`
	Id   []string `json:"id"`
}

type Referal struct {
	Id        string
	Code      string
	Result    []ReferalResult
	IsClaimed bool
}

func ReferalCheck(code string, guildId string, claimer string) Referal {
	referal := []ReferalResult{}
	row := DB.QueryRow("SELECT id, benefit FROM referal WHERE code = ? AND guild = ?", code, guildId)
	var id string
	var benefit = []byte("")
	err := row.Scan(&id, &benefit)
	if err != nil {
		// If no rows or other error, return empty result
		return Referal{}
	}

	var is_claimed string
	row = DB.QueryRow("SELECT identifier AS is_claimed FROM referal_used WHERE code = ? AND guild = ? AND `by` = ? AND ignored = 0", code, guildId, claimer)
	err = row.Scan(&is_claimed)
	if err != nil {
		log.Println("Error while checking referal claimer")
		log.Fatalln(err)
	}
	if is_claimed == id {
		return Referal{
			Id:        id,
			Code:      code,
			IsClaimed: true,
		}
	}

	err = json.Unmarshal(benefit, &referal)
	if err != nil {
		return Referal{}
	}

	if len(referal) > 0 {
		return Referal{
			Id:     id,
			Code:   code,
			Result: referal,
		}
	}

	return Referal{}
}

func ReferalApply(s *discordgo.Session, grant_to string, code string, guildId string) Referal {
	referal := ReferalCheck(code, guildId, grant_to)
	_, err := s.Guild(guildId)
	if len(referal.Result) == 0 || err != nil {
		if referal.IsClaimed {
			return referal
		}
		return Referal{}
	}
	_, err = DB.Exec("INSERT INTO referal_used (code, guild, date, `by`, identifier) VALUES (?, ?, ?, ?, ?)", code, guildId, time.Now(), grant_to, referal.Id)
	if err != nil {
		log.Println("Error while inserting referal claimer")
		log.Fatalln(err)
		return Referal{}
	}
	for _, refer_to := range referal.Result {
		switch refer_to.Type {
		case "channel":
			for _, channelId := range refer_to.Id {
				_, err := s.Channel(channelId)
				if err == nil {
					s.ChannelEdit(channelId, &discordgo.ChannelEdit{
						PermissionOverwrites: []*discordgo.PermissionOverwrite{
							{
								ID:   grant_to,
								Type: discordgo.PermissionOverwriteTypeMember,
								Allow: discordgo.PermissionViewChannel |
									discordgo.PermissionSendMessages,
							},
						},
					})
				}
			}
		case "role":
			for _, roleId := range refer_to.Id {
				_, err := s.GuildRoles(roleId)
				if err == nil {
					s.ChannelEdit(roleId, &discordgo.ChannelEdit{
						PermissionOverwrites: []*discordgo.PermissionOverwrite{
							{
								ID:   grant_to,
								Type: discordgo.PermissionOverwriteTypeMember,
								Allow: discordgo.PermissionViewChannel |
									discordgo.PermissionSendMessages,
							},
						},
					})
				}
			}
		}
	}
	_, err = DB.Exec("UPDATE referal SET used = used + 1 WHERE id = ?", referal.Id)
	if err != nil {
		log.Println("Error while updating referal used count")
	}
	return referal
}
