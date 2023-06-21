/*
Utility function to obtain the access token from the Google
OAuth2 token endpoint
*/

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"github.com/rs/zerolog"
	"github.com/joho/godotenv"
)

// Define Logger Instance
var logger = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()

type GoogleOAuthToken struct{
	Access_token string
	Id_token string
}

type GoogleUserResult struct{
	Username string
	FirstName string
	LastName string
	Gender string
	EmailAddress string
	Password string
	City string
	PhoneNumber string
	IsEmailVerified bool
}


// retrieve OAuth2 Access Token
func GetGoogleOAuthToken(code string) (*GoogleOAuthToken, error) {
	rootURL := "https://oauth2.googleapis.com/token"

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error().Err(errors.New("reading environment variables failed")).Msgf("%v",err)
	}

	values := url.Values{}
	values.Set("code", code)
	values.Set("client_id", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	values.Set("client_secret", os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	values.Set("redirect_uri", os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"))
	values.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", rootURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve token: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var token GoogleOAuthToken
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Get the Google User's Account Information
func GetGoogleUser(access_token string, id_token string) (*GoogleUserResult, error){
	rootURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", access_token)

	req, err := http.NewRequest("GET", rootURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", id_token))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &GoogleUserResult{
		Username: GoogleUserRes["name"].(string),
		FirstName: GoogleUserRes["given_name"].(string),
		LastName: GoogleUserRes["family_name"].(string),
		Gender: GoogleUserRes["gender"].(string),
		EmailAddress: GoogleUserRes["email"].(string),
		Password: "",
		City: GoogleUserRes["city"].(string),
		PhoneNumber: GoogleUserRes["phone_number"].(string),
		IsEmailVerified: GoogleUserRes["verified_email"].(bool),
	}

	return userBody, nil
}
