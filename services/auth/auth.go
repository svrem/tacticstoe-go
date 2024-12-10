package auth_service

import (
	"encoding/json"
	"net/http"
	db "tacticstoe/database"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	expiration = 7 * 24 * time.Hour
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	switch provider {
	case "google":
		gh := newGoogleHandler()
		http.Redirect(w, r, gh.GetAuthURL(), http.StatusTemporaryRedirect)
	}

	w.WriteHeader(http.StatusNotFound)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	provider := r.PathValue("provider")

	var user *db.User
	switch provider {
	case "google":
		gh := newGoogleHandler()
		if googleUser, err := gh.handleCallback(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			user = googleUser
		}

	default:
		http.NotFound(w, r)
		return
	}

	db.CreateUser(database, user)

	csrfToken := uuid.New().String()
	jwtToken := generateJWT(user, csrfToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(expiration),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(expiration),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func MeHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	crsfToken := r.Header.Get("X-CSRF-TOKEN")
	jwtToken, err := r.Cookie("jwt")

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := parseJWTToUser(database, jwtToken.Value, crsfToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)

}
