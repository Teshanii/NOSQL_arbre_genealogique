package operations

import (
	"Autriche/database"
	"Autriche/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// Compter le nombre total de personnes
func CountIndividuals() (int64, error) {
	collection := database.IndividualsCollection()                          // la fonction pour acceder a la collection individu..
	count, err := collection.CountDocuments(context.Background(), bson.M{}) // pour compter les documents ca veut dire les individus
	return count, err                                                       // pas de filtre on recup tout le monde
}

// Compter le nombre de personnes par genre
func CountByGender(gender string) (int64, error) {
	collection := database.IndividualsCollection()
	count, err := collection.CountDocuments(context.Background(), bson.M{"gender": gender})
	return count, err
}

// Trouver les personnes sans date de naissance
func GetIndividualsWithoutBirthDate() ([]models.Individual, error) { // on affiche en liste
	collection := database.IndividualsCollection()
	var results []models.Individual

	cursor, err := collection.Find(context.Background(), bson.M{"birth_date": ""})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &results) // recupere tous les resultats
	return results, err
}

// Calculer l'âge moyen
func GetAverageAge() (float64, error) {
	collection := database.IndividualsCollection()
	var results []models.Individual

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return 0, err
	}

	err = cursor.All(context.Background(), &results)
	if err != nil {
		return 0, err
	}

	if len(results) == 0 { //donc le nbr de personnes
		return 0, nil
	}

	total := 0 // on initialise a 0 total c pour le total dage et count c un compteur pour les personnes
	count := 0
	for _, person := range results { // ignorer les inex (_)
		if person.BirthDate != "" {
			age := 2025 - getYear(person.BirthDate)
			total = total + age
			count = count + 1
		}
	}

	if count == 0 {
		return 0, nil
	}

	return float64(total) / float64(count), nil // c pour pouvoir avoir les chiffres en virgule
}

// Extraire l'année d'une date "1985-03-15" → 1985
func getYear(date string) int {
	if len(date) < 4 { // psk si c moins de 4 caractere c que c bzr ya pas dannée
		return 0
	}
	year := 0                // on initialise year a 0 pour la boucle
	for i := 0; i < 4; i++ { // ca commmence par lindex 0
		year = year*10 + int(date[i]-'0') // date [i] donc c le nombre a chaque case de la liste
	}
	return year
}

// Trouver les dates impossibles (ici on verifie que ya pas qlq qui est né dans le futur genre incoherenece )
func GetInconsistencies() ([]string, error) {
	collection := database.IndividualsCollection()
	var results []models.Individual
	var problems []string

	cursor, err := collection.Find(context.Background(), bson.M{}) // on recuperer tous les individus
	if err != nil {
		return nil, err
	}
	cursor.All(context.Background(), &results) // on les met tous dans resultats

	for _, person := range results { // pour chaque personne
		//on ignore les postions
		if person.BirthDate != "" && getYear(person.BirthDate) > 2025 { // si ce nest pas vide
			problems = append(problems, person.FirstName+" "+person.LastName+" : né dans le futur") // on rajoute dans problemes person.FirstName+" "+person.LastName
		}
	}

	return problems, nil
}
