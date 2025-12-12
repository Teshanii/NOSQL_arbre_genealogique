package operations

import (
	"Autriche/database"
	"Autriche/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// Chercher des personnes par nom
func SearchByName(name string) ([]models.Individual, error) {
	collection := database.IndividualsCollection()
	var results []models.Individual
	// regex ( quand tu cherches lin ca te met lina )
	filter := bson.M{"last_name": bson.M{"$regex": name, "$options": "i"}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &results)
	return results, err
}

// Chercher des personnes nées entre deux dates
func SearchByDateRange(startDate string, endDate string) ([]models.Individual, error) {
	collection := database.IndividualsCollection()
	var results []models.Individual
	//superieur ou egal (gte) //inferieur ou egal (lte)
	filter := bson.M{"birth_date": bson.M{"$gte": startDate, "$lte": endDate}}
	cursor, err := collection.Find(context.Background(), filter)
	//nil vide (retourne rien)
	if err != nil {
		return nil, err
	}
	//on retourne le resultat en liste
	err = cursor.All(context.Background(), &results)
	return results, err
}

// Chercher des personnes par genre
func SearchByGender(gender string) ([]models.Individual, error) {
	collection := database.IndividualsCollection()
	var results []models.Individual

	filter := bson.M{"gender": gender}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &results)
	return results, err
}
