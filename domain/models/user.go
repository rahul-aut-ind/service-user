package models

import (
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
)
