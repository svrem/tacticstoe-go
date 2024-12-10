package auth_service

import "net/http"

type AuthUser struct {
	ProviderId    string `json:"provider_id"`
	OAuthProvider string `json:"provider"`

	Username       string `json:"user_name"`
	ProfilePicture string `json:"profile_picture"`
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

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")

	var user *AuthUser
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

	println(user.Username)

	w.WriteHeader(http.StatusNotFound)
}
