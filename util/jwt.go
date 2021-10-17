package util

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// GetUserAndFirmIDFromJWT parse jwt and take userid and firmid
func GetUserAndFirmIDFromJWT(r *http.Request) (string, string, error) {
	var userID string
	var firmID string
	var ok bool

	tokenStringBearer := r.Header.Get("Authorization")
	tokenString := tokenStringBearer[7:]
	token, err := jwt.Parse(tokenString, nil)
	if token == nil {
		return userID, firmID, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return firmID, firmID, errors.New("Failed to get Authorization token")
	}

	userID, ok = claims["sub"].(string)
	if !ok {
		return userID, firmID, errors.New("Failed to get id from token")
	}

	firmID, ok = claims["firm_id"].(string)
	if !ok {
		return userID, firmID, errors.New("Failed to get firm id from token")
	}

	return userID, firmID, nil
}

// GetUserAndFirmIDFromJWT parse jwt and take userid and firmid
func GetUserFromJWT(r *http.Request) (string, error) {
	var userID string
	var ok bool

	tokenStringBearer := r.Header.Get("Authorization")
	tokenString := tokenStringBearer[7:]
	token, err := jwt.Parse(tokenString, nil)
	if token == nil {
		return userID, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return userID, errors.New("Failed to get Authorization token")
	}

	userID, ok = claims["sub"].(string)
	if !ok {
		return userID, errors.New("Failed to get id from token")
	}
	return userID, nil
}
