package operations

import (
	"Autriche/database"
	"Autriche/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// Trouver les parents d'une personne
func GetParents(childID string) ([]models.Individual, error) {
	relCollection := database.RelationsCollection()
	indCollection := database.IndividualsCollection()
	var parents []models.Individual

	// Chercher les relations où cette personne est l'enfant
	cursor, err := relCollection.Find(context.Background(), bson.M{"child_id": childID})
	if err != nil {
		return nil, err
	}

	var relations []models.Relation
	err = cursor.All(context.Background(), &relations) // on cherche et onmet tout dans relations
	if err != nil {
		return nil, err
	}

	// Pour chaque relation, trouver le parent
	for _, rel := range relations { // pour chaque relation
		var parent models.Individual
		err := indCollection.FindOne(context.Background(), bson.M{"_id": rel.ParentID}).Decode(&parent)
		// cherche la personne avec cet ID celui quon avait dans relations donc chaque une relation
		if err == nil {
			parents = append(parents, parent) //ajoute le parent à la liste
		}
	}

	return parents, nil
}

// Trouver les enfants d'une personne
func GetChildren(parentID string) ([]models.Individual, error) {
	relCollection := database.RelationsCollection()
	indCollection := database.IndividualsCollection()
	var children []models.Individual

	// Chercher les relations où cette personne est le parent
	cursor, err := relCollection.Find(context.Background(), bson.M{"parent_id": parentID})
	if err != nil {
		return nil, err
	}

	var relations []models.Relation
	err = cursor.All(context.Background(), &relations)
	if err != nil {
		return nil, err
	}

	// Pour chaque relation, trouver l'enfant
	for _, rel := range relations {
		var child models.Individual
		err := indCollection.FindOne(context.Background(), bson.M{"_id": rel.ChildID}) //Decode traduit le BSON en une variable Go qu'on peut utiliser.
		//met le résultat dans la variable child
		//Decode a besoin de l'adresse pour savoir où mettre le résultat.
		if err == nil {
			children = append(children, child)
		}
	}

	return children, nil
}
