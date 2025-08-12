package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/golang-jwt/jwt/v5"
)

type (
	IJWT interface {
		GenerateToken(adminID, roleName string) (string, string, error)
		ValidateToken(token string) (*jwt.Token, error)
		GetAdminIDByToken(tokenString string) (string, error)
		GetAdminRoleNameByToken(tokenString string) (string, error)
	}

	jwtCustomClaim struct {
		AdminID  string `json:"admin_id"`
		RoleName string `json:"role_name"`
		jwt.RegisteredClaims
	}

	JWT struct {
		secretKey string
		issuer    string
	}
)

func NewJWT() *JWT {
	return &JWT{
		secretKey: getSecretKey(),
		issuer:    "Template",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "Template"
	}

	return secretKey
}

func (j *JWT) GenerateToken(userID, roleName string) (string, string, error) {
	accessClaims := jwtCustomClaim{
		userID,
		roleName,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 300)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", dto.ErrGenerateAccessToken
	}

	refreshClaims := jwtCustomClaim{
		userID,
		roleName,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 3600 * 24 * 7)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", dto.ErrGenerateRefreshToken
	}

	return accessTokenString, refreshTokenString, nil
}

func (j *JWT) parseToken(t_ *jwt.Token) (any, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, dto.ErrUnexpectedSigningMethod
	}

	return []byte(j.secretKey), nil
}

func (j *JWT) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, j.parseToken)
	if err != nil {
		return nil, err
	}

	return token, err
}

func (j *JWT) GetAdminIDByToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", dto.ErrValidateToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", dto.ErrTokenInvalid
	}

	userID := fmt.Sprintf("%v", claims["user_id"])

	return userID, nil
}

func (j *JWT) GetAdminRoleNameByToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", dto.ErrValidateToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", dto.ErrTokenInvalid
	}

	roleName := fmt.Sprintf("%v", claims["role_name"])

	return roleName, nil
}
