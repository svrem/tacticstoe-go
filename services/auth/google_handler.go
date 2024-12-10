package auth_service

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	db "tacticstoe/database"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	GivenName string `json:"given_name"`
	Picture   string `json:"picture"`
}

type GoogleHandler struct {
	oauth2conf *oauth2.Config
}

func newGoogleHandler() *GoogleHandler {
	var clientId = os.Getenv("GOOGLE_CLIENT_ID")
	var clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	var redirectUrl = os.Getenv("GOOGLE_REDIRECT_URI")

	println(clientId)
	println(clientSecret)
	println(redirectUrl)

	var oauth2conf = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"profile"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleHandler{
		oauth2conf: oauth2conf,
	}
}

func (gh *GoogleHandler) GetAuthURL() string {
	return gh.oauth2conf.AuthCodeURL("state")
}

func (gh *GoogleHandler) handleCallback(w http.ResponseWriter, r *http.Request) (*db.User, error) {
	code := r.URL.Query().Get("code")

	t, err := gh.oauth2conf.Exchange(context.Background(), code)

	if err != nil {
		return nil, err
	}

	client := gh.oauth2conf.Client(context.Background(), t)

	// get user id
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, err
	}

	var user googleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &db.User{
		ProviderId:     user.Id,
		Provider:       "google",
		Username:       user.GivenName,
		ProfilePicture: user.Picture,
	}, nil
}
