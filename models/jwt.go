package models

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/opensourceways/app-robot-server/config"
)

var (
	TokenExpired     = errors.New("Token is expired ")
	TokenMalformed   = errors.New("That's not even a token ")
	TokenNotValidYet = errors.New("Token not active yet ")
	TokenInvalid     = errors.New("Couldn't handle this token: ")
)

type JWT struct {
	SignKey []byte
}

type CustomClaims struct {
	ID   string
	Name string
	jwt.StandardClaims
}

//NewJwt create a JWTConfig object
func NewJwt() *JWT {
	return &JWT{[]byte(config.Application.Jwt.SigningKey)}
}

//CreateToken create a token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SignKey)
}

//ParseToken parse token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SignKey, nil
	})
	if err != nil {
		if ve, ok := err.(jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
		return nil, err
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenInvalid
}

//ParseToken refresh Token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

func genToken(userId, name string) (string, error) {
	j := NewJwt()
	claims := CustomClaims{
		ID:   userId,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: time.Now().Unix() + config.Application.Jwt.TokenExpiration*60*60,
			Issuer:    "app-robot-server",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
