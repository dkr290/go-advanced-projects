package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

func GenerateToken(email string, userId int64) (string, error) {
    //make sure you use HS256 for signingMethods

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// because the registration is by email address
		"email":  email,
		//this us the user id as well
		"userId": userId,
		//this is the expiration time of the token 2 hours
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		//checking with ok pattern where type checking type if the value stored is in that type check
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method,%v", t.Header["alg"]))
		}
		return []byte(secretKey), nil
	})
	// check for errors
	if err != nil {
		return 0, errors.New("could not parse token "+ err.Error())

	}
     // check if the token is valid
	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userId := int64(claims["userId"].(float64))
	return userId, nil
}
