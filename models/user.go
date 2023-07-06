package models

import (
    "fmt"
	"time"
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User data structure
type User struct {
    gorm.Model
	Username string `gorm:"unique" json:"username"`
    FirstName string `json:"firstname"`
    LastName string `json:"lastname"`
	Gender string `json:"gender"`	
    EmailAddress string `json:"email"`
    Password string `json:"password"`
    City string `json:"city"`
    PhoneNumber string `json:"phone_number"`
    RedeemablePoints int `json:"redeemable_points"`
    IsEmailVerified bool `json:"is_email_verified" gorm:"not null; default:false"`
    VirtualCardNumber string `json:"card_number" gorm:"unique"` 
}	

// User Creation Request Data
type CreateUserInput struct {
    Username string `json:"username" validate:"required"`
    FirstName string `json:"firstname" validate:"required"`
    LastName string `json:"lastname" validate:"required"`
    Gender string `json:"gender" validate:"required"`
    EmailAddress string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
    City string `json:"city" validate:"required"`
    PhoneNumber string `json:"phone_number" validate:"required"`    
}

type VerifyEmails struct {
    gorm.Model
    Username     string         `gorm:"unique" json:"username"`
    EmailAddress string         `json:"email"`
    SecretCode   string         `json:"secret_code"`
    ExpiredAt    time.Time `json:"expired_at" gorm:"default: (now() + interval '15 minutes')"`
    User         User      `gorm:"foreignKey:Username"`
}

type PasswordReset struct{
    gorm.Model
    Username     string         `gorm:"unique" json:"username"`
    EmailAddress string         `json:"email"`
    PasswordResetCode   string         `json:"reset_code"`
    ExpiredAt    time.Time `json:"expired_at" gorm:"default: (now() + interval '15 minutes')"`
    User         User      `gorm:"foreignKey:Username"`
}

// User Login Request data
type LoginInput struct {
	EmailAddress    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

// User Data Response
type UserResponse struct {
	Username      string    `json:"username,omitempty"`
    FirstName   string  `json:"first_name,omitempty"`
    LastName    string  `json:"last_name,omitempty"`
	EmailAddress    string    `json:"email,omitempty"`
    PhoneNumber     string  `json:"phone_number,omitempty"`
    RedeemablePoints  int  `json:"redeemable_points"`
    City    string  `json:"city,omitempty"`
    VirtualCardNumber string `json:"card_number"` 
}

// ForgotPasswordInput
type ForgotPasswordInput struct{
    EmailAddress string `json:"email" validate:"required, email"`
}

// ResetPasswordInput
type ResetPasswordInput struct{
    NewPassword string  `json:"new_password" validate:"required"`
    ConfirmPassword string  `json:"confirm_password" validate:"required"`
}


// MISC FUNCTIONS
// Save a new user record to the database
func (user *User) SaveUser() (*User, error){

	var err error = DB.Create(&user).Error
	if err != nil {
        logger.Error().Err(errors.New("saving user to database failed")).Msgf("%v",err)
		return nil, err
	}
	return &User{}, nil
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

// Function to update the user's password in the database
func UpdateUserPassword(userId uint, hashedPassword string) error {

    // Find the user by id
	var user User
	if err := DB.First(&user, userId).Error; err != nil {
		return err
	}

	// Update the user's password
	user.Password = string(hashedPassword)
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
