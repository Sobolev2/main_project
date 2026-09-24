package repository

import (
	"time"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"errors"
	"semen_project/internal/models"
)
type AuthRepo interface {
	GetUserByUsername(username string) (models.UserPublic, error)
	GetPasswordById(id int) (string, error)
	CreateUser(username, firstName, lastName, password string) (*models.UserPublic, error)
	DeleteAllUserRefreshTokensExceptThis(userID int, tokenHash string) error
	CheckPassword(hash string, password string) error


	SaveRefreshToken(userID int, tokenHash string) error
	GetRefreshTokenByTokenHash(hash string) (models.RefreshToken, error)

	DeleteRefreshTokenByTokenHash(tokenHash string) error
	GenerateAccessToken(UserID int, secret string) (string, error)
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
func (s *Store) GenerateAccessToken(UserID int, secret string) (string, error) {
	claims := Claims{
	UserID: UserID,
	RegisteredClaims: jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	},
	}
	secretKey := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token_string, nil
}
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func ValidateAccessToken(TokenString string, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
    TokenString,
    claims,
    func(token *jwt.Token) (interface{}, error) {

        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }

        return []byte(secret), nil
    },
)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
	    return nil, errors.New("invalid token")
	}
	return claims, nil
}
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}


