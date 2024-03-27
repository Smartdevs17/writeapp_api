package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `gorm:"unique"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	Documents []Document
}

type Document struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Count   int    `json:"count"`
	UserID  uint   // Foreign key referencing the ID field of the User struct
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Define the relationship
}
