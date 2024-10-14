package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yervsil/auth_service/domain"
)

var jwtKey = []byte("n3T@9fhY*klZ23%*vQJY6p@7sWx8Q9xF")

type TokenPair struct {
	RefreshToken string    `json:"refreshToken"`
	AccessToken string    `json:"accessToken"`
}

func Token(user *domain.User) (*TokenPair, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": user.Id,
		"email": user.Email,
        "name": user.Name,
		"exp":	time.Now().Add(3 * time.Hour).Unix(),
	})

    accessToken, err := at.SignedString(jwtKey)
    if err != nil {
        return nil, err
    }
    fmt.Println(user)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	    "uuid": user.Id,
		"email":  user.Email,
        "name": user.Name,
		"exp":	time.Now().Add(72 * time.Hour).Unix(),
	})

    refreshToken, err := rt.SignedString(jwtKey)
    if err != nil {
        return nil, err
    }

	var tp TokenPair

	tp.AccessToken = accessToken
	tp.RefreshToken = refreshToken

	return &tp, nil
}


func ParseToken(refreshToken string) (*domain.User, error) {
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unsupported signing type: %v", token.Header["alg"])
        }
        return []byte(jwtKey), nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return nil, errors.New("token is invalid")
    }
    fmt.Println(claims)
    return &domain.User{Id: int(claims["uuid"].(float64)), 
                        Email: claims["email"].(string), 
                        Name: claims["name"].(string)}, nil
} 