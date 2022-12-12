package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/pkg/config"
	jwt "github.com/dgrijalva/jwt-go"
)

var ErrInvalidToken = errors.New("invalid token")

type AuthenticatedClaim struct {
	Authorized bool   `json:"authorized"`
	CompanyId  uint32 `json:"companyId"`
	UserId     uint32 `json:"userId"`
	Email      string `json:"email"`
	Purpose    string `json:"purpose"`
	ExpiryTime int64  `json:"exp"`
}

type AuthenticationResponse struct {
	AccessToken  string     `json:"accessToken"`
	RefreshToken string     `json:"refreshToken"`
	User         model.User `json:"user"`
}

func GenerateAuthenticationResponse(user *model.User,
	accessTimeoutDuration time.Duration, refreshTimeoutDuration time.Duration) (response AuthenticationResponse, err error) {
	t := time.Now()
	claim := AuthenticatedClaim{
		Authorized: true,
		CompanyId:  user.CompanyId,
		UserId:     user.ID,
		Email:      user.Email,
		Purpose:    "access",
		ExpiryTime: t.Add(accessTimeoutDuration).Unix(),
	}
	token, err := CreateToken(claim)
	if err != nil {
		return
	}
	response.AccessToken = token

	claim.Purpose = "refresh"
	claim.ExpiryTime = t.Add(refreshTimeoutDuration).Unix()
	token, err = CreateToken(claim)
	if err != nil {
		return
	}
	response.RefreshToken = token
	response.User = *user

	return
}

func CreateToken(claim AuthenticatedClaim) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = claim.Authorized
	claims["userId"] = claim.UserId
	claims["companyId"] = claim.CompanyId
	claims["email"] = claim.Email
	claims["purpose"] = claim.Purpose
	claims["exp"] = claim.ExpiryTime

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.ServerConfiguration.Application.JWTSecret))
}

func ValidateToken(r *http.Request) (response *AuthenticatedClaim, err error) {
	tokenString := ExtractToken(r)

	return ValidateTokenString(tokenString)
}

func ValidateTokenString(tokenString string) (response *AuthenticatedClaim, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ServerConfiguration.Application.JWTSecret), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		response, err = parseAuthenticatedClaim(claims)
	} else {
		err = ErrInvalidToken
	}

	return
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// parseAuthenticatedClaim display the claims licely in the terminal
func parseAuthenticatedClaim(data interface{}) (response *AuthenticatedClaim, err error) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return
	}

	var claim AuthenticatedClaim
	err = json.Unmarshal(b, &claim)
	if err == nil {
		response = &claim
	}

	return
}
