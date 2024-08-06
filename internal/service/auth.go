package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/GlebKirsan/go-final-project/internal/config"
	"github.com/GlebKirsan/go-final-project/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	alg jwt.SigningMethod
}

func NewAuthService() *AuthService {
	return &AuthService{
		alg: jwt.SigningMethodHS256,
	}
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *AuthService) Authorize(auth *models.AuthRequest) (string, error) {
	cfg := config.Get()
	pass := cfg.Pass
	if auth.Password != pass {
		return "", errors.New("wrong password")
	}

	claims := jwt.MapClaims{
		"hash": GetMD5Hash(pass),
	}

	jwtToken := jwt.NewWithClaims(s.alg, claims)
	token, err := jwtToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
