package models

import (
    "html"
	"errors"
	"strings"
	"api/utils/token"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User data structure
type User struct {
    gorm.Model

    // UID uint `json:"id" gorm:"primary_key;unique"`
	Username string `gorm:"unique" json:"username"`
    FirstName string `json:"firstname"`
    LastName string `json:"lastname"`
	Gender string `json:"gender"`	
    EmailAddress string `json:"email"`
    Password string `json:"password"`
    City string `json:"city"`
    PhoneNumber string `json:"phone_number"`
}	

type CreateUserInput struct {
    Username string `json:"username" binding:"required"`
    FirstName string `json:"firstname" binding:"required"`
    LastName string `json:"lastname" binding:"required"`
    Gender string `json:"gender" binding:"required"`
    EmailAddress string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
    City string `json:"city" binding:"required"`
    PhoneNumber string `json:"phone_number" binding:"required"`    
}
type LoginInput struct{
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}


// MISC FUNCTIONS
// Save a new user record to the database
func (user *User) SaveUser() (*User, error){

	var err error = DB.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// Hash password input
func (user *User) BeforeSave() error{

	// Hash Given Password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil{
        return err
    }

    user.Password = string(hashedPassword)

    // remove spaces in username
    user.Username = html.EscapeString(strings.TrimSpace(user.Username))

    return nil
}

// Verify password hash and password given are identical
func VerifyPassword(password, hashedPassword string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Verify Login Details
func LoginCheck(username string, password string) (string, error){
	
	user := User{}

    var err error = DB.Model(User{}).Where("username = ?", username).Take(&user).Error

	if err != nil{
		return "", err
	}

	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword{
		return "", err
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil{
		return "", err
	}

	return token, nil
}

// Get Data of logged in user/ user providing token
func GetUser(uid uint) (User, error){

    var user User
    if err := DB.First(&user,uid).Error; err != nil{
        return user, errors.New("User Not Found")
    }

    user.PrepareGive()

    return user, nil
}

// Hide Password in returned data
func (user *User)PrepareGive(){
    user.Password = ""
}
