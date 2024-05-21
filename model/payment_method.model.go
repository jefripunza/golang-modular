package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//-> main collection
type PaymentMethod struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name    string `bson:"name"`
	LogoURL string `bson:"logo_url"`
	Code    string `bson:"code"`
}
