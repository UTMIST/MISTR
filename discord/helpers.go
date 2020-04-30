package discord

import (
	"bufio"
	"fmt"
	"log"
	"os"

	discordgo "github.com/bwmarrin/discordgo"
)

// Return MISTA's manual.
func help(s *discordgo.Session, m *discordgo.MessageCreate) {

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
func correctChannel(m *discordgo.MessageCreate) (bool, bool) {

	// Check that environment variables.
	env, envExists := os.LookupEnv("ENVIRONMENT")
	devChannel, devChExists := os.LookupEnv("DEV_CHANNEL")
	if !envExists || !devChExists {
		log.Println("Missing env variables DEV_CHANNEL and/or ENVIRONMENT")
		return false, false
	}

	// Check if channel matches environment.
	inDev := env == "DEV" && m.ChannelID == devChannel
	inProd := env != "DEV" && m.ChannelID != devChannel
	if !inDev && !inProd {
		log.Printf("Environment mismatch for %s.\n", env)
		return false, false
	}

	return true, inDev
}
