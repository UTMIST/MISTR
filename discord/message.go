package discord

import (
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

	// Switch on message for reply.
	switch message {
	case " update":
		if updateChannel, exists := os.LookupEnv("UPDATE_CHANNEL"); exists && updateChannel == m.ChannelID {
			s.ChannelMessageSend(m.ChannelID, gitlab.PagesUpdate())
		}
	case " roles":
		if rolesChannel, exists := os.LookupEnv("ROLES_CHANNEL"); exists && rolesChannel == m.ChannelID {
			Roles(s, m)
		}
	default:
		s.ChannelMessageSend(m.ChannelID, Help())
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
