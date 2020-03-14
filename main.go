package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
	gitlab "gitlab.com/utmist/mista/gitlab"
)

const guildID = "673778422291628033"
const discordTokenEnv = "DISCORD_BOT_TOKEN"

var rolesMessageID string
var prefixes = []string{"mista! ", "m! "}
var emojis = []string{
	"3️⃣",
	"2️⃣",
	"1️⃣",
	"🦝",
	"🦅",
	"🍁",
	"🎲",
	"⚖️",
	"🔎",
	"🧮",
	"💻",
	"🧠",
	"🔌",
	"👨‍🔬",
}
var roleIDs = []string{
	"678759295998885888", // :three: 2T3
	"678759268845092864", // :two: 2T2
	"678759179967528961", // :one: 2T1
	"678847616691208203", // :raccoon: UTSC
	"678847667190890496", // :eagle: UTM
	"683102967582294037", // :maple_leaf: UTSG
	"678759570776260624", // :game_die: Stats
	"678759600413212702", // :scales: Phys
	"678842446683176980", // :mag_right: Phil
	"678759542808379415", // :abacus: Math
	"678759477008138252", // :computer: CompSci
	"678764218094321685", // :brain: CogSci
	"678759424705429518", // :electric_plug: ECE
	"678759380589740040", // :man_scientist: EngSci
}
var roleMap = make(map[string]string)

// Role add/remove helper
func role(s *discord.Session, roles []string, authorID string, rID string, hasRole bool, msgID string) {
	if len(rolesMessageID) == 0 || msgID != rolesMessageID {
		return
	}
	if hasRole {
		for _, ID := range roles {
			if ID == rID {
				err := s.GuildMemberRoleRemove(guildID, authorID, rID)
				if err != nil {
					log.Fatal(err)
				}

			}
		}
	} else {
		for _, ID := range roles {
			if ID == rID {
				return
			}
		}
		err := s.GuildMemberRoleAdd(guildID, authorID, rID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

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

// Creates a designated reaction message for adding roles
func deploy(s *discord.Session, m *discord.MessageCreate) {

	str := "**Incoming Members**\n\n" +
		"Please give yourself all applicable roles by reacting to this message.\n\n" +
		"> add reaction => assign role\n" +
		"> remove reaction => unassign role\n\n"
	for _, r := range emojis {
		role, _ := s.State.Role(guildID, roleMap[r])
		str += r + ": " + role.Name + "\n"
	}
	message, _ := s.ChannelMessageSend(m.ChannelID, str)
	for _, r := range emojis {
		s.MessageReactionAdd(message.ChannelID, message.ID, r)
	}
	file, err := os.Create("roles-id.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(message.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	rolesMessageID = message.ID
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
	switch message {
	case "update":
		s.ChannelMessageSend(m.ChannelID, gitlab.PagesUpdate())
	case "deploy":
		deploy(s, m)
	case "help":
		s.ChannelMessageSend(m.ChannelID, help())
	}

}

// Message add reaction handler
func messageAddReact(s *discord.Session, m *discord.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}
	member, _ := s.GuildMember(guildID, m.UserID)
	role(s, member.Roles, m.UserID, roleMap[m.Emoji.Name], false, m.MessageID)
}

// Message remove reaction handler
func messageRemoveReact(s *discord.Session, m *discord.MessageReactionRemove) {
	if m.UserID == s.State.User.ID {
		return
	}
	member, _ := s.GuildMember(guildID, m.UserID)
	role(s, member.Roles, m.UserID, roleMap[m.Emoji.Name], true, m.MessageID)
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

	// Building reaction: roleID map
	for i := 0; i < len(roleIDs); i++ {
		roleMap[emojis[i]] = roleIDs[i]
	}

	r, err := ioutil.ReadFile("roles-id.txt")
	if err != nil {
		log.Println(err)
	} else {
		rolesMessageID = string(r)
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
	dg.AddHandler(messageAddReact)
	dg.AddHandler(messageRemoveReact)

	// Open client and run on a loop.
	if err = dg.Open(); err != nil {
		log.Fatalln(err)
	}
	log.Printf("MISTA (ID: %s) is running...\n", botID)
	<-make(chan struct{})

	return
}
