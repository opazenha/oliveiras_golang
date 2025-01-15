package database

import (
	"context"
	"encoding/json"
	"log"

	"github.com/zenha/oliveiras/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client represents a MongoDB client
type Client struct {
	client *mongo.Client
}

// NewClient creates a new MongoDB client
func NewClient(mongoURI string) (*Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

// Disconnect closes the MongoDB connection
func (c *Client) Disconnect() error {
	return c.client.Disconnect(context.TODO())
}

// GetAirbnbByDate retrieves Airbnb listings for a date range
func (c *Client) GetAirbnbByDate(startDate, endDate string) ([]models.AirbnbData, error) {
	collection := c.client.Database("oliveiras").Collection("airbnb")

	filter := bson.M{
		"start_date": bson.M{"$gte": startDate},
		"end_date":   bson.M{"$lte": endDate},
	}

	log.Println("Performing query with filter:", filter)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []models.AirbnbData
	for cursor.Next(context.TODO()) {
		var data models.AirbnbData
		if err := cursor.Decode(&data); err != nil {
			log.Println("Error decoding document:", err)
			return nil, err
		}
		results = append(results, data)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	log.Printf("Query completed. Found %d results.\n", len(results))
	return results, nil
}

// GetBookingByDate retrieves booking listings for a date range
func (c *Client) GetBookingByDate(startDate, endDate string) ([]models.BookingData, error) {
	collection := c.client.Database("oliveiras").Collection("booking")

	filter := bson.M{
		"start_date": bson.M{"$gte": startDate},
		"end_date":   bson.M{"$lte": endDate},
	}

	log.Println("Performing query with filter:", filter)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []models.BookingData
	for cursor.Next(context.TODO()) {
		var rawDoc bson.M
		if err := cursor.Decode(&rawDoc); err != nil {
			log.Println("Error decoding raw document:", err)
			return nil, err
		}

		// Extract ObjectID
		if id, ok := rawDoc["_id"].(primitive.ObjectID); ok {
			rawDoc["_id"] = bson.M{"$oid": id.Hex()}
		}

		// Convert the document back to JSON
		jsonData, err := json.Marshal(rawDoc)
		if err != nil {
			log.Printf("Error marshaling document to JSON: %v", err)
			continue
		}

		// Decode into BookingData
		var data models.BookingData
		if err := json.Unmarshal(jsonData, &data); err != nil {
			log.Printf("Error unmarshaling into BookingData: %v", err)
			continue
		}

		results = append(results, data)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	log.Printf("Query completed. Found %d results.\n", len(results))
	return results, nil
}

// Helper function to safely get string values from map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
