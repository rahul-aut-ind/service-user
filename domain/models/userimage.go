package models

import "time"

type (
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
