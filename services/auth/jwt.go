package auth_service

import (
	"log/slog"
	"os"
	db "tacticstoe/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func generateJWT(user *db.User) string {
	claims := jwt.MapClaims{
		"user_id":         user.ID,
		"exp":             time.Now().Add(expiration).Unix(),
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

func parseJWTToUser(database *gorm.DB, tokenString string) *db.User {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		slog.Error("Failed to parse token: " + err.Error())
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error("Failed to parse claims")
		return nil
	}

	userId := claims["user_id"].(string)

	user := db.GetUserByID(database, userId)

	return user
}
