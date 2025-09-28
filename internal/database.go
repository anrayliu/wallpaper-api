package internal

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	uri := os.Getenv("DB_URI")

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func QueryName(client *mongo.Client, name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M

	err := client.Database("wallpapers").Collection("names").FindOne(ctx, bson.M{"access_name": name}).Decode(&result)
	if err != nil {
		return "", err
	}

	return result["stored_name"].(string), nil
}
