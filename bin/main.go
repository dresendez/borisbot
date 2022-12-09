package main

import (
	"database/sql"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-sql-driver/mysql"
)

// Variables used for command line parameters
var (
	Token string
)
var (
	DBUSER string
)
var (
	DBPASS string
)

var db *sql.DB
var buffer = make([][]byte, 0)

var borisResponsesSlice = []string{"I'm Inveencible!",
	"The Americans are slugheads. They'll never detect me.",
	"Nobody screws with Boris Grishenko.",
	"Natalya. Shh, shh, it's me. It's Boris. It's Boris. It's Boris. Hello.",
	"I'm going for a cigarette.",
	"She's a moron. A second level programmer.  She works on the guidance system. She doesn't even have access to the filing codes. ",
	" [shaking a malfunctioning computer monitor vigorously]  Speak to me! ",
	"https://tenor.com/view/computer-angry-tech-rage-techno-rage-anger-gif-25680709"}

var franResponseSlice = []string{"Where the FUCK is my team",
	"No shot I lose this...",
	"No shot I'm the only one that landed Control!",
	"He doesn't really wanna be here, he's just pushin by himself..... I'm foaming",
	"I can't stand this guy",
	"*pushes alone and dies* Where is everybody",
	"no shot im the only one here!",
	"Is it boat gaming time yet...",
	"https://tenor.com/view/groundhog-chicken-chalken-chalk-gh-gif-21808428",
	"oop.. game crashed. :)",
	"I hate it here..."}

var peekResponseSlice = []string{"No shot that just happened to you like that...",
	"Damage check?",
	"Score check?",
	"I just got a new gaming chair too",
	"Not the teammate we want but the teammate we need...",
	"I just picked up a Welgun so you know shit gunna go well",
	"You good teammate?",
	"Purple gimme your money purple",
	"https://tenor.com/view/eric-andre-eric-andre-eric-andre-show-nightmare-gif-25187817",
	"https://tenor.com/view/ramsey-jalen-jaguars-jags-whatever-gif-10769417",
	"https://cdn.discordapp.com/attachments/773763353319440395/1047732590863454268/image.png"}

var rp1ResponseSlice = []string{"While yall were <insert childish activity>, I was <insert degenerate activity that rhymes>",
	"Can you-....",
	"Score check?",
	"!fran",
	"Rrrrrrrrfffuck",
	"https://tenor.com/view/trigg-tmnt-gif-9342541",
	"https://tenor.com/view/when-gif-20190303",
	"Hes just stirrin the pot... like a little witch :mage ",
	"You know how he gets..."}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&DBUSER, "u", "", "DB Username")
	flag.StringVar(&DBPASS, "p", "", "DB Pass")
	flag.Parse()
}

func main() {

	// Load the sound file.
	err := loadSound()
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		fmt.Println("Please copy $GOPATH/src/github.com/bwmarrin/examples/airhorn/airhorn.dca to this directory.")
		return
	}

	// Capture connection properties.
	cfg := mysql.Config{
		User:   DBUSER,
		Passwd: DBPASS,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "botdb",
	}
	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println(pingErr)
	}
	fmt.Println("Connected!")

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
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	// if m.Author.ID == s.State.User.ID {
	// 	return
	// }

	// check if the message is "!airhorn"
	if strings.HasPrefix(m.Content, "!airhorn") {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, g.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error playing sound:", err)
				}

				return
			}
		}
	}
	//TODO: refactor to map lookup, if successful, get random message from slice pertaining to map entry

	if m.Content == "!boris" {
		// send a random response
		message := fmt.Sprint(borisResponsesSlice[rand.Intn(len(borisResponsesSlice))])
		s.ChannelMessageSend(m.ChannelID, message)

	}
	if m.Content == "!fran" {
		// send a random response
		message := fmt.Sprint(franResponseSlice[rand.Intn(len(franResponseSlice))])
		s.ChannelMessageSend(m.ChannelID, message)

	}
	if m.Content == "!peek" {
		// send a random response
		message := fmt.Sprint(peekResponseSlice[rand.Intn(len(peekResponseSlice))])
		s.ChannelMessageSend(m.ChannelID, message)

	}
	if m.Content == "!rp1" {
		// send a random response
		message := fmt.Sprint(rp1ResponseSlice[rand.Intn(len(rp1ResponseSlice))])
		s.ChannelMessageSend(m.ChannelID, message)

	}

	if m.Content == "!commands" {
		// send a random response
		s.ChannelMessageSend(m.ChannelID, "Available commands: !boris")
	}

}

// loadSound attempts to load an encoded sound file from disk.
func loadSound() error {

	file, err := os.Open("airhorn.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
