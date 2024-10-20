package models

import "gorm.io/gorm"

type (
	// User represents a user in the system
	User struct {
		gorm.Model
		ID    uint   `gorm:"primaryKey"`
		Name  string `json:"name"`
		Email string `json:"email" gorm:"unique"`
	}

	Response struct {
		Data interface{} `json:"data"`
	}
)
