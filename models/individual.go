package models

type Individual struct {
	ID        string   `json:"id" bson:"_id,omitempty"`
	FirstName string   `json:"first_name" bson:"first_name"`
	LastName  string   `json:"last_name" bson:"last_name"`
	BirthDate string   `json:"birth_date" bson:"birth_date"`
	DeathDate string   `json:"death_date,omitempty" bson:"death_date,omitempty"`
	Gender    string   `json:"gender" bson:"gender"`
	Events    []string `json:"events,omitempty" bson:"events,omitempty"`
}
