package database

import (
	"context"
	"log"

	"github.com/zenha/oliveiras/internal/models"
	"go.mongodb.org/mongo-driver/bson"
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
