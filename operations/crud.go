package operations

import (
	"Autriche/database"
	"Autriche/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// Ajouter une personne
func InsertIndividual(ind models.Individual) error {
	collection := database.IndividualsCollection()
	_, err := collection.InsertOne(context.Background(), ind)
	return err
}

// trouver une personne a partir de son ID
func FindIndividualByID(id string) (models.Individual, error) {
	collection := database.IndividualsCollection()
	var individual models.Individual
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&individual)
	return individual, err
}

// modifier une personne
func UpdateIndividual(id string, updates bson.M) error {
	collection := database.IndividualsCollection()
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

// supprimer une personne
func DeleteIndividual(id string) error {
	collection := database.IndividualsCollection()
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
