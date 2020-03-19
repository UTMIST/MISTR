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
