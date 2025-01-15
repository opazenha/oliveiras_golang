package gemini

import (
	"context"
	"fmt"
	"log"

	"github.com/zenha/oliveiras/internal/models"
	"google.golang.org/genai"
)

func NewClient(apiKey string) (*genai.Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGoogleAI,
	})
	return client, err
}

// Call the GenerateContent method
func GenerateContent(client *genai.Client, prompt string) (string, error) {
	config := genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Role: "system",
			Parts: []*genai.Part{
				{
					Text: "You are a very talented and experiences hotel manager. Your task is according to the information provided regarding the listings from the same location where your rent house is, provide an appropriate price for each of the dates provided on the listings. Respond with JUST the date: price on each line. >>> Example: 2025-01-14: 112.99 >>> Rent House Information: 3 rooms. 2 double bed rooms. 1 dual single bed room. 2 bathrooms. No pool.",
				},
			},
		},
	}

	result, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.0-flash-exp",
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{
						Text: prompt,
					},
				},
			},
		},
		&config)
	if err != nil {
		log.Println("API call failed: ", err)
		return "", err
	}

	response := string(result.Candidates[0].Content.Parts[0].Text)

	return response, nil
}

func PrepareAirbnbPrompt(airbnbListings []models.AirbnbData) string {
	var prompt string
	for _, listing := range airbnbListings {
		formattedString := fmt.Sprintf("%s to %s: %.2f > Info: %s - ", listing.StartDate, listing.EndDate, listing.Listing.Price, listing.Listing.BedConfiguration)
		prompt += formattedString
	}
	// log.Println("Airbnb Prompt: ", prompt)
	return prompt
}

func PrepareBookingPrompt(bookingListings []models.BookingData) string {
	var prompt string
	for _, listing := range bookingListings {
		formattedString := fmt.Sprintf("%s to %s: %.2f > Info: %s - ", listing.StartDate, listing.EndDate, listing.Price, listing.BedConfiguration)
		prompt += formattedString
	}
	// log.Println("Booking Prompt: ", prompt)
	return prompt
}
