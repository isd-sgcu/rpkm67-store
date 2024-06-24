package model

import (
	"github.com/google/uuid"
)

type Object struct {
	Base
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	ImageUrl  string    `json:"image_url" gorm:"mediumtext"`
	ObjectKey string    `json:"object_key" gorm:"mediumtext"`
}
