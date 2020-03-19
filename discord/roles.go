package discord

import (
	"bufio"
	"log"
	"os"
	"strings"

	discordgo "github.com/bwmarrin/discordgo"
)

var rolesMessageID string

var emoji = []string{}

var roleMap = map[string]string{}

// LoadRoleIDs loads roleIDs from files.
func LoadRoleIDs() {

	if envRolesMessageID, exists := os.LookupEnv("ROLES_MESSAGE"); exists {
		rolesMessageID = envRolesMessageID
	}

	file, err := os.Open("roles.txt")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, line := range lines {
		parts := strings.Split(line, "#")
		emoji = append(emoji, parts[0])
		roleMap[parts[0]] = parts[1]
	}
}

// Roles creates a designated reaction message for adding roles.
func Roles(s *discordgo.Session, m *discordgo.MessageCreate) {
	str := "**Roles**\n\n" +
		"Please give yourself all applicable roles by reacting to this message.\n\n" +
		"> add reaction: assign role\n" +
		"> remove reaction: unassign role\n\n"

	for _, r := range emoji {
		role, _ := s.State.Role(guildID, roleMap[r])
		str += r + ": " + role.Name + "\n"
	}

	message, _ := s.ChannelMessageSend(m.ChannelID, str)
	for _, r := range emoji {
		s.MessageReactionAdd(message.ChannelID, message.ID, r)
	}

	RewriteRolesMessageID(message.ID)

	rolesMessageID = message.ID
}

// Role add/remove helper.
func Role(s *discordgo.Session, roles []string, authorID, rID, msgID string, addRole bool) {
	if len(rolesMessageID) == 0 || msgID != rolesMessageID {
		return
	}

	if !addRole {
		if err := s.GuildMemberRoleAdd(guildID, authorID, rID); err != nil {
			log.Fatalln(err)
		}
		log.Printf("Adding role %s to %s.\n", rID, authorID)
		return
	}

	for _, ID := range roles {
		if ID != rID {
			continue
		}

		if err := s.GuildMemberRoleRemove(guildID, authorID, rID); err != nil {
			log.Fatalln(err)
		}

		log.Printf("Removing role %s from %s.\n", rID, authorID)
	}
}
