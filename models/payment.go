package models

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	Cod     bool `json:"cod"     bson:"cod"`
}
