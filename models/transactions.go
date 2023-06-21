package models

import (
    "errors"
    TIME "time"
)

// Create Transactions Table
type Transactions struct{
    ID  uint `json:"-" gorm:"primaryKey"`
    Username    string  `json:"-"`
    Points int  `json:"points"`
    State string    `json:"state"`	// Either EARN or REDEEM
	Description string  `json:"description"`	// Choices: WELCOME BONUS, COMPLETE EMAIL VERIFICATION BONUS, BASIC EARN MPESA TRANSACTION
    Date    string  `json:"date" gorm:"type:date"`
    Time    string  `json:"time" gorm:"type:time"`
    Product string  `json:"redeemed_product,omitempty"`
}
    

// MISC functions
// Save Transactions to Database
func SaveToTransactions(
    username string, 
    state string, 
    description string, 
    points int,
    redeemed_product string){

        var transaction Transactions
        currentTime := TIME.Now()

        transaction.Username = username
        transaction.Points = points
        transaction.State = state
        transaction.Description = description
        transaction.Date = currentTime.Format("01-02-2006") // TODO: convert to string
        // TODO: Get current time
        transaction.Time = currentTime.Format("15:04:05")
        transaction.Product = redeemed_product
        
       if err := DB.Create(&transaction).Error; err != nil{
        logger.Error().Err(errors.New("saving transaction to db failed")).Msgf("%v",err)
       }
}