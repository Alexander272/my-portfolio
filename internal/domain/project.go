package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessType string

const (
	All    AccessType = "all"
	Link   AccessType = "link"
	Nobody AccessType = "nobody"
)

type ProjectMin struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
type SelfProjectMin struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	Access    AccessType         `json:"access" bson:"access"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Project struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description string             `json:"description" bson:"description"`
	Files       []File             `json:"files" bson:"files"`
	Access      AccessType         `json:"-" bson:"access"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
type SelfProject struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description string             `json:"description" bson:"description"`
	Files       []File             `json:"files" bson:"files"`
	Access      AccessType         `json:"access" bson:"access"`
	Published   bool               `json:"published" bson:"published,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type ProjectInput struct {
	UserId      primitive.ObjectID `json:"userId" bson:"userId"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Files       []File             `json:"files" bson:"files"`
	Access      AccessType         `json:"access" bson:"access"`
	Published   bool               `json:"published" bson:"published"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
