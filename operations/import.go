package operations

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"Autriche/models"
)

func LoadIndividualsFromFile(path string) ([]models.Individual, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var individuals []models.Individual
	if err := json.Unmarshal(bytes, &individuals); err != nil {
		return nil, err
	}

	return individuals, nil
}
