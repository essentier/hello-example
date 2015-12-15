package authutil

import (
	"crypto/rsa"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
)

func PublicKey() *rsa.PublicKey {
	return publicKey
}

func GetSubjectInToken(tokenString string) (subject string, err error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	sub, exists := token.Claims["sub"]
	if !exists {
		return "", errors.New("The token string does not contain the subject.")
	}

	return sub.(string), nil
}

func ParseToken(tokenString string) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, jwtKeyFunc)
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	log.Printf("token: %#v", token)
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		log.Printf("Unexpected signing method: %v", token.Header["alg"])
		return nil, errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
	} else {
		log.Printf("okay")
		return publicKey, nil
	}
}
