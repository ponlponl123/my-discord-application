package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Referal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "referal_" + i.Interaction.Member.User.ID,
			Title:    "Referal",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "code",
							Label:       "Referal",
							Style:       discordgo.TextInputShort,
							Placeholder: "XXX-XXX",
							Required:    true,
							MaxLength:   100,
							MinLength:   1,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "target",
							Label:       "Target",
							Style:       discordgo.TextInputShort,
							Placeholder: "0123456789",
							Required:    true,
							MaxLength:   100,
							MinLength:   1,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return
	}
}

func ReferalModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Processing...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return
	}
	data := i.ModalSubmitData()
	if !strings.HasPrefix(data.CustomID, "referal") {
		return
	}
	userid := strings.Split(data.CustomID, "_")[1]
	referal := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	target := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	resultMsg := "Failed, target is not valid or invalid referal code."
	if referal == "my-secret-referal" && target != "" {
		_, err := s.Channel(target)
		if err == nil {
			s.ChannelEdit(i.ChannelID, &discordgo.ChannelEdit{
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:   userid,
						Type: discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionViewChannel |
							discordgo.PermissionSendMessages,
					},
				},
			})
			resultMsg = "## Referal success!\nPermission granted to <@" + userid + ">"
		}
	}
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &resultMsg,
	})
}
