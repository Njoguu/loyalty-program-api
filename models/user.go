package models

import (
	"errors"
	"fmt"
	// "html"
	// "strings"
	"time"

	// token "api/utils/token"
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
    IsEmailVerified bool `json:"is_email_verified" gorm:"not null; default:false"`
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

type VerifyEmails struct {
    gorm.Model
    Username     string         `gorm:"unique" json:"username"`
    EmailAddress string         `json:"email"`
    SecretCode   string         `json:"secret_code"`
    ExpiredAt    time.Time `json:"expired_at" gorm:"default: (now() + interval '15 minutes')"`
    User         User      `gorm:"foreignKey:Username"`
}

type LoginInput struct {
	EmailAddress    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type Points struct {
    // gorm.Model
    UserID   uint `gorm:"unique" json:"userid"`
    Username string `gorm:"unique" json:"username"`
    Points   int32 `json:"points"`
}

type UserResponse struct {
	// ID        uuid.UUID `json:"id,omitempty"`
	Username      string    `json:"name,omitempty"`
    FirstName   string  `json:"first_name,omitempty"`
    LastName    string  `json:"last_name,omitempty"`
	EmailAddress    string    `json:"email,omitempty"`
    IsEmailVerified bool    `json:"is_email_verified,omitempty"`
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

// Verify password hash and password given are identical
func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

// Hash Password function
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
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

// On Account creation, user is allocated 500 points
func AllocatePoints(uid uint, username string) (User, Points, error){

    points := Points{}
    user := User{}

    points.UserID = uid
    points.Username = username
    points.Points = 500

    if err := DB.Create(&points).Error; err != nil {
        return user, points, err
    }

    return user, points, nil
}

func GetPointsByID(uid uint) (Points, error){
    var userPointsData Points
    if err := DB.First(&userPointsData).Error; err != nil{
        return userPointsData, errors.New("could not find user points data")
    }

    return userPointsData, nil
}