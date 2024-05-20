package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type IProjectData struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}
