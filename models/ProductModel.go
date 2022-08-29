package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	TabType     string             `json:"tabType,omitempty" bson:"tabType,omitempty" validate:"required"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty" validate:"required"`
	Category    string             `json:"category,omitempty" bson:"category,omitempty" validate:"required"`
	Sizes       []int              `json:"sizes,omitempty" bson:"sizes,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	Quantity    int                `json:"quantity,omitempty" bson:"quantity,omitempty" validate:"required"`
}
