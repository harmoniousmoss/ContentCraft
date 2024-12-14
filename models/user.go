package models // models: package for all models

import (
	"time" // time: package for time related operations

	"go.mongodb.org/mongo-driver/bson/primitive" // primitive: package for ObjectID from MongoDB
)

type Role string // Role: type for user role

const ( // const: keyword to declare a constant
	Administrator Role = "administrator" // Administrator: constant for "administrator"
	Viewer        Role = "viewer"        // Viewer: constant for "viewer"
)

type Status string // Status: type for user status

const ( // const: keyword to declare a constant
	Approved Status = "approved" // Approved: constant for "approved"
	Pending  Status = "pending"  // Pending: constant for "pending"
)

type User struct { // User: struct for the user
	ID                 primitive.ObjectID `bson:"_id"`
	FullName           string             `json:"full_name" bson:"full_name" validate:"required"`
	Email              string             `json:"email" bson:"email" validate:"required,email"`
	Password           string             `json:"password,omitempty" bson:"password"`
	Role               Role               `json:"role" bson:"role" validate:"omitempty,oneof=administrator viewer po it comms hr cmas"` // omitempty allows default setting for the field
	Status             Status             `json:"status" bson:"status" validate:"omitempty,oneof=approved pending"`                     // omitempty allows for default setting for the field
	TermsAndConditions bool               `json:"terms_and_conditions" bson:"terms_and_conditions" validate:"required"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}
