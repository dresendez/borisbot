# BorisBot - Discord Quote Manager

A Discord bot that manages and recalls quotes from your friends, with context and usage tracking. Perfect for preserving memorable moments and inside jokes from your Discord community.

## Features

- Store and recall quotes for authorized users
- Track quote usage and context
- Backup and restore functionality
- Admin controls for user management
- SQLite database for reliable storage
- Easy to set up and customize for your server

## Setup

### Prerequisites

1. Go 1.16 or higher
2. A Discord Bot Token (see below for creation)
3. SQLite3

### Creating a Discord Bot

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Go to the "Bot" section and click "Add Bot"
4. Copy the bot token (you'll need this later)
5. Go to "OAuth2" -> "URL Generator"
6. Select the following permissions:
   - Scopes: `bot`
   - Bot Permissions: `Send Messages`, `Read Message History`
7. Use the generated URL to invite the bot to your server

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/borisbot.git
cd borisbot
```

2. Install dependencies:
```bash
go mod init borisbot
go get github.com/bwmarrin/discordgo
go get github.com/mattn/go-sqlite3
```

3. Configure the bot:
   - Copy `config.json.example` to `config.json`
   - Edit `config.json` with your bot token and admin Discord user ID:
   ```json
   {
       "discord_token": "YOUR_BOT_TOKEN_HERE",
       "admin_user_id": "YOUR_DISCORD_USER_ID"
   }
   ```
   Note: To get your Discord user ID, enable Developer Mode in Discord (Settings -> App Settings -> Advanced -> Developer Mode), then right-click your name and select "Copy ID"

4. Initialize the database (as admin):
   - Start the bot: `go run bin/main.go`
   - In Discord, use the `!initdb` command to initialize the database with default quotes

### Security Notes

- Never commit `config.json` to Git - it contains sensitive information
- Keep regular backups of `quotes.db` and `quotes_backup.json`
- The following files are ignored by Git for security:
  - `config.json` (contains sensitive tokens)
  - `*.db` and `*.db-journal` (database files)
  - `quotes_backup.json` (backup of quotes)

## Commands

### User Commands
- `!quote <user>` - Get a random quote from a user
- `!listquotes <user>` - List all quotes from a user
- `!addquote <user> <quote>` - Add a new quote for an authorized user
- `!context <quote_id> <text>` - Add context to a specific quote
- `!commands` - List all available commands

### Admin Commands
- `!adduser <user_id> <name>` - Add a new authorized user
- `!listusers` - List all authorized users
- `!delquote <quote_id>` - Delete a specific quote
- `!backup` - Create a backup of all quotes
- `!restore` - Restore quotes from backup
- `!initdb` - Initialize database with default quotes (one-time use)

## File Structure

- `bin/main.go` - Main bot code
- `config.json` - Bot configuration (gitignored)
- `config.json.example` - Example configuration template
- `quotes.db` - SQLite database (gitignored)
- `quotes_backup.json` - Backup of quotes (gitignored)

## Database Management

The bot uses SQLite for storage, with two main tables:

### authorized_users
- `id` - Discord user ID or custom identifier
- `name` - User's name
- `added_by` - Admin who added the user
- `added_at` - Timestamp
- `is_active` - Whether the user is active

### quotes
- `id` - Quote ID
- `user_id` - Reference to authorized_users.id
- `text` - The quote text
- `added_by` - Who added the quote
- `added_at` - When it was added
- `context` - Optional context
- `use_count` - Times the quote was used
- `last_used` - Last time it was used

## Backup and Restore

The bot maintains quotes in both the SQLite database and a JSON backup:

1. Regular operation uses the SQLite database
2. `quotes_backup.json` serves as a backup and initial data source
3. Use `!backup` to save current state to JSON
4. Use `!restore` to restore from backup

## Customization

You can customize this bot for your server by:
1. Adding your own quotes to `quotes_backup.json`
2. Modifying the command prefixes in the code
3. Adding new features or commands

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - feel free to use and modify as needed.
