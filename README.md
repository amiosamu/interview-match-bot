# Interview Match Bot

A Telegram bot that helps developers find interview preparation partners based on their field of interest and experience level.

## Features

- No authentication required
- Match with other developers based on:
  - Field of interest (Backend, Frontend, Full Stack)
  - Experience level (Intern, Junior, Middle, Senior)
- Instant notifications when matches are found

## Setup

1. Create a new bot with [@BotFather](https://t.me/botfather) on Telegram
2. Get your bot token
3. Replace `YOUR_BOT_TOKEN` in `main.go` with your actual bot token
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Run the bot:
   ```bash
   go run main.go
   ```

## Usage

1. Start the bot with `/start`
2. Select your field of interest:
3. Select your experience level:
   - `/intern` - Intern
   - `/junior` - Junior
   - `/middle` - Middle
   - `/senior` - Senior
4. The bot will notify you when it finds someone matching your criteria

## Development

The bot is built using:
- Go
- [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)