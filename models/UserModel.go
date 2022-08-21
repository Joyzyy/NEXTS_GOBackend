package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username   string             `json:"username,omitempty" bson:"username,omitempty" validate:"required"`
	Password   string             `json:"password,omitempty" bson:"password,omitempty" validate:"required,min=6"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty" validate:"required"`
	Created_at time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated_at time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
