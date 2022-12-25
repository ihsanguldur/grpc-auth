package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"grpc-auth/pkg/api/models"
	"time"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	Id       uint
	UserName string
}

type JwtWrapper struct {
	SecretKey string
	Issuer    string
	Expire    int
}

func (w *JwtWrapper) GenerateToken(user *models.User) (signedToken string, err error) {
	claims := jwtClaims{
		Id:       user.ID,
		UserName: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    w.Issuer,
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Duration(w.Expire) * time.Minute)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(w.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (w *JwtWrapper) ValidateToken(signedToken string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(w.SecretKey), nil
		})

	if token.Valid {
		return token.Claims.(*jwtClaims), nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, errors.New("that's not even a token")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, errors.New("timing is everything")
	} else {
		return nil, err
	}
}
