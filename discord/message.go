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

	var message string
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.GuildID != guildID {
		return
	}

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

	if hostname, err := os.Hostname(); message == " host" && err == nil {
		log.Println("Printing hostname.")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I'm running on %s.", hostname))
		return
	}
	if !correctChannel(m) {
		return
	}

	// Switch on message for reply.
	switch message {
	case " roles":
		if rolesChannel, exists := os.LookupEnv("ROLES_CHANNEL"); exists && rolesChannel == m.ChannelID {
			Roles(s, m)
		}
	case " update":
		if updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL"); exists && updateChannel == m.ChannelID {
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
func MessageAddReact(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}
	member, _ := s.GuildMember(guildID, m.UserID)
	Role(s, member.Roles, m.UserID, roleMap[m.Emoji.Name], m.MessageID, false)
}

// MessageRemoveReact is the handler for when someone removes a react.
func MessageRemoveReact(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.UserID == s.State.User.ID {
		return
	}
	member, _ := s.GuildMember(guildID, m.UserID)
	Role(s, member.Roles, m.UserID, roleMap[m.Emoji.Name], m.MessageID, true)
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

func correctChannel(m *discordgo.MessageCreate) bool {
	env, envExists := os.LookupEnv("ENVIRONMENT")
	devChannel, devChExists := os.LookupEnv("DEV_CHANNEL")
	if !envExists || !devChExists {
		log.Println("Missing env variables for DEV_CHANNEL and/or ENVIRONMENT")
		return false
	}
	inDev := env == "DEV" && m.ChannelID == devChannel
	inProd := env != "DEV" && m.ChannelID != devChannel
	if !inDev && !inProd {
		log.Printf("Environment mismatch for %s.\n", env)
		return false
	}

	return true
}
