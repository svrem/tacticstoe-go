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
		if google_user, err := gh.handleCallback(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			user = google_user
		}

	default:
		http.NotFound(w, r)
		return
	}

	db.CreateUser(database, user)

	csrf_token := uuid.New().String()
	jwt_token := generateJWT(user, csrf_token)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwt_token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(expiration),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrf_token,
		Expires:  time.Now().Add(expiration),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func MeHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	crsf_token := r.Header.Get("X-CSRF-TOKEN")
	jwt_token, err := r.Cookie("jwt")

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := parseJWTToUser(database, jwt_token.Value, crsf_token)

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
