package main

import (
	"fmt"
	"net/http"
	"strings"
)

func NewRouter(
	homepageURL string,
	statusHandler StatusHandler,
	userHandler UserHandler,
	googleSingleSignOn SingleSignOn,
	facebookSingleSignOn SingleSignOn,
	githubSingleSignOn SingleSignOn,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/status", statusHandler.GetStatus)
	mux.HandleFunc("/api/v1/me", userHandler.getSignedInUser)

	googleSingleSignOnHandler := NewSingleSignOnHandler(googleSingleSignOn, homepageURL)
	mux.HandleFunc("/api/v1/single-sign-on/google/sign-in", googleSingleSignOnHandler.SignIn)
	mux.HandleFunc("/api/v1/single-sign-on/google/callback", googleSingleSignOnHandler.Callback)

	facebookSingleSignOnHandler := NewSingleSignOnHandler(facebookSingleSignOn, homepageURL)
	mux.HandleFunc("/api/v1/single-sign-on/facebook/sign-in", facebookSingleSignOnHandler.SignIn)
	mux.HandleFunc("/api/v1/single-sign-on/facebook/callback", facebookSingleSignOnHandler.Callback)

	githubSingleSignOnHandler := NewSingleSignOnHandler(githubSingleSignOn, homepageURL)
	mux.HandleFunc("/api/v1/single-sign-on/github/sign-in", githubSingleSignOnHandler.SignIn)
	mux.HandleFunc("/api/v1/single-sign-on/github/callback", githubSingleSignOnHandler.Callback)

	return mux
}

type StatusHandler struct{}

func NewStatusHandler() StatusHandler {
	return StatusHandler{}
}

func (h StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	rsp := struct {
		Status string `json:"status"`
	}{Status: "OK"}
	HttpReplyJson(w, http.StatusOK, rsp)
}

type UserHandler struct {
	authenticator Authenticator
	userManager   UserManager
}

func NewUserHandler(
	authenticator Authenticator,
	userManager UserManager,
) UserHandler {
	return UserHandler{
		authenticator: authenticator,
		userManager:   userManager,
	}
}

func (h UserHandler) getSignedInUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	token := getToken(r)
	userID, err := h.authenticator.GetUserID(token)
	if err != nil {
		err = fmt.Errorf("could not authorize user: %v", err)
		HttpReplyError(w, http.StatusUnauthorized, err)
		return
	}

	user, err := h.userManager.GetUserByID(userID)
	if err != nil {
		err = fmt.Errorf("could not retrieve authorized user: %v", err)
		HttpReplyError(w, http.StatusUnauthorized, err)
		return
	}

	rsp := struct {
		User User `json:"user"`
	}{User: user}
	HttpReplyJson(w, http.StatusOK, rsp)
}

type SingleSignOnHandler struct {
	singleSignOn SingleSignOn
	homepageURL  string
}

func NewSingleSignOnHandler(
	singleSignOn SingleSignOn,
	homepageURL string,
) SingleSignOnHandler {
	return SingleSignOnHandler{
		singleSignOn: singleSignOn,
		homepageURL:  homepageURL,
	}
}

func (h SingleSignOnHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	token := getToken(r)
	if h.singleSignOn.IsSignedIn(token) {
		http.Redirect(w, r, h.homepageURL, http.StatusSeeOther)
		return
	}

	authorizationURL, err := h.singleSignOn.GetAuthorizationURL()
	if err != nil {
		err = fmt.Errorf("invalid authorization url: %v", err)
		HttpReplyError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, authorizationURL, http.StatusSeeOther)
}

func (h SingleSignOnHandler) Callback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.singleSignOn.SignIn(code)
	if err != nil {
		err = fmt.Errorf("invalid authorization code: %v", err)
		HttpReplyError(w, http.StatusBadRequest, err)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/"})
	http.Redirect(w, r, h.homepageURL, http.StatusSeeOther)
}

func getToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) < 1 {
		return ""
	}
	words := strings.Split(auth, " ")
	if len(words) != 2 {
		return ""
	}
	if words[0] != "Bearer" {
		return ""
	}
	return words[1]
}
