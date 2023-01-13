package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = true
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "bind",
		Description: "bind user dms to a channel",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "id of user to bind",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "channel",
				Description: "channel to bind user to",
				Required:    true,
			},
		},
	},
	{
		Name:        "reset",
		Description: "reset user binding to default channel",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "id of user to reset bindings for",
				Required:    true,
			},
		},
	},
	{
		Name:        "savebinds",
		Description: "save all bindings to disk",
	},
	{
		Name:        "listbinds",
		Description: "print all binds to stdout",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"bind": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.User.ID != ownerID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot use this command!",
				},
			})
		} else {
			bindUS := i.ApplicationCommandData().Options[0].StringValue()
			bindCID := i.ApplicationCommandData().Options[1].StringValue()
			//fmt.Println(bindUS)
			if find_bind(bindUS) == "" {
				create_bind(bindUS, bindCID)
			}
			modify_bind(bindUS, bindCID)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Bound User %s to Channel %s", bindUS, bindCID),
				},
			})
		}
	},
	"reset": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.User.ID != ownerID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot use this command!",
				},
			})
		} else {
			bindUS := i.ApplicationCommandData().Options[0].StringValue()
			//fmt.Println(bindUS)
			if find_bind(bindUS) == "" {
				create_bind(bindUS, defaultChannel)
			}
			modify_bind(bindUS, defaultChannel)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Reset user %s", bindUS),
				},
			})
		}
	},
	"savebinds": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.User.ID != ownerID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot use this command!",
				},
			})
		} else {
			save_bindings()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Saved Binds to Disk",
				},
			})
		}
	},
	"listbinds": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.User.ID != ownerID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot use this command!",
				},
			})
		} else {
			fmt.Println(bindings)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Bindings printed to stdout",
				},
			})
		}
	},
}
