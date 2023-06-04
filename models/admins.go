package models

import (
	"github.com/jinzhu/gorm"
)

// Admin DB Table
type Admins struct{
	gorm.Model
	Username	string	
	Password	string
	IsAdmin bool `json:"is_admin" gorm:"default:true"`
}

// User Login Request data
type AdminLoginInput struct {
	Username    string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}