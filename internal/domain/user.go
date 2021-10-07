package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	Name         string             `json:"name" bson:"name, omitempty"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Role         string             `json:"role" bson:"role"`
	RegisteredAt time.Time          `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  time.Time          `json:"lastVisitAt" bson:"lastVisitAt"`
	Verification Verification       `json:"verification" bson:"verification"`
}

type Verification struct {
	Code     string    `json:"code" bson:"code"`
	Verified bool      `json:"verified" bson:"verified"`
	Expires  time.Time `json:"expires" bson:"expires"`
}
