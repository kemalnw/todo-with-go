package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo represent the todo model
type Todo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

// AddTimeStamps handle timestamp
func (t *Todo) AddTimeStamps() {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}
