package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
	discord "gitlab.com/utmist/mistr/discord"
)

const discordTokenEnv = "DISCORD_BOT_TOKEN"

func main() {

	if err := dotenv.Load(); err != nil {
		log.Println("Could not load .env")
	}

	// Load bot token.
	token := os.Getenv(discordTokenEnv)
	dg, err := discordgo.New("Bot " + token)
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
	dg.AddHandler(discord.Ready)
	dg.AddHandler(discord.MessageCreate)

	// Open client and run on a loop.
	if err = dg.Open(); err != nil {
		log.Fatalln(err)
	}
	log.Printf("MISTR (ID: %s) is running...\n", botID)
	<-make(chan struct{})

	return
}
