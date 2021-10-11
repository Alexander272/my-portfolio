package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	UserUrl      string             `json:"userUrl" bson:"userUrl"`
	Name         string             `json:"name" bson:"name, omitempty"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Role         string             `json:"role" bson:"role"`
	AvatarUrl    string             `json:"avatarUrl" bson:"avatarUrl"`
	RegisteredAt time.Time          `json:"-" bson:"registeredAt"`
	LastVisitAt  time.Time          `json:"-" bson:"lastVisitAt"`
	Verification Verification       `json:"-" bson:"verification"`
}

type Verification struct {
	Code     string    `json:"code" bson:"code"`
	Verified bool      `json:"verified" bson:"verified"`
	Expires  time.Time `json:"expires" bson:"expires"`
}

type UserUpdate struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserUrl   string `json:"userUrl"`
	Role      string `json:"role"`
	AvatarUrl string `json:"avatarUrl"`
}
