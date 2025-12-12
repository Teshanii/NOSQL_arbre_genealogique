package models

type Event struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Type        string `json:"type" bson:"type"`
	Date        string `json:"date" bson:"date"`
	Description string `json:"description" bson:"description"`
}
