package main

import (
	"Autriche/database"
	"Autriche/operations"
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(" .env non trouvé, on continue avec les variables d'env existantes")
	}

	// 1) Connexion à MongoDB
	database.Connect()
	log.Println(" Application démarrée")

	// 2) Contexte avec timeout pour les opérations Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3) Récupérer la collection "individuals"
	coll := database.IndividualsCollection()

	// 4) Charger les individus depuis le fichier JSON
	individuals, err := operations.LoadIndividualsFromFile("data/individuals.json")
	if err != nil {
		log.Fatal("Erreur lecture JSON:", err)
	}

	// 5) Préparer les documents pour InsertMany
	docs := make([]interface{}, len(individuals))
	for i, ind := range individuals {
		docs[i] = ind
	}

	// 6) Insérer en base
	_, err = coll.InsertMany(ctx, docs)
	if err != nil {
		log.Println(" Données déjà importées ou erreur:", err)
	} else {
		log.Println(" Import terminé !")
	}

	// 7) Afficher le menu
	ShowMenu()
}
