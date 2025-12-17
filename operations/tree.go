package operations

import (
	"Autriche/database"
	"Autriche/models"
	"context"
	"log"
	"math/rand"

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
	defer cursor.Close(context.Background())

	var relations []models.Relation
	if err = cursor.All(context.Background(), &relations); err != nil {
		return nil, err
	}

	// Pour chaque relation, trouver le parent
	for _, rel := range relations {
		var parent models.Individual
		err := indCollection.FindOne(context.Background(), bson.M{"_id": rel.ParentID}).Decode(&parent)
		if err == nil {
			parents = append(parents, parent)
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
	defer cursor.Close(context.Background())

	var relations []models.Relation
	if err = cursor.All(context.Background(), &relations); err != nil {
		return nil, err
	}

	// Pour chaque relation, trouver l'enfant
	for _, rel := range relations {
		var child models.Individual
		err := indCollection.FindOne(context.Background(), bson.M{"_id": rel.ChildID}).Decode(&child)
		if err == nil {
			children = append(children, child)
		}
	}

	return children, nil
}

// Générer des relations parent/enfant avec un ancêtre commun
func GenerateRandomRelations(ctx context.Context) error {
	collInd := database.IndividualsCollection()
	collRel := database.RelationsCollection()

	// Vider les relations avant de régénérer
	if _, err := collRel.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}

	// Récupérer tous les individus
	cur, err := collInd.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	var individuals []models.Individual
	if err := cur.All(ctx, &individuals); err != nil {
		return err
	}

	if len(individuals) == 0 {
		return nil
	}

	// 1) Ancêtre commun = premier individu du JSON
	root := individuals[0]
	log.Printf("Ancêtre commun choisi : %s %s (%s)", root.FirstName, root.LastName, root.ID)

	var rels []interface{}

	// 2) Lier la racine à tous les autres individus
	for _, ind := range individuals {
		if ind.ID == root.ID {
			continue
		}
		rels = append(rels, models.Relation{
			ParentID: root.ID,
			ChildID:  ind.ID,
			Relation: "parent",
		})
	}

	// 3) Relations supplémentaires aléatoires pour enrichir l'arbre
	// (optionnel, mais attention aux cycles -> protégé par visited dans BuildTree)
	for i := 0; i < len(individuals)/2; i++ {
		parent := individuals[rand.Intn(len(individuals))]
		child := individuals[rand.Intn(len(individuals))]

		// éviter boucle triviale et éviter de mettre quelqu'un au-dessus de la racine
		if parent.ID == child.ID || child.ID == root.ID {
			continue
		}

		rels = append(rels, models.Relation{
			ParentID: parent.ID,
			ChildID:  child.ID,
			Relation: "parent",
		})
	}

	if len(rels) == 0 {
		return nil
	}

	_, err = collRel.InsertMany(ctx, rels)
	return err
}

// Construire l'arbre à partir d'un ID racine
func BuildTree(ctx context.Context, rootID string, visited map[string]bool) (*models.TreeNode, error) {
	indCollection := database.IndividualsCollection()

	// si déjà visité, on stoppe pour éviter la boucle
	if visited[rootID] {
		return nil, nil
	}
	visited[rootID] = true

	// Récupérer l'individu racine
	var root models.Individual
	if err := indCollection.FindOne(ctx, bson.M{"_id": rootID}).Decode(&root); err != nil {
		return nil, err
	}

	node := &models.TreeNode{Individual: root}

	// Récupérer ses enfants via GetChildren
	children, err := GetChildren(rootID)
	if err != nil {
		return nil, err
	}

	// Construire récursivement l'arbre pour chaque enfant
	for _, child := range children {
		childNode, err := BuildTree(ctx, child.ID, visited)
		if err != nil {
			return nil, err
		}
		if childNode != nil {
			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}

// Afficher l'arbre dans la console
func PrintTree(node *models.TreeNode, level int) {
	if node == nil {
		return
	}
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}
	log.Printf("%s +- %s %s (%s) %s", indent, node.Individual.FirstName, node.Individual.LastName, node.Individual.BirthDate, node.Individual.Gender)
	for _, child := range node.Children {
		PrintTree(child, level+1)
	}
}
