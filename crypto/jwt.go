package crypto

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_TOKEN_CHARSET = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const JWT_TOKEN_LENGTH = 64
const JWT_TOKEN_EXPIRY = 12 * time.Hour

type JWTClaims struct {
	Origin string `json:"origin"`
	jwt.RegisteredClaims
}

func JWTGenerateTokenSecret() string {
	result := make([]rune, JWT_TOKEN_LENGTH)
	for i := range result {
		result[i] = JWT_TOKEN_CHARSET[rand.Intn(len(JWT_TOKEN_CHARSET))]
	}

	return string(result)
}

func JWTGenerateToken(secret string, origin string, subject string) (token string, expires time.Time, err error) {
	expires = time.Now().Add(JWT_TOKEN_EXPIRY)

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		origin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			Subject:   subject,
		},
	})

	token, err = jwt.SignedString([]byte(secret))
	return
}

func JWTValidateToken(secret string, token string) (valid bool, subject string, err error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return false, "", err
	}

	if !parsed.Valid {
		return false, "", nil
	}

	sub, err := parsed.Claims.GetSubject()
	if err != nil {
		return false, "", err
	}

	return true, sub, nil
}
