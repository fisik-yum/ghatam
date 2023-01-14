/*
   kavaca- a discord bot that acts as a mail forwarder
   Copyright (C) 2021  fisik_yum

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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
