package commands

import (
	"my-discord-bot/src/utils"
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
	resultMsg := "Failed, target is not valid or invalid referal code."
	result := utils.ReferalApply(s, userid, referal, i.GuildID)
	embeds := []*discordgo.MessageEmbed{}
	if result.IsClaimed {
		resultMsg = "## Referal already claimed!\nYou're not allowed to claim this referal again."
	}
	if result.Error != nil {
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:       "Internal Application Error!",
			Description: "```" + result.Error.Error() + "```",
			Color:       0xff0000,
		})
	}
	if len(result.Result) > 0 && !result.IsClaimed && result.Error == nil {
		resultMsg = "## Referal claimed!\nPermission granted to <@" + userid + ">"
		fields := []*discordgo.MessageEmbedField{}
		for _, v := range result.Result {
			var name string
			var value string
			switch v.Type {
			case "channel":
				name = "Access to " + v.Type
				for _, b := range v.Id {
					value += "<#" + b + ">" + "\n"
				}
			case "role":
				name = "Grant " + v.Type
				for _, b := range v.Id {
					value += "<@&" + b + ">" + "\n"
				}
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   name,
				Value:  value,
				Inline: false,
			})
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:  "Your benefit",
			Fields: fields,
			Color:  0x00ff00,
		})
	}
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &resultMsg,
		Embeds:  &embeds,
	})
}
