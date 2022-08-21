package utils

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SignedDetails struct {
	Id string
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret")

func Issue(id string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "localJWT",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "cookieJWT",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 168)),
		},
	}

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)
		return
	}

	return token, refreshToken, err
}

func Verify(signedToken string) (claims *SignedDetails, returnStatus int, errMsg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	claims, response := token.Claims.(*SignedDetails)
	if !response {
		errMsg = "reponse error"
		returnStatus = http.StatusBadRequest
		return
	}

	if claims.Issuer == "cookieJWT" {
		errMsg = "Can't use the a JWT Cookie to retrieve data"
		returnStatus = http.StatusBadRequest
		return
	}

	if err != nil {
		errMsg = err.Error()
		// gucci
		if strings.Contains(errMsg, "expired") && claims.Issuer == "localJWT" {
			errMsg = "Token has expired."
			returnStatus = http.StatusBadRequest
		}
		return
	}

	return claims, http.StatusOK, errMsg
}
