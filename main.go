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

// this is a modified version of ping pong from discordgo examples

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// variables found in config.json, which needs to exist
var (
	token          string
	ownerID        string
	prefix         string
	defaultChannel string
	bindings       []binding
)

func init() {
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		panic("config.json is missing")
	}
	_, err = os.Stat("bindings.json")
	if os.IsNotExist(err) {
		os.Create("bindings.json")
	}
	load_bindings()
	save_bindings()
	read_config()
}

func main() {

	dg, err := discordgo.New("Bot " + token)
	check(err)
	dg.AddHandler(messageCreate)
	dg.ShouldReconnectOnError = true
	dg.Identify.Intents = 12800
	err = dg.Open()
	check(err)
	if defaultChannel == "" {
		dc, err := dg.UserChannelCreate(ownerID)
		check(err)
		defaultChannel = dc.ID
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	save_bindings()
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	//basic functionality
	if s.State.User.ID == m.Author.ID {
		return
	} else if m.Author.ID != ownerID {
		c, err := s.Channel(m.ChannelID)
		check(err)
		if c.Type != discordgo.ChannelTypeDM { //only allow DM channels
			return
		}
		if find_bind(m.Author.ID) == "" {
			create_bind(m.Author.ID, defaultChannel)
		}
		cID := find_bind(m.Author.ID)
		s.ChannelMessageSend(cID, (m.Author.ID + " **" + m.Author.Username + "#" + m.Author.Discriminator + ":** " + m.Content))
		return
	}

	//command handling, modular prefixes, accessible to only the owner
	cmd := trim_index(m.Content, 0)
	if m.Author.ID == ownerID && m.Type != discordgo.MessageTypeReply {
		switch strings.HasPrefix(cmd, prefix) {

		case strings.HasPrefix(cmd, prefix+"bind "): //also modifies existing binds
			bindUS := trim_index(m.Content, 1)
			bindCID := trim_index(m.Content, 2)
			if find_bind(bindUS) == "" {
				create_bind(bindUS, bindCID)
			}
			modify_bind(bindUS, bindCID)
			s.ChannelMessageSend(m.ChannelID, "rebound user "+bindUS)
			return

		case strings.HasPrefix(cmd, prefix+"listbinds"):
			fmt.Println(bindings)
			save_bindings()
			return

		case strings.HasPrefix(cmd, prefix+"reset "):
			bindUS := trim_index(m.Content, 1)
			if find_bind(bindUS) == "" {
				create_bind(bindUS, defaultChannel)
			}
			modify_bind(bindUS, defaultChannel)
			s.ChannelMessageSend(m.ChannelID, "reset bind for user "+bindUS)
			return

		case strings.HasPrefix(cmd, prefix+"savebinds"):
			save_bindings()
			return

		case strings.HasPrefix(cmd, prefix+"info"):
			s.ChannelMessageSend(m.ChannelID, ("` kavaca built with discordgo " + discordgo.VERSION + "`"))
			return

		default:
			s.ChannelMessageSend(m.ChannelID, "Invalid Command")
			return
		}
	} else if m.Author.ID == ownerID && m.Type == discordgo.MessageTypeReply {
		im, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
		check(err)
		truncate := len(im.Author.ID+" **"+im.Author.Username+"#"+im.Author.Discriminator+":** ") - 1
		if (len(im.Content) < truncate-1) || im.Author.ID != s.State.User.ID {
			s.ChannelMessageSend(m.ChannelID, "An error occurred")
			return
		}
		originalID := im.Content[:18]
		c, err := s.UserChannelCreate(originalID)
		check(err)
		rmessagecontent := fmt.Sprintf("`Reply to message sent @ %s :` *%s*\n%s", string(im.Timestamp)[:19], im.Content[truncate:], m.Content)
		s.ChannelMessageSend(c.ID, rmessagecontent)
		return
	}
}
