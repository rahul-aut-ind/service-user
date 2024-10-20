package models

import "gorm.io/gorm"

// User represents a user in the system
type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}
