package main

import (
	"Autriche/models"
	"Autriche/operations"
	"fmt"
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
		fmt.Println("6. Ajouter une personne")
		fmt.Println("7. Supprimer une personne")
		fmt.Println("8. Exporter en JSON")
		fmt.Println("9. Détecter les incohérences")
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
			ajouterPersonne()
		case 7:
			supprimerPersonne()
		case 8:
			exporterJSON()
		case 9:
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
