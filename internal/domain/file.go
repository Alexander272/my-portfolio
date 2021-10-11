package domain

type File struct {
	FileType string `json:"type" bson:"type, omitempty"`
	Name     string `json:"name" bson:"name"`
	OrigName string `json:"origName" bson:"origName"`
	Url      string `json:"url" bson:"url"`
}
