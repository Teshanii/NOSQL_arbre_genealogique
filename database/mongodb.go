package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal(" MONGO_URI manquant dans le fichier .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf(" Erreur de connexion MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf(" Impossible de pinger MongoDB: %v", err)
	}

	log.Println(" Connecté à MongoDB Atlas")
	Client = client
}

func IndividualsCollection() *mongo.Collection {
	return Client.Database("genealogy").Collection("individuals")
}

func RelationsCollection() *mongo.Collection {
	return Client.Database("genealogy").Collection("relations")
}
