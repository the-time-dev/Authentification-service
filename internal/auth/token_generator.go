package auth

import (
	"auth-service/internal/storage"
	"encoding/hex"
	"fmt"
	"time"

	"crypto/sha256"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("supersecretkey")

type Authorizer struct {
	storage *storage.Storage
}

func NewAuthorizer(storage *storage.Storage) *Authorizer {
	return &Authorizer{storage: storage}
}

func (a *Authorizer) GenerateTokenPair(userID, clientIP string) (string, string, error) {

	tokenUUID := uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":       userID,
		"ip":        clientIP,
		"tokenUUID": tokenUUID,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
	})
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s.%s", userID, time.Now().String())))
	refreshToken := hex.EncodeToString(hash[:])
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	err = a.storage.AddToken(tokenUUID, string(hashedRefreshToken))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshToken, nil
}

func (a *Authorizer) ValidateAccessToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return false, err
	}

	return true, nil
}

func (a *Authorizer) RefreshAccessToken(accessToken, refreshToken string) (string, string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	tokenUUID := claims["tokenUUID"].(string)

	hashToken, err := a.storage.GetHash(tokenUUID)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashToken), []byte(refreshToken))
	if err != nil {
		return "", "", err
	}

	claims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s.%s", claims["sub"].(string), time.Now().String())))
	refreshToken = hex.EncodeToString(hash[:])
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	err = a.storage.UpdateToken(tokenUUID, string(hashedRefreshToken))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
