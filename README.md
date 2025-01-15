# Oliveiras Bot

A Telegram bot that helps track and analyze vacation rental listings from Airbnb and Booking.com. The bot provides real-time price analysis and listing information for specified date ranges.

## Features

- Real-time scraping of Airbnb and Booking.com listings
- Price analysis including average, highest, and lowest prices
- Total listings count for both platforms
- MongoDB integration for data persistence
- Telegram bot interface for easy interaction

## Project Structure

```
oliveiras/
├── cmd/
│   └── bot/              # Main application entry point
├── internal/
│   ├── bot/             # Bot message handling logic
│   ├── database/        # MongoDB operations
│   ├── models/          # Data structures and types
│   ├── scraper/         # Web scraping functionality
│   └── telegram/        # Telegram API client
├── pkg/
│   └── config/          # Configuration management
```

## Prerequisites

- Go 1.23 or higher
- MongoDB
- Python environment for the scraper
- Telegram Bot Token

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
MONGO_ATLAS_URI=your_mongo_uri
TELEGRAM_BOT_TOKEN=your_bot_token
PYTHON_PATH=/path/to/python
SCRAPER_PATH=/path/to/scraper/script
SERVER_PORT=7771
```

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set up your environment variables in `.env`
4. Run the bot:
   ```bash
   go run cmd/bot/main.go
   ```

## Usage

The bot responds to the following commands:

- `/scrape [start_date] [end_date]` - Scrapes and analyzes listings for the specified date range
  Example: `/scrape 2025-01-14 2025-01-16`

## Architecture

- **Bot Handler**: Manages incoming Telegram messages and command routing
- **Scraper Service**: Interfaces with Python scraping script
- **Database Layer**: Handles MongoDB operations for data persistence
- **Telegram Client**: Manages Telegram API communication
- **Configuration**: Centralized configuration management

## Future Enhancements

- AI integration for smart price predictions
- Additional booking platforms support
- Advanced analytics and reporting
- User preferences and saved searches
- Price alert notifications

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.