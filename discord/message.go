package discord

import (
	"fmt"
	"log"
	"os"
	"strings"

	discordgo "github.com/bwmarrin/discordgo"
	"gitlab.com/utmist/mista/gitlab"
)

// MessageCreate is the handler for when a message is created.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// If the bot and/or server don't match.
	if m.Author.ID == s.State.User.ID || m.GuildID != guildID {
		return
	}

	var message string

	// Search for a matching prefix for alias.
	for i, p := range prefixes {
		if strings.HasPrefix(m.Content, p) {
			message = strings.TrimPrefix(m.Content, p)
			break
		}

		// If none of the prefixes match, it's not a message for MISTA.
		if i == len(prefixes)-1 {
			return
		}
	}

	// Return hostname.
	if hostname, err := os.Hostname(); message == " host" && err == nil {
		log.Println("Printing hostname.")
		s.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("I'm running on %s.", hostname))
		return
	}

	// Other than hostname, return if the channel doesn't match environment.
	correct, inDev := correctChannel(m)
	if !correct {
		return
	}

	// Switch on message for reply.
	switch message {
	case " flush":
		updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL")
		if inDev || exists && (m.ChannelID == updateChannel) {
			s.ChannelMessageSend(m.ChannelID, gitlab.PagesFlush())
		}
	case " roles":
		rolesChannel, exists := os.LookupEnv("ROLES_CHANNEL")
		if exists && m.ChannelID == rolesChannel {
			Roles(s, m)
		}
	case " update":
		updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL")
		if inDev || exists && (m.ChannelID == updateChannel) {
			s.ChannelMessageSend(m.ChannelID, gitlab.PagesUpdate())
		}
	case " restart":
		updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL")
		if inDev || exists && (m.ChannelID == updateChannel) {
			s.ChannelMessageSend(m.ChannelID, "I'm restarting and getting some upgrades :)")
			log.Println("Restarting...")
			os.Exit(0)
		}

	case "":
		help(s, m)
	case " help":
		help(s, m)
	case " manual":
		help(s, m)
	}
}

// MessageAddReact is the handler for when someone adds a react.
func MessageAddReact(s *discordgo.Session,
	m *discordgo.MessageReactionAdd) {

	// Don't act on the bot's reacts.
	if m.UserID == s.State.User.ID {
		return
	}

	member, _ := s.GuildMember(guildID, m.UserID)
	Role(s, member.Roles, m.UserID, roleMap[m.Emoji.Name], m.MessageID, true)
}

// MessageRemoveReact is the handler for when someone removes a react.
func MessageRemoveReact(s *discordgo.Session,
	m *discordgo.MessageReactionRemove) {

	// Don't act on the bot's reacts.
	if m.UserID == s.State.User.ID {
		return
	}

	member, _ := s.GuildMember(guildID, m.UserID)
	Role(s, member.Roles, m.UserID, roleMap[m.Emoji.Name], m.MessageID, false)
}
