package discord

import (
	"bufio"
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
	if !correctChannel(m) {
		return
	}

	// Switch on message for reply.
	switch message {
	case " roles":
		rolesChannel, exists := os.LookupEnv("ROLES_CHANNEL")
		if exists && rolesChannel == m.ChannelID {
			Roles(s, m)
		}
	case " update":
		updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL")
		if exists && updateChannel == m.ChannelID {
			s.ChannelMessageSend(m.ChannelID, gitlab.PagesUpdate())
		}
	case "":
		Help(s, m)
	case " help":
		Help(s, m)
	case " manual":
		Help(s, m)
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

// Help returns MISTA's manual.
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Open and defer-close the file.
	file, err := os.Open("manual.md")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Build the multi-line manual.
	line := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprintf("%s\n%s", line, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Printing manual.")
	s.ChannelMessageSend(m.ChannelID, line)
}

// Return whether the channel matches the environment.
func correctChannel(m *discordgo.MessageCreate) bool {

	// Check that environment variables.
	env, envExists := os.LookupEnv("ENVIRONMENT")
	devChannel, devChExists := os.LookupEnv("DEV_CHANNEL")
	if !envExists || !devChExists {
		log.Println("Missing env variables DEV_CHANNEL and/or ENVIRONMENT")
		return false
	}

	// Check if channel matches environment.
	inDev := env == "DEV" && m.ChannelID == devChannel
	inProd := env != "DEV" && m.ChannelID != devChannel
	if !inDev && !inProd {
		log.Printf("Environment mismatch for %s.\n", env)
		return false
	}

	return true
}
