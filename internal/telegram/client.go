package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Client represents a Telegram bot client
type Client struct {
	token  string
	client *http.Client
}

// NewClient creates a new Telegram client
func NewClient(token string) *Client {
	return &Client{
		token:  token,
		client: &http.Client{},
	}
}

// SendMessage sends a message to a Telegram chat
func (c *Client) SendMessage(chatID int, text string) error {
	responseBytes, err := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+c.token+"/sendMessage", bytes.NewBuffer(responseBytes))
	if err != nil {
		log.Println(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(respBody) + "\n")
	return nil
}
