package discord

import (
	"fmt"
	"log"
	"os"
	"strings"

	discordgo "github.com/bwmarrin/discordgo"
	"gitlab.com/utmist/mistr/gitlab"
	"gitlab.com/utmist/mistr/update"
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

	// clean up the message
	message = strings.TrimSpace(message)

	// parse messages into each substring for hash
	subMessages := strings.Split(message, " ")

	// handle prefix only
	if len(subMessages) == 0 {
		help(s, m)
		return
	}

	// Switch on message for reply.
	switch subMessages[0] {
	case "flush":
		updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL")
		if inDev || exists && (m.ChannelID == updateChannel) {
			s.ChannelMessageSend(m.ChannelID, gitlab.PagesFlush())
		}
	case "update":
		updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL")
		if inDev || exists && (m.ChannelID == updateChannel) {
			s.ChannelMessageSend(m.ChannelID, gitlab.PagesUpdate())
		}
	case "restart":
		updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL")
		if inDev || exists && (m.ChannelID == updateChannel) {
			if update.IsUpdated() {
				s.ChannelMessageSend(m.ChannelID, "I'm already up to date :)")
				return
			}

			reply := "I'm restarting and getting some upgrades :D"
			if URL, exists := os.LookupEnv("REPO_URL"); exists {
				reply = fmt.Sprintf("%s; see %s.", reply, URL)
			}
			s.ChannelMessageSend(m.ChannelID, reply)
			log.Println(reply)
			os.Exit(0)
		}
	case "hash":
		botSandboxChannel, exists := os.LookupEnv("DEV_CHANNEL")
		if exists && (m.ChannelID == botSandboxChannel) {
			if len(subMessages) > 2 {
				s.ChannelMessageSend(m.ChannelID, "Hash input should only contain one string, try again :P")
				return
			}
			hashString(s, m, subMessages[1])
		}
	default:
		help(s, m)
	}
}
