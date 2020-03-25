package discord

import (
	"bufio"
	"log"
	"os"
	"strings"

	discordgo "github.com/bwmarrin/discordgo"
)

var emoji = []string{}
var roleMap = map[string]string{}
var rolesMessageID string

// LoadRoleIDs loads roleIDs from files.
func LoadRoleIDs() {

	// Check for the environment ROLES_MESSAGE variable.
	if envRolesMessageID, exists := os.LookupEnv("ROLES_MESSAGE"); exists {
		rolesMessageID = envRolesMessageID
	}

	// Open roles.txt and read lines into an array.
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

	// Load roles via emojis into the map.
	for _, line := range lines {
		parts := strings.Split(line, "#")
		emoji = append(emoji, parts[0])
		roleMap[parts[0]] = parts[1]
	}
}

// Roles creates a designated reaction message for adding roles.
func Roles(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Create roles list and instructions message and send it.
	str := "**Roles**\n\n" +
		"Give yourself applicable roles by reacting to this message.\n\n" +
		"> add reaction: assign role\n" +
		"> remove reaction: unassign role\n\n"
	for _, r := range emoji {
		role, _ := s.State.Role(guildID, roleMap[r])
		str += r + ": " + role.Name + "\n"
	}
	message, _ := s.ChannelMessageSend(m.ChannelID, str)

	// React to with all role-emoji pairs.
	for _, r := range emoji {
		s.MessageReactionAdd(message.ChannelID, message.ID, r)
	}

	// Set messageID as the roles message.
	RewriteRolesMessageID(message.ID)
	rolesMessageID = message.ID
}

// Role add/remove helper.
func Role(s *discordgo.Session, roles []string,
	authorID, rID, msgID string, addRole bool) {

	// Return if the messageID is wrong.
	if len(rolesMessageID) == 0 || msgID != rolesMessageID {
		return
	}

	// Try adding the role if that's how this function is called.
	if addRole {
		if err := s.GuildMemberRoleAdd(guildID, authorID, rID); err != nil {
			log.Fatalln(err)
		}
		log.Printf("Adding role %s to %s.\n", rID, authorID)
		return
	}

	// Try removing the role if that's how this function is called.
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
