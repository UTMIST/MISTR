package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	gitlab "gitlab.com/utmist/mista/gitlab"

	discord "github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
)

const guildID = "673778422291628033"
const discordTokenEnv = "DISCORD_BOT_TOKEN"

var prefixes = []string{"mista!", "m!"}

// Bot is ready handler
func ready(s *discord.Session, r *discord.Ready) {
	s.UpdateStatus(0, "training models...")
}

// Return MISTA's manual.
func help() string {

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

	return line
}

// Message created handler
func messageCreate(s *discord.Session, m *discord.MessageCreate) {
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
	var reply string
	switch message {
	case " update":
		reply = gitlab.PagesUpdate()
	default:
		reply = help()
	}

	s.ChannelMessageSend(m.ChannelID, reply)

}

// Message reaction handler
func messageReact(s *discord.Session, m discord.MessageReactionAdd) {
	// user, _ := s.User(m.UserID)
	switch m.Emoji.Name {
	case ":one:":

	}
}

func main() {
	if err := dotenv.Load(); err != nil {
		log.Println("Could not load .env")
	}

	// Load bot token.
	token := os.Getenv(discordTokenEnv)
	dg, err := discord.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Get the bot's user.
	u, err := dg.User("@me")
	if err != nil {
		log.Println(err)
	}
	botID := u.ID

	// Add Handlers
	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	// Open client and run on a loop.
	if err = dg.Open(); err != nil {
		log.Fatalln(err)
	}
	log.Printf("MISTA (ID: %s) is running...\n", botID)
	<-make(chan struct{})

	return
}
