package models

type Relation struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	ParentID string `json:"parent_id" bson:"parent_id"`
	ChildID  string `json:"child_id" bson:"child_id"`
	Relation string `json:"relation" bson:"relation"`
}
