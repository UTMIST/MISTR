package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
	gitlab "github.com/xanzy/go-gitlab"
)

const guildID = "673778422291628033"
const discordTokenEnv = "DISCORD_BOT_TOKEN"
const gitlabTokenEnv = "GITLAB_TOKEN"
const projectIDEnv = "PROJECT_ID"

var prefixes = []string {"mista! ", "m! "}

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

func gitlabPagesUpdate() string {
	var reply string
	log.Println("Running CI for GitLab Pages...")

	// Look for the two environment variables.
	token, exists := os.LookupEnv(gitlabTokenEnv)
	if !exists {
		reply = "No discord bot token found."
		log.Println(reply)
		return reply
	}
	projectID, exists := os.LookupEnv(projectIDEnv)
	if !exists {
		reply = "No GitLab project ID found."
		log.Println(reply)
		return reply
	}

	// Set up GL Client
	git := gitlab.NewClient(nil, token)
	git.SetBaseURL("https://gitlab.com/api/v4")

	// Search for the most recent successful pipeline on master.
	listOpts := &gitlab.ListProjectPipelinesOptions{
		Status: gitlab.BuildState(gitlab.Success),
		Ref:    gitlab.String("master"),
	}
	pipelines, _, err := git.Pipelines.ListProjectPipelines(projectID, listOpts)
	if err != nil || len(pipelines) == 0 {
		log.Fatal(err)
	}
	pipeline := pipelines[0]
	log.Printf("Found successful pipeline: %d\n", pipeline.ID)


	// Get variables of last successful pipeline on master.
	vars, _, err := git.Pipelines.GetPipelineVariables(projectID, pipeline.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Use variables to run new pipeline.
	opt := &gitlab.CreatePipelineOptions{
		Ref:       gitlab.String("master"),
		Variables: vars,
	}
	newPipeline, _, err := git.Pipelines.CreatePipeline(projectID, opt)
	if err != nil {
		log.Fatal(err)
	}
	reply = fmt.Sprintf("Successfully rerun pipeline: %d as %d\n", pipeline.ID, newPipeline.ID)
	log.Print(reply)
	return reply
}

// 'On bot is ready' event
func ready(s *discord.Session, r *discord.Ready) {
	s.UpdateStatus(0, "defragmenting disk...")
}

// Message created event
func messageCreate(s *discord.Session, m *discord.MessageCreate) {
	var message string
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.GuildID != guildID {
		return
	}
	for _, p := range prefixes {
		if strings.HasPrefix(m.Content, p) {
			message = strings.TrimPrefix(m.Content, p)
		}
	}
	// message exists iff prefix was detected
	if len(message) == 0 {
		return}

	switch message {
	case "glp":
		reply := gitlabPagesUpdate()
		s.ChannelMessageSend(m.ChannelID, reply)
	}

}
