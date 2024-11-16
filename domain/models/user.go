package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	// User represents a user in the system
	User struct {
		gorm.Model
		ID      int64  `gorm:"primaryKey"`
		Name    string `json:"name"`
		Email   string `json:"email" gorm:"unique"`
		Address string `json:"address"`
		Age     int    `json:"age"`
	}

	Response struct {
		Data interface{} `json:"data"`
	}

	Request struct {
		FirstName string `json:"firstName" validate:"required,min=2,max=100,alpha"`
		LastName  string `json:"lastName" validate:"required,min=2,max=100,alpha"`
		Email     string `json:"email" validate:"required,email"`
		Address   string `json:"address" validate:"min=5,max=300"`
		Age       int    `json:"age" validate:"gte=18,lte=100"`
	}

	UserImage struct {
		IsDeleted bool      `json:"isDeleted" validate:"required"`
		UserID    string    `json:"userId" validate:"required"`
		ImageID   string    `json:"imageId" validate:"required"`
		Path      string    `json:"path" validate:"required"`
		TakenAt   time.Time `json:"takenAt" validate:"required"`
		UpdatedAt time.Time `json:"updatedAt" validate:"required"`
	}

	UserImageResult struct {
		UserImages []UserImage
		Page       Page
	}

	Page struct {
		LastEvaluatedKey map[string]string `json:"params"`
	}

	UploadResponse struct {
		ID string `json:"id"`
	}

	ImageResponse struct {
		ImageID string    `json:"id"`
		Path    string    `json:"path"`
		TakenAt time.Time `json:"takenAt"`
	}

	PaginatedImageResponse struct {
		Images []ImageResponse `json:"items"`
		Page   Page            `json:"nextPage"`
	}

	PaginatedInput struct {
		UserID           string
		LastImageID      string
		LastImageTakenAt string
		Limit            int32
	}

	Metadata struct {
		TakenAt time.Time `json:"takenAt" validate:"required"`
		Type    string    `json:"type" validate:"required"`
	}
)

const (
	ImageKey            = "image"
	MetadataKey         = "metadata"
	DefaultHistoryLimit = 10
)
