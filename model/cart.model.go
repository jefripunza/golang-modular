package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//-> main collection
type Cart struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	UserID    string `bson:"user_id"`
	ProductID string `bson:"product_id"`

	Qty int `bson:"qty"`

	CreatedAt primitive.DateTime `bson:"created_at"`
}
