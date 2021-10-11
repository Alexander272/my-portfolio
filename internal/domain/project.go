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

type Project struct {
	Id          primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId, omitempty"`
	Name        string             `json:"name" bson:"name, omitempty"`
	Description string             `json:"description" bson:"description"`
	Files       []File             `json:"files" bson:"files"`
	Access      AccessType         `json:"access" bson:"access"`
	Published   bool               `json:"published" bson:"published,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
