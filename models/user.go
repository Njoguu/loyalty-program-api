package models

// User data structure
type User struct {
    UID             uint       `json:"id" gorm:"primary_key"`
	Username		string		`json:"username"`
    FirstName      string    `json:"first_name"`
    LastName       string    `json:"last_name"`
	Gender			string		`json:"gender"`	
    EmailAddress   string    `json:"email_address"`
    Password       string    `json:"password"`
    City        string    `json:"city"`
    PhoneNumber   string    `json:"phone_number"`
}	

type CreateUserInput struct {
    Username string `json:"username" binding:"required"`
    FirstName string `json:"first_name" binding:"required"`
    LastName string `json:"last_name" binding:"required"`
    Gender string `json:"gender" binding:"required"`
    EmailAddress string `json:"email_address" binding:"required"`
    Password string `json:"password" binding:"required"`
    PhoneNumber string `json:"phone_number" binding:"required"`    
}