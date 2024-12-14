package database

import (
	"context"
	"fmt"
	"log"
	"myfiberproject/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
)

// ConnectMongoDB establishes a connection to MongoDB using the URI from the environment.
func ConnectMongoDB() error {
	// Load environment variables
	config.LoadEnv()

	// Fetch MongoDB URI from environment variables
	mongoURI := config.GetEnv("MONGO_URI", "")
	if mongoURI == "" {
		return fmt.Errorf("MONGO_URI environment variable is not set in .env or system environment")
	}

	log.Println("Attempting to connect to MongoDB...")

	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error

	// Establish MongoDB connection with retries
	MongoClient, err = connect(clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB!")

	// Create indexes for collections
	CreateIndexesForCollections()

	return nil
}

// connect attempts to connect to MongoDB with retries.
func connect(clientOptions *options.ClientOptions) (*mongo.Client, error) {
	var client *mongo.Client
	var err error

	for attempt := 1; attempt <= 5; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			err = client.Ping(ctx, nil)
			if err == nil {
				return client, nil
			}
		}

		log.Printf("Failed to connect to MongoDB: %v, retrying in %d seconds...", err, attempt*2)
		time.Sleep(time.Duration(attempt*2) * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to MongoDB after several attempts: %v", err)
}

// CreateTextIndex creates a text index on the specified fields of a collection.
func CreateTextIndex(collectionName string, fields []string) {
	collection := MongoClient.Database("gogemini").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the text index keys
	indexKeys := bson.D{}
	for _, field := range fields {
		indexKeys = append(indexKeys, bson.E{Key: field, Value: "text"})
	}

	indexModel := mongo.IndexModel{
		Keys: indexKeys,
	}

	// Set options for creating indexes
	indexOptions := options.CreateIndexes().SetMaxTime(10 * time.Second)
	createdIndexName, err := collection.Indexes().CreateOne(ctx, indexModel, indexOptions)
	if err != nil {
		log.Fatalf("Failed to create text index on %s collection: %v", collectionName, err)
	} else if createdIndexName != "" {
		log.Printf("Text index on %s collection created or verified successfully", collectionName)
	}
}

// CreateIndexesForCollections initializes indexes for all collections.
func CreateIndexesForCollections() {
	// Define the collections and their text index fields
	collections := map[string][]string{

		"activity_category": {"name"},
	}

	// Create indexes for each collection
	for collectionName, fields := range collections {
		CreateTextIndex(collectionName, fields)
	}
}

// GetMongoClient provides access to the MongoDB client instance.
func GetMongoClient() *mongo.Client {
	return MongoClient
}
