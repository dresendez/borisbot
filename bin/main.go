package main

import (
<<<<<<< Updated upstream
	"flag"
	"fmt"
	"math/rand"
=======
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
>>>>>>> Stashed changes
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
<<<<<<< Updated upstream
=======
	_ "github.com/mattn/go-sqlite3"
>>>>>>> Stashed changes
)

// Config holds the bot configuration
type Config struct {
	DiscordToken string `json:"discord_token"`
	AdminUserID  string `json:"admin_user_id"` // Discord User ID of the admin
}

// Quote represents a quote with metadata
type Quote struct {
	ID       int64
	UserID   string
	Text     string
	AddedBy  string
	AddedAt  time.Time
	Context  string
	UseCount int
	LastUsed *time.Time
}

// AuthorizedUser represents a user who can have quotes
type AuthorizedUser struct {
	ID       string // Discord User ID
	Name     string
	AddedBy  string
	AddedAt  time.Time
	IsActive bool
}

// BackupQuotes represents the structure of the quotes backup file
type BackupQuotes struct {
	Quotes map[string][]string `json:"quotes"`
}

// Variables used for command line parameters
var (
	Token string
<<<<<<< Updated upstream
)

var borisResponsesSlice = []string{"I'm Inveencible!",
	"The Americans are slugheads. They'll never detect me.",
	"Nobody screws with Boris Grishenko.",
	"Natalya. Shh, shh, it's me. It's Boris. It's Boris. It's Boris. Hello.",
	"I'm going for a cigarette.",
	"She's a moron. A second level programmer.  She works on the guidance system. She doesn't even have access to the filing codes. ",
	" [shaking a malfunctioning computer monitor vigorously]  Speak to me! "}

var franResponseSlice = []string{"Where the FUCK is my team",
	"No shot I lose this...",
	"No shot I'm the only one that landed Control!",
	"He doesn't really wanna be here, he's just pushin by himself..... I'm foaming",
	"I can't stand this guy",
	"*pushes alone and dies* Where is everybody",
	"no shot im the only one here!",
	"Is it boat gaming time yet..."}

var peekResponseSlice = []string{"No shot that just happened to you like that...",
	"Damage check?",
	"Score check?",
	"I just got a new gaming chair too",
	"Not the teammate we want but the teammate we need...",
	"I just picked up a Welgun so you know shit gunna go well",
	"You good teammate?",
	"Purple gimme your money purple",
	"https://tenor.com/view/eric-andre-eric-andre-eric-andre-show-nightmare-gif-25187817"}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
=======
	db    *sql.DB
)

// CommandHandler represents a function that handles a command
type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate) error

// Bot represents our discord bot
type Bot struct {
	commands map[string]CommandHandler
	config   Config
>>>>>>> Stashed changes
}

// NewBot creates a new instance of our bot
func NewBot(config Config) (*Bot, error) {
	bot := &Bot{
		commands: make(map[string]CommandHandler),
		config:   config,
	}

<<<<<<< Updated upstream
=======
	// Initialize database
	if err := bot.initDB(); err != nil {
		return nil, fmt.Errorf("error initializing database: %v", err)
	}

	// Register command handlers
	bot.commands["!quote"] = bot.handleQuote
	bot.commands["!addquote"] = bot.handleAddQuote
	bot.commands["!listquotes"] = bot.handleListQuotes
	bot.commands["!delquote"] = bot.handleDeleteQuote
	bot.commands["!adduser"] = bot.handleAddUser
	bot.commands["!listusers"] = bot.handleListUsers
	bot.commands["!commands"] = bot.handleCommands
	bot.commands["!context"] = bot.handleAddContext
	bot.commands["!backup"] = bot.handleBackup
	bot.commands["!restore"] = bot.handleRestore
	bot.commands["!initdb"] = bot.handleInitDB // New command for one-time initialization

	return bot, nil
}

func (b *Bot) initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "quotes.db")
	if err != nil {
		return err
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS authorized_users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			added_by TEXT NOT NULL,
			added_at DATETIME NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT 1
		);
		
		CREATE TABLE IF NOT EXISTS quotes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			text TEXT NOT NULL,
			added_by TEXT NOT NULL,
			added_at DATETIME NOT NULL,
			context TEXT,
			use_count INTEGER NOT NULL DEFAULT 0,
			last_used DATETIME,
			FOREIGN KEY (user_id) REFERENCES authorized_users(id)
		);
	`)
	return err
}

// isAdmin checks if the user is the admin
func (b *Bot) isAdmin(userID string) bool {
	return userID == b.config.AdminUserID
}

// isAuthorizedUser checks if the user is authorized
func (b *Bot) isAuthorizedUser(userID string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM authorized_users WHERE id = ? AND is_active = 1)", userID).Scan(&exists)
	return err == nil && exists
}

// handleAddUser adds a new authorized user (admin only)
func (b *Bot) handleAddUser(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !b.isAdmin(m.Author.ID) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Only the admin can add users")
		return err
	}

	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 3 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !adduser <user_id> <name>")
		return err
	}

	userID := parts[1]
	name := parts[2]

	_, err := db.Exec(
		"INSERT INTO authorized_users (id, name, added_by, added_at) VALUES (?, ?, ?, ?)",
		userID, name, m.Author.ID, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("error adding user: %v", err)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added user %s (%s)", name, userID))
	return err
}

// handleQuote returns a random quote for a user
func (b *Bot) handleQuote(s *discordgo.Session, m *discordgo.MessageCreate) error {
	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !quote <user>")
		return err
	}

	userID := parts[1]
	fmt.Printf("Fetching quote for user: %s\n", userID)

	// Check if user exists
	if !b.isAuthorizedUser(userID) {
		msg := fmt.Sprintf("User %s is not authorized. Available users can be listed with !listusers", userID)
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		return err
	}

	var (
		quote       Quote
		contextNull sql.NullString
	)

	err := db.QueryRow(`
		SELECT id, text, context, use_count 
		FROM quotes 
		WHERE user_id = ? 
		ORDER BY RANDOM() 
		LIMIT 1
	`, userID).Scan(&quote.ID, &quote.Text, &contextNull, &quote.UseCount)

	if err == sql.ErrNoRows {
		_, err := s.ChannelMessageSend(m.ChannelID, "No quotes found for this user")
		return err
	}
	if err != nil {
		return err
	}

	// Only use context if it's valid (not NULL)
	if contextNull.Valid {
		quote.Context = contextNull.String
	}

	fmt.Printf("Found quote %d for user %s\n", quote.ID, userID)

	// Update use count and last used
	_, err = db.Exec(
		"UPDATE quotes SET use_count = use_count + 1, last_used = ? WHERE id = ?",
		time.Now(), quote.ID,
	)
	if err != nil {
		return err
	}

	message := quote.Text
	if quote.Context != "" {
		message += fmt.Sprintf("\n\nContext: %s", quote.Context)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, message)
	return err
}

// handleAddQuote adds a new quote
func (b *Bot) handleAddQuote(s *discordgo.Session, m *discordgo.MessageCreate) error {
	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 3 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !addquote <user_id> <quote>")
		return err
	}

	userID := parts[1]
	quote := parts[2]

	// Check if the target user is authorized
	if !b.isAuthorizedUser(userID) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Cannot add quotes for unauthorized users")
		return err
	}

	_, err := db.Exec(
		"INSERT INTO quotes (user_id, text, added_by, added_at) VALUES (?, ?, ?, ?)",
		userID, quote, m.Author.ID, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("error adding quote: %v", err)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Quote added successfully")
	return err
}

// handleAddContext adds context to a quote
func (b *Bot) handleAddContext(s *discordgo.Session, m *discordgo.MessageCreate) error {
	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 3 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !context <quote_id> <context>")
		return err
	}

	var quoteID int64
	fmt.Sscanf(parts[1], "%d", &quoteID)
	context := parts[2]

	_, err := db.Exec("UPDATE quotes SET context = ? WHERE id = ?", context, quoteID)
	if err != nil {
		return fmt.Errorf("error adding context: %v", err)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Context added successfully")
	return err
}

// handleListQuotes lists all quotes for a user
func (b *Bot) handleListQuotes(s *discordgo.Session, m *discordgo.MessageCreate) error {
	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !listquotes <user_id>")
		return err
	}

	userID := parts[1]
	rows, err := db.Query(`
		SELECT id, text, context, use_count, added_at, added_by 
		FROM quotes 
		WHERE user_id = ?
		ORDER BY added_at DESC
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var quotes []Quote
	for rows.Next() {
		var (
			q           Quote
			contextNull sql.NullString
		)
		err := rows.Scan(&q.ID, &q.Text, &contextNull, &q.UseCount, &q.AddedAt, &q.AddedBy)
		if err != nil {
			return err
		}
		if contextNull.Valid {
			q.Context = contextNull.String
		}
		quotes = append(quotes, q)
	}

	if len(quotes) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "No quotes found for this user")
		return err
	}

	message := fmt.Sprintf("Quotes for user:\n")
	for _, q := range quotes {
		message += fmt.Sprintf("%d. %s", q.ID, q.Text)
		if q.Context != "" {
			message += fmt.Sprintf(" (Context: %s)", q.Context)
		}
		message += fmt.Sprintf(" [Used %d times, Added by %s on %s]\n",
			q.UseCount, q.AddedBy, q.AddedAt.Format("2006-01-02"))
	}

	_, err = s.ChannelMessageSend(m.ChannelID, message)
	return err
}

// handleListUsers lists all authorized users (admin only)
func (b *Bot) handleListUsers(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !b.isAdmin(m.Author.ID) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Only the admin can list users")
		return err
	}

	rows, err := db.Query("SELECT id, name, added_at FROM authorized_users WHERE is_active = 1")
	if err != nil {
		return err
	}
	defer rows.Close()

	var users []AuthorizedUser
	for rows.Next() {
		var u AuthorizedUser
		err := rows.Scan(&u.ID, &u.Name, &u.AddedAt)
		if err != nil {
			return err
		}
		users = append(users, u)
	}

	message := "Authorized Users:\n"
	for _, u := range users {
		message += fmt.Sprintf("- %s (%s) [Added: %s]\n",
			u.Name, u.ID, u.AddedAt.Format("2006-01-02"))
	}

	_, err = s.ChannelMessageSend(m.ChannelID, message)
	return err
}

// handleCommands shows available commands
func (b *Bot) handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) error {
	commands := make([]string, 0, len(b.commands))
	for cmd := range b.commands {
		commands = append(commands, cmd)
	}
	message := fmt.Sprintf("Available commands: %s", strings.Join(commands, ", "))
	_, err := s.ChannelMessageSend(m.ChannelID, message)
	return err
}

// handleDeleteQuote deletes a quote by ID
func (b *Bot) handleDeleteQuote(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !b.isAdmin(m.Author.ID) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Only the admin can delete quotes")
		return err
	}

	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !delquote <quote_id>")
		return err
	}

	var quoteID int64
	fmt.Sscanf(parts[1], "%d", &quoteID)

	result, err := db.Exec("DELETE FROM quotes WHERE id = ?", quoteID)
	if err != nil {
		return fmt.Errorf("error deleting quote: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		_, err = s.ChannelMessageSend(m.ChannelID, "No quote found with that ID")
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Deleted quote %d", quoteID))
	return err
}

// messageCreate handles incoming messages
func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Printf("Received message: %s\n", m.Content)

	// Split the message into command and arguments
	parts := strings.SplitN(m.Content, " ", 2)
	command := parts[0]

	if handler, exists := b.commands[command]; exists {
		fmt.Printf("Executing command: %s\n", command)
		if err := handler(s, m); err != nil {
			fmt.Printf("Error handling command %s: %v\n", command, err)
		}
	}
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token (optional, can be set in config.json)")
	flag.Parse()
}

func loadConfig() (Config, error) {
	var config Config

	file, err := os.Open("config.json")
	if err != nil {
		// If config file doesn't exist or can't be opened, only use command line token
		if Token == "" {
			return config, fmt.Errorf("no token provided. Either create config.json or use -t flag")
		}
		config.DiscordToken = Token
		return config, nil
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, fmt.Errorf("error decoding config file: %v", err)
	}

	// Only use token from config if not provided via command line
	if Token != "" {
		config.DiscordToken = Token
	}

	if config.DiscordToken == "" {
		return config, fmt.Errorf("no token provided. Set it in config.json or use -t flag")
	}

	return config, nil
}

// handleInitDB initializes the database with quotes from backup file (admin only, one-time use)
func (b *Bot) handleInitDB(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !b.isAdmin(m.Author.ID) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Only the admin can initialize the database")
		return err
	}

	// Check if database already has quotes
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM quotes").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Database already contains quotes. Use !restore if you want to reload from backup.")
		return err
	}

	// Load quotes from backup file
	if err := b.restoreFromFile("quotes_backup.json"); err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to initialize database: %v", err))
		return err
	}

	// Verify quotes were loaded
	err = db.QueryRow("SELECT COUNT(*) FROM quotes").Scan(&count)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Database initialized successfully with %d quotes", count))
	return err
}

func (b *Bot) restoreFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var backup BackupQuotes
	if err := json.NewDecoder(file).Decode(&backup); err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Add users and quotes
	for userID, quotes := range backup.Quotes {
		// Add user if doesn't exist
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO authorized_users (id, name, added_by, added_at) VALUES (?, ?, ?, ?)",
			userID, userID, b.config.AdminUserID, time.Now(),
		)
		if err != nil {
			return err
		}

		// Add quotes
		for _, quote := range quotes {
			_, err = tx.Exec(
				"INSERT INTO quotes (user_id, text, added_by, added_at) VALUES (?, ?, ?, ?)",
				userID, quote, b.config.AdminUserID, time.Now(),
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// handleBackup creates a backup of all quotes
func (b *Bot) handleBackup(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !b.isAdmin(m.Author.ID) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Only the admin can create backups")
		return err
	}

	// Get all quotes grouped by user
	rows, err := db.Query("SELECT user_id, text FROM quotes ORDER BY user_id, added_at")
	if err != nil {
		return err
	}
	defer rows.Close()

	backup := BackupQuotes{
		Quotes: make(map[string][]string),
	}

	for rows.Next() {
		var userID, text string
		if err := rows.Scan(&userID, &text); err != nil {
			return err
		}
		backup.Quotes[userID] = append(backup.Quotes[userID], text)
	}

	// Save to file
	file, err := os.Create("quotes_backup.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(backup); err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Backup created successfully")
	return err
}

// handleRestore restores quotes from backup
func (b *Bot) handleRestore(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !b.isAdmin(m.Author.ID) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Only the admin can restore backups")
		return err
	}

	if err := b.restoreFromFile("quotes_backup.json"); err != nil {
		return fmt.Errorf("error restoring backup: %v", err)
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "Backup restored successfully")
	return err
}

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

>>>>>>> Stashed changes
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Create a new bot instance with the loaded config
	bot, err := NewBot(config)
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(bot.messageCreate)

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
<<<<<<< Updated upstream

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Got msg: ", m.Content)
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
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

	if m.Content == "!commands" {
		// send a random response
		s.ChannelMessageSend(m.ChannelID, "Available commands: !boris")
	}

}
=======
>>>>>>> Stashed changes
