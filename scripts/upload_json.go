package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	// "github.com/shivam-cse/contextual-news-api/internal/models"
	"github.com/shivam-cse/contextual-news-api/pkg/startup"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const FILE = "../data/news_data.json"

func helper(database *mongo.Database) error {
	// Initialize the MongoDB with the JSON data
	collection := database.Collection("news")

	log.Println("Clearing existing data in the news collection")
	// Clear existing data
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Warning: Could not clear existing data: %v", err)
	}

	dataByte, err := os.ReadFile(FILE)
	if err != nil {
		return err
	}

	var newsData []interface{}
	err = json.Unmarshal(dataByte, &newsData)
	if err != nil {
		return err
	}

	for i, article := range newsData {
		// Add location field as GeoJSON and map id to _id
		if articleMap, ok := article.(map[string]interface{}); ok {
			// Map the id field to _id for MongoDB
			if id, ok := articleMap["id"]; ok {
				articleMap["_id"] = id
				delete(articleMap, "id") // Remove the original id field
			}
			
			// Add location field as GeoJSON
			if lat, ok := articleMap["latitude"].(float64); ok {
				if long, ok := articleMap["longitude"].(float64); ok {
					articleMap["location"] = map[string]interface{}{
						"type":        "Point",
						"coordinates": []float64{long, lat},
					}
				}
			}
		}
		newsData[i] = article // Update the article in the slice
	}

	log.Println("Inserting news articles into MongoDB")
	// insert many
	_, err = collection.InsertMany(context.Background(), newsData)
	if err != nil {
		return err
	}
	return nil
}

func UploadJSON() {
	// Load configuration
	config, err := startup.LoadConfig("../.env")
	if err != nil {
		panic("Error loading configuration: " + err.Error())
	}
	log.Println("Configuration loaded successfully")

	// Connect to MongoDB
	mongoClient, err := startup.ConnectMongoDB(config.MongoConnectionString)
	if err != nil {
		panic("Error connecting to MongoDB: " + err.Error())
	}
	defer startup.Close(mongoClient)
	log.Println("Connected to MongoDB successfully")

	// Get a handle to the specific database we want to use.
	database := mongoClient.Database(config.MongoDatabase)

	log.Println("Initializing MongoDB with JSON data")
	// Initialize the MongoDB with the JSON data
	err = helper(database)
	if err != nil {
		panic("Failed to upload news json data: " + err.Error())
	}
	println("News JSON data uploaded successfully!")
}

func main() {
	UploadJSON()
}
