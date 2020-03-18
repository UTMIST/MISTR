package discord

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	discordgo "github.com/bwmarrin/discordgo"
)

var prefixes = []string{"mista!", "m!"}

const guildID = "673778422291628033"

// Ready handler for when MISTA is active.
func Ready(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateStatus(0, "training models...")
}

// Help returns MISTA's manual.
func Help() string {

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

	return line
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

// RewriteRolesMessageID rewrites .env with the new channel ID for roles.
func RewriteRolesMessageID(messageID string) {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}

	lines := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "ROLES_MESSAGE") >= 0 {
			line = fmt.Sprintf("%s=%s", line[:strings.Index(line, "=")], messageID)
		}
		lines = append(lines, line)
	}
	file.Close()

	file, err = os.Create(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, line := range lines {
		fmt.Fprintln(file, line)
	}
}
