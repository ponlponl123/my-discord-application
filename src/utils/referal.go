package utils

import (
	"encoding/json"
	"errors"
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
	Error     error
}

func ReferalCheck(code string, guildId string, claimer string) Referal {
	referal := []ReferalResult{}
	row := DB.QueryRow("SELECT id, benefit FROM referal WHERE code = ? AND guild = ? AND disabled = 0", code, guildId)
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
		// If no rows found, it means the referral is not claimed yet
		if err.Error() == "sql: no rows in result set" {
			is_claimed = ""
		} else {
			log.Println("Error while checking referal claimer:", err)
			return Referal{}
		}
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
		log.Println("Error:", err)
		return Referal{
			Error: errors.New("error while making referal claimer"),
		}
	}
	for _, refer_to := range referal.Result {
		switch refer_to.Type {
		case "channel":
			for _, channelId := range refer_to.Id {
				channel, err := s.Channel(channelId)
				if err == nil {
					_, err = s.ChannelEdit(channelId, &discordgo.ChannelEdit{
						PermissionOverwrites: append(channel.PermissionOverwrites, &discordgo.PermissionOverwrite{
							ID:    grant_to,
							Type:  discordgo.PermissionOverwriteTypeMember,
							Allow: discordgo.PermissionViewChannel,
						}),
					})
					if err != nil {
						log.Println("Error while granting channel permission")
						log.Println("Error:", err)
						return Referal{
							Error: errors.New("error while granting channel permission"),
						}
					}
				}
			}
		case "role":
			roles, err := s.GuildRoles(guildId)
			if err != nil {
				log.Println("Error while getting guild roles")
				log.Println("Error:", err)
				return Referal{}
			}
			roleSet := make(map[string]struct{}, len(roles))
			for _, role := range roles {
				roleSet[role.ID] = struct{}{}
			}
			for _, roleId := range refer_to.Id {
				if _, exists := roleSet[roleId]; exists {
					err = s.GuildMemberRoleAdd(guildId, grant_to, roleId)
					if err != nil {
						log.Println("Error while adding role to member")
						log.Println("Error:", err)
						return Referal{
							Error: errors.New("error while giving role to member"),
						}
					}
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
