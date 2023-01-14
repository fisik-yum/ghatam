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

/*
This is a somewhat heavily modified version of the discordgo examples for slash commands
Major changes include initialization code, message handling, and message parsing for replies and commands, since slash commands make getting args a lot easier.
*/
package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// variables found in config.json, which needs to exist
var (
	token          string
	ownerID        string
	defaultChannel string
	bindings       []binding
)

var s *discordgo.Session

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

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	if defaultChannel == "" {
		dc, err := s.UserChannelCreate(ownerID)
		check(err)
		defaultChannel = dc.ID
	}

}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	s.ShouldReconnectOnError = true
	s.Identify.Intents = 12800

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		log.Printf("Registering command: %v", v.Name)
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	log.Println("Adding forwarding Handler")
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if (m.Author.ID != s.State.User.ID) && m.Author.ID != ownerID {
			c, err := s.Channel(m.ChannelID)
			check(err)
			if c.Type != discordgo.ChannelTypeDM { //only allow DM channels
				return
			}
			if find_bind(m.Author.ID) == "" {
				create_bind(m.Author.ID, defaultChannel)
			}
			cID := find_bind(m.Author.ID)
			s.ChannelMessageSendEmbed(cID, &discordgo.MessageEmbed{Fields: []*discordgo.MessageEmbedField{{Name: (m.Author.Username + "#" + m.Author.Discriminator), Value: (m.Content), Inline: true}}, Author: &discordgo.MessageEmbedAuthor{Name: m.Author.ID}})

			return
		} else if (m.Author.ID == ownerID) && (m.Message.Type == discordgo.MessageTypeReply) && (m.ReferencedMessage.Author.ID == s.State.User.ID) {
			refMes := m.ReferencedMessage
			if len(refMes.Embeds) < 1 {
				return
			}
			mesAuthID := refMes.Embeds[0].Author.Name
			fchannel, err := s.UserChannelCreate(mesAuthID)
			check(err)
			s.ChannelMessageSend(fchannel.ID, m.Content)

			check(err)
		}

	})

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("Saving channel bindings")
	save_bindings()

	log.Println("Gracefully shutting down.")
}
