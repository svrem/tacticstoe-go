package auth

import (
	"encoding/json"
	"net/http"
	db "tacticstoe/internal/database"
	"time"

	"gorm.io/gorm"
)

const (
	expiration = 7 * 24 * time.Hour
)

type AuthUser struct {
	Username       string
	ProfilePicture string
	Provider       string
	ProviderId     string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	switch provider {
	case "google":
		gh := newGoogleHandler()
		http.Redirect(w, r, gh.GetAuthURL(), http.StatusTemporaryRedirect)
	}

	w.WriteHeader(http.StatusNotFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	provider := r.PathValue("provider")

	var authUser *AuthUser
	switch provider {
	case "google":
		gh := newGoogleHandler()
		if googleUser, err := gh.handleCallback(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			authUser = googleUser
		}

	default:
		http.NotFound(w, r)
		return
	}

	user, err := db.CreateUser(database, authUser.Username, authUser.ProfilePicture, authUser.Provider, authUser.ProviderId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jwtToken := generateJWT(user)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(expiration),
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func AutherizeUser(w http.ResponseWriter, r *http.Request, database *gorm.DB) *db.User {
	jwtToken, err := r.Cookie("token")

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	user := parseJWTToUser(database, jwtToken.Value)

	if user == nil {
		// delete cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Unix(0, 0),
			Path:     "/",
		})

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	return user
}

func MeHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	user := AutherizeUser(w, r, database)

	if user == nil {
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
