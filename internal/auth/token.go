package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	GitHubUser string `json:"github_user"`
	jwt.RegisteredClaims
}

// ParseJWTToken parses and validates a JWT token
func ParseJWTToken(tokenString string, secretKey []byte) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// IsTokenExpired checks if a JWT token is expired
func IsTokenExpired(tokenString string) (bool, error) {
	// Parse without verification to check expiration
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return true, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return true, fmt.Errorf("invalid claims")
	}

	if claims.ExpiresAt == nil {
		// No expiration time set
		return false, nil
	}

	return claims.ExpiresAt.Before(time.Now()), nil
}

// RefreshJWTToken refreshes a JWT token using the Portal API refresh endpoint
// This is a placeholder implementation - the actual refresh logic would depend on the Portal API
func RefreshJWTToken(refreshToken string, portalURL string) (string, error) {
	// TODO: Implement JWT token refresh logic
	// This would typically involve:
	// 1. Making a POST request to the Portal API refresh endpoint
	// 2. Sending the refresh token
	// 3. Receiving a new JWT token
	return "", fmt.Errorf("JWT token refresh not yet implemented")
}
