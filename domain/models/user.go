package models

import "gorm.io/gorm"

type (
	// User represents a user in the system
	User struct {
		gorm.Model
		ID    int64  `gorm:"primaryKey"`
		Name  string `json:"name"`
		Email string `json:"email" gorm:"unique"`
	}

	Response struct {
		Data interface{} `json:"data"`
	}
)

const (
	RequestAccepted   = "ok"
	ErrMsgNoUserfound = "finding user"
)
