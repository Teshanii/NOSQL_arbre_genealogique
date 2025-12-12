package operations

import (
	"Autriche/database"
	"Autriche/models"
	"context"
	"encoding/json"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

// Cette fonction est pour exporter tous les individus dans un fichier JSON
func ExportToJSON(filename string) error {
	collection := database.IndividualsCollection()
	var individuals []models.Individual // on definie une liste vide

	cursor, err := collection.Find(context.Background(), bson.M{}) // prend tout le monde
	if err != nil {
		return err
	}

	err = cursor.All(context.Background(), &individuals) // cursor.All(...) récupère tous les résultats
	// et mets les resultats dans individuals
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)   // prépare l'écriture JSON dans le fichier
	return encoder.Encode(individuals) // ecrit la liste des personnes en json dans le fichier
}

// Importer des individus depuis un fichier JSON
func ImportFromJSON(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var individuals []models.Individual
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&individuals) //lit le JSON et met les données dans individuals
	if err != nil {
		return err
	}

	collection := database.IndividualsCollection()
	for _, ind := range individuals { //pour chaque personne dans la liste et on ignore l'index
		_, err := collection.InsertOne(context.Background(), ind)
		if err != nil {
			return err
		}
	}

	return nil
}
