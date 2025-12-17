package main

import (
	"Autriche/database"
	"Autriche/models"
	"Autriche/operations"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func ShowMenu() {
	var choix int

	for {
		fmt.Println("\n========== ARBRE GÉNÉALOGIQUE ==========")
		fmt.Println("1. Afficher les statistiques")
		fmt.Println("2. Chercher une personne par nom")
		fmt.Println("3. Chercher par période de naissance")
		fmt.Println("4. Voir les parents d'une personne")
		fmt.Println("5. Voir les enfants d'une personne")
		fmt.Println("6. Afficher l'arbre complet")
		fmt.Println("7. Afficher tous les individus")
		fmt.Println("8. Ajouter une personne")
		fmt.Println("9. Modifier une personne")
		fmt.Println("10. Supprimer une personne")
		fmt.Println("11. Exporter en JSON")
		fmt.Println("12. Détecter les incohérences")
		fmt.Println("0. Quitter")
		fmt.Println("=========================================")
		fmt.Print("Votre choix : ")
		fmt.Scan(&choix)

		switch choix {
		case 1:
			afficherStats()
		case 2:
			chercherParNom()
		case 3:
			chercherParDate()
		case 4:
			voirParents()
		case 5:
			voirEnfants()
		case 6:
			afficherArbreComplet()
		case 7:
			afficherTousIndividus()
		case 8:
			ajouterPersonne()
		case 9:
			modifierPersonne()
		case 10:
			supprimerPersonne()
		case 11:
			exporterJSON()
		case 12:
			detecterIncoherences()
		case 0:
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Choix invalide")
		}
	}
}

func afficherStats() {
	fmt.Println("\n--- STATISTIQUES ---")

	total, _ := operations.CountIndividuals()
	fmt.Println("Nombre total d'individus :", total)

	hommes, _ := operations.CountByGender("male")
	femmes, _ := operations.CountByGender("female")
	fmt.Println("Hommes :", hommes)
	fmt.Println("Femmes :", femmes)

	ageMoyen, _ := operations.GetAverageAge()
	fmt.Printf("Âge moyen : %.1f ans\n", ageMoyen)

	sansDate, _ := operations.GetIndividualsWithoutBirthDate()
	fmt.Println("Personnes sans date de naissance :", len(sansDate))
}

func chercherParNom() {
	var nom string
	fmt.Print("Entrez le nom à chercher : ")
	fmt.Scan(&nom)

	results, err := operations.SearchByName(nom)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("\n--- RÉSULTATS ---")
	if len(results) == 0 {
		fmt.Println("Aucune personne trouvée")
		return
	}

	for _, p := range results {
		fmt.Printf("- %s %s (né: %s) ID: %s\n", p.FirstName, p.LastName, p.BirthDate, p.ID)
	}
}

func chercherParDate() {
	var debut, fin string
	fmt.Print("Date de début (YYYY-MM-DD) : ")
	fmt.Scan(&debut)
	fmt.Print("Date de fin (YYYY-MM-DD) : ")
	fmt.Scan(&fin)

	results, err := operations.SearchByDateRange(debut, fin)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("\n--- RÉSULTATS ---")
	if len(results) == 0 {
		fmt.Println("Aucune personne trouvée")
		return
	}

	for _, p := range results {
		fmt.Printf("- %s %s (né: %s)\n", p.FirstName, p.LastName, p.BirthDate)
	}
}

func voirParents() {
	var id string
	fmt.Print("Entrez l'ID de la personne : ")
	fmt.Scan(&id)

	parents, err := operations.GetParents(id)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("\n--- PARENTS ---")
	if len(parents) == 0 {
		fmt.Println("Aucun parent trouvé")
		return
	}

	for _, p := range parents {
		fmt.Printf("- %s %s (né: %s)\n", p.FirstName, p.LastName, p.BirthDate)
	}
}

func voirEnfants() {
	var id string
	fmt.Print("Entrez l'ID de la personne : ")
	fmt.Scan(&id)

	enfants, err := operations.GetChildren(id)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("\n--- ENFANTS ---")
	if len(enfants) == 0 {
		fmt.Println("Aucun enfant trouvé")
		return
	}

	for _, e := range enfants {
		fmt.Printf("- %s %s (né: %s)\n", e.FirstName, e.LastName, e.BirthDate)
	}
}

func afficherArbreComplet() {
	ctx := context.Background()

	fmt.Println("\n--- ARBRE GÉNÉALOGIQUE COMPLET ---")

	coll := database.IndividualsCollection()
	var root models.Individual
	if err := coll.FindOne(ctx, bson.M{}).Decode(&root); err != nil {
		fmt.Println("Impossible de trouver un ancêtre:", err)
		return
	}

	visited := make(map[string]bool)
	tree, err := operations.BuildTree(ctx, root.ID, visited)
	if err != nil {
		fmt.Println("Erreur BuildTree:", err)
		return
	}

	operations.PrintTree(tree, 0)
}

// func afficherDescendants(parentID string, niveau int) {
// 	enfants, err := operations.GetChildren(parentID)
// 	if err != nil || len(enfants) == 0 {
// 		return
// 	}

// 	indent := ""
// 	for i := 0; i < niveau; i++ {
// 		indent += "  "
// 	}

// 	for _, e := range enfants {
// 		fmt.Printf("%s +- %s %s (né: %s)\n", indent, e.FirstName, e.LastName, e.BirthDate)
// 		afficherDescendants(e.ID, niveau+1)
// 	}
// }

func afficherTousIndividus() {
	fmt.Println("\n--- TOUS LES INDIVIDUS ---")

	results, err := operations.SearchByName("")
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	for i, p := range results {
		fmt.Printf("%d. %s %s (né: %s) - %s\n", i+1, p.FirstName, p.LastName, p.BirthDate, p.Gender)
	}

	fmt.Printf("\nTotal : %d individus\n", len(results))
}

func ajouterPersonne() {
	var id, prenom, nom, dateNaissance, genre string

	fmt.Print("ID : ")
	fmt.Scan(&id)
	fmt.Print("Prénom : ")
	fmt.Scan(&prenom)
	fmt.Print("Nom : ")
	fmt.Scan(&nom)
	fmt.Print("Date de naissance (YYYY-MM-DD) : ")
	fmt.Scan(&dateNaissance)
	fmt.Print("Genre (male/female) : ")
	fmt.Scan(&genre)

	personne := models.Individual{
		ID:        id,
		FirstName: prenom,
		LastName:  nom,
		BirthDate: dateNaissance,
		Gender:    genre,
	}

	err := operations.InsertIndividual(personne)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("Personne ajoutée !")
}

func modifierPersonne() {
	var id string
	fmt.Print("Entrez l'ID de la personne à modifier : ")
	fmt.Scan(&id)

	// Vérifier que la personne existe
	personne, err := operations.FindIndividualByID(id)
	if err != nil {
		fmt.Println("Personne non trouvée")
		return
	}

	fmt.Printf("Personne trouvée : %s %s\n", personne.FirstName, personne.LastName)
	fmt.Println("\nQue voulez-vous modifier ?")
	fmt.Println("1. Prénom")
	fmt.Println("2. Nom")
	fmt.Println("3. Date de naissance")
	fmt.Println("4. Genre")

	var choix int
	fmt.Print("Votre choix : ")
	fmt.Scan(&choix)

	var nouvelleValeur string
	var champ string

	switch choix {
	case 1:
		fmt.Print("Nouveau prénom : ")
		fmt.Scan(&nouvelleValeur)
		champ = "first_name"
	case 2:
		fmt.Print("Nouveau nom : ")
		fmt.Scan(&nouvelleValeur)
		champ = "last_name"
	case 3:
		fmt.Print("Nouvelle date (YYYY-MM-DD) : ")
		fmt.Scan(&nouvelleValeur)
		champ = "birth_date"
	case 4:
		fmt.Print("Nouveau genre (male/female) : ")
		fmt.Scan(&nouvelleValeur)
		champ = "gender"
	default:
		fmt.Println("Choix invalide")
		return
	}

	err = operations.UpdateIndividual(id, bson.M{champ: nouvelleValeur})
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("Personne modifiée !")
}

func supprimerPersonne() {
	var id string
	fmt.Print("Entrez l'ID de la personne à supprimer : ")
	fmt.Scan(&id)

	err := operations.DeleteIndividual(id)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("Personne supprimée !")
}

func exporterJSON() {
	err := operations.ExportToJSON("export_individuals.json")
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("Export terminé ! Fichier : export_individuals.json")
}

func detecterIncoherences() {
	problems, err := operations.GetInconsistencies()
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("\n--- INCOHÉRENCES ---")
	if len(problems) == 0 {
		fmt.Println("Aucune incohérence trouvée")
		return
	}

	for _, p := range problems {
		fmt.Println("-", p)
	}
}
