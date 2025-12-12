package main

import (
	"log"

	//"NOSQL_arbre_genealogique/database"
	"Autriche/database"
)

func main() {
	database.Connect()
	log.Println("🚀 Application démarrée")
}
