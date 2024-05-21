package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -> main collection
type Report struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// auth, unique:username
	Username string `bson:"username"`
	Password string `bson:"password"`

	// info
	Name     string `bson:"name"`
	ImageURL string `bson:"image_url"`

	CreatedAt primitive.DateTime  `bson:"created_at"`
	UpdatedAt *primitive.DateTime `bson:"updated_at,omitempty"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty"`
}
