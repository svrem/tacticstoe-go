package auth_service

import (
	"log/slog"
	"os"
	db "tacticstoe/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func generateJWT(user *db.User, csrf_token string) string {
	claims := jwt.MapClaims{
		"user_id":         user.ID,
		"exp":             time.Now().Add(expiration).Unix(),
		"csrf":            csrf_token,
		"username":        user.Username,
		"profile_picture": user.ProfilePicture,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		slog.Error("Failed to sign token: " + err.Error())
		return ""
	}

	return signedToken
}

func parseJWTToUser(database *gorm.DB, tokenString string, csrf_token string) (*db.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		slog.Error("Failed to parse token: " + err.Error())
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("Failed to parse claims")
		return nil, err
	}

	if claims["csrf"] != csrf_token {
		slog.Error("CSRF token mismatch")
		return nil, err
	}

	user_id := uint(claims["user_id"].(float64))
	user := db.GetUserByID(database, user_id)

	return user, nil
}
