package startup

import (
	"context"
	"time"

	"github.com/shivam-cse/contextual-news-api/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetContext(dur int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		context.Background(),
		time.Duration(dur)*time.Second,
	)
}

func ConnectMongoDB(uri string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(uri)
    connectCtx, connectCancel := GetContext(10)
    defer connectCancel()

	client, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify that the connection is established and the server is available.
	pingCtx, pingCancel := GetContext(5)
	defer pingCancel()
	err = client.Ping(pingCtx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Close(client *mongo.Client) {
	if client == nil {
		return
	}

	ctx, cancel := GetContext(5)
	defer cancel()
	client.Disconnect(ctx)
}

func IsConnected(client *mongo.Client) bool {
	if client == nil {
		return false
	}

	ctx, cancel := GetContext(5)
	defer cancel()

	err := client.Ping(ctx, nil)
	return err == nil
}

func CreateIndexOnNewsColl(db *mongo.Database) error {
	// Create a 2dsphere index on the location field
	// and a text index on the title and description fields
	// which allows for efficient querying of news articles based on their location and user query.
	collection := db.Collection(constants.NEWS)

	indexModel := []mongo.IndexModel{
		{
			Keys: bson.D{
				primitive.E{Key: "location", Value: "2dsphere"},
			},
		},
		{
			Keys: bson.D{
				primitive.E{Key: "title", Value: "text"},
				primitive.E{Key: "description", Value: "text"},
			},
			Options: options.Index().SetWeights(
				bson.M{
					"title": 5,
					"description": 3,
				},
			),
		},
    }

	_, err := collection.Indexes().CreateMany(context.Background(), indexModel)
	return err
}