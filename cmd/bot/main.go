package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/zenha/oliveiras/internal/bot"
	"github.com/zenha/oliveiras/internal/database"
	"github.com/zenha/oliveiras/internal/models"
	"github.com/zenha/oliveiras/internal/scraper"
	"github.com/zenha/oliveiras/internal/telegram"
	"github.com/zenha/oliveiras/pkg/config"
)

var lastUpdateID int

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize MongoDB client
	mongoClient, err := database.NewClient(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect()

	// Initialize services
	telegramClient := telegram.NewClient(cfg.TelegramToken)
	scraperService := scraper.NewService(cfg.PythonPath, cfg.ScraperPath)
	botHandler := bot.NewHandler(telegramClient, scraperService)

	// Setup webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(string(body) + "\n")

		var update models.TelegramUpdate
		if err := json.Unmarshal(body, &update); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if we've already processed this update
		if update.UpdateID <= lastUpdateID {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Update the last processed update ID
		lastUpdateID = update.UpdateID

		if err := botHandler.HandleMessage(update.Message.Chat.ID, update.Message.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// Start the server
	log.Printf("Starting server on port %s...\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
