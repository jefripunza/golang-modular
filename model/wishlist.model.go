package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//-> main collection
type Wishlist struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	UserID    string `bson:"user_id"`
	ProductID string `bson:"product_id"`

	CreatedAt primitive.DateTime `bson:"created_at"`
}
