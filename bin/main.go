package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

var borisResponsesSlice = []string{"I'm Inveencible!",
	"The Americans are slugheads. They'll never detect me.",
	"Nobody screws with Boris Grishenko.",
	"Natalya. Shh, shh, it's me. It's Boris. It's Boris. It's Boris. Hello.",
	"I'm going for a cigarette.",
	"She's a moron. A second level programmer.  She works on the guidance system. She doesn't even have access to the filing codes. ",
	" [shaking a malfunctioning computer monitor vigorously]  Speak to me! "}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Got msg: ", m.Content)
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	message := fmt.Sprint(borisResponsesSlice[rand.Intn(len(borisResponsesSlice))])
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!boris" {
		// send a random response
		s.ChannelMessageSend(m.ChannelID, message)

	}

	if m.Content == "!commands" {
		// send a random response
		s.ChannelMessageSend(m.ChannelID, "Available commands: !boris")
	}

}
