package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type ISqlCredential struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// UserID    string             `bson:"user_id"`
	ProjectID string `bson:"project_id"` // use

	Type     string `bson:"type"`
	Host     string `bson:"host"`
	Port     any    `bson:"port"`
	User     string `bson:"user"`
	Password string `bson:"pass"`
	Name     string `bson:"name"` // use
	// Timeout  *any   `bson:"timeout"`
}

type INoSqlCredential struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	ProjectID string             `bson:"project_id"` // use

	Host     string `bson:"host"`
	Port     string `bson:"port"`
	User     string `bson:"user"`
	Password string `bson:"pass"`
	Name     string `bson:"name"` // use
	Timeout  *any   `bson:"timeout"`
}

type IEmailCredential struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	ProjectID string             `bson:"project_id"` // use

	Host     string `bson:"host"`
	Port     string `bson:"port"`
	Email    string `bson:"email"` // use
	Password string `bson:"pass"`
}
