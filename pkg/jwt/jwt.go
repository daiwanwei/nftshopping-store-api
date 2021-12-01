package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	secretKey      string = "secret"
	issuer         string = "rental"
	ExpirationDays int64  = 7
)

type Claim struct {
	jwt.StandardClaims
}

type UserDetail interface {
	GetAuthorities() []string
	GetName() string
}

func ExtractUserName(tokenString string) (name string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return
	}
	claims := token.Claims.(*Claim)
	name = claims.Subject
	return
}

func ExtractExpiration(tokenString string) (expiration time.Time, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return
	}
	claims := token.Claims.(*Claim)
	expiration = time.Unix(claims.ExpiresAt, 0)
	return
}

func IsTokenExpired(tokenString string) (isExpired bool, err error) {
	expiration, err := ExtractExpiration(tokenString)
	if err != nil {
		return
	}
	if expiration.Before(time.Now()) {
		isExpired = true
	} else {
		isExpired = false
	}
	return
}

func GenerateToken(userDetail UserDetail) (tokenString string, err error) {
	claims := &Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(ExpirationDays)).Unix(),
			Issuer:    issuer,
			Subject:   userDetail.GetName(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return
	}
	return
}

func ValidateToken(tokenString string, tenant UserDetail) (isValid bool, err error) {
	name, err := ExtractUserName(tokenString)
	if err != nil {
		return
	}
	isExpired, err := IsTokenExpired(tokenString)
	if err != nil {
		return
	}
	if name != tenant.GetName() || isExpired {
		isValid = false
	} else {
		isValid = true
	}
	return
}
