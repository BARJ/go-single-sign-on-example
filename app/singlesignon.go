package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type SingleSignOnUser struct {
	Email   string
	Name    string
	Picture string
}

type IdentityProvider interface {
	GetAuthorizationURL() (string, error)
	GetIdentityToken(code string) (string, error)
	GetSingleSignOnUser(identityToken string) (SingleSignOnUser, error)
}

type SingleSignOn struct {
	identityProvider IdentityProvider
	authenticator    Authenticator
	repository       Repository
}

func NewSingleSignOn(
	identityProvider IdentityProvider,
	authenticator Authenticator,
	repository Repository,
) SingleSignOn {
	return SingleSignOn{
		identityProvider: identityProvider,
		authenticator:    authenticator,
		repository:       repository,
	}
}

func (s SingleSignOn) GetAuthorizationURL() (string, error) {
	return s.identityProvider.GetAuthorizationURL()
}

func (s SingleSignOn) IsSignedIn(token string) bool {
	_, err := s.authenticator.GetUserID(token)
	return err == nil
}

func (s SingleSignOn) SignIn(code string) (string, error) {
	if len(code) < 1 {
		return "", errors.New("authorization code cannot be empty")
	}

	identityToken, err := s.identityProvider.GetIdentityToken(code)
	if err != nil {
		return "", err
	}

	singleSignOnUser, err := s.identityProvider.GetSingleSignOnUser(identityToken)
	if err != nil {
		return "", err
	}

	user, err := s.getOrCreateUser(singleSignOnUser)
	if err != nil {
		return "", err
	}

	return s.authenticator.CreateToken(user.ID)
}

func (s SingleSignOn) getOrCreateUser(singleSignOnUser SingleSignOnUser) (User, error) {
	if len(singleSignOnUser.Email) < 1 {
		return User{}, errors.New("email cannot be empty")
	}

	user, err := s.repository.GetUserByEmail(singleSignOnUser.Email)
	if err == nil {
		return user, nil
	}
	var userNotFound ErrUserNotFound
	if !errors.As(err, &userNotFound) {
		return User{}, err
	}

	return s.repository.CreateUser(User{
		Email:   singleSignOnUser.Email,
		Name:    singleSignOnUser.Name,
		Picture: singleSignOnUser.Picture,
	})
}

var _ IdentityProvider = (*GoogleIdentityProvider)(nil)

type GoogleIdentityProvider struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

func NewGoogleIdentityProvider(
	clientID string,
	clientSecret string,
	redirectURI string,
) GoogleIdentityProvider {
	return GoogleIdentityProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

func (g GoogleIdentityProvider) GetAuthorizationURL() (string, error) {
	u, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		return "", err
	}

	q := url.Values{}
	q.Add("client_id", g.clientID)
	q.Add("redirect_uri", g.redirectURI)
	q.Add("response_type", "code")
	q.Add("scope", "email profile")
	q.Add("access_type", "online")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (g GoogleIdentityProvider) GetIdentityToken(code string) (string, error) {
	u, err := url.Parse("https://oauth2.googleapis.com/token")
	if err != nil {
		return "", err
	}

	body := url.Values{}
	body.Add("client_id", g.clientID)
	body.Add("client_secret", g.clientSecret)
	body.Add("code", code)
	body.Add("grant_type", "authorization_code")
	body.Add("redirect_uri", g.redirectURI)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	res := struct {
		AccessToken string `json:"access_token"`
	}{}
	if err := DefaultHttpClient.HttpRequestJson(http.MethodPost, u.String(), headers, body.Encode(), &res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

func (g GoogleIdentityProvider) GetSingleSignOnUser(identityToken string) (SingleSignOnUser, error) {
	u, err := url.Parse("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return SingleSignOnUser{}, err
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", identityToken),
	}

	res := struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}{}

	if err := DefaultHttpClient.HttpRequestJson(http.MethodGet, u.String(), headers, "", &res); err != nil {
		return SingleSignOnUser{}, err
	}

	return SingleSignOnUser{
		Name:    res.Name,
		Email:   res.Email,
		Picture: res.Picture,
	}, nil
}

var _ IdentityProvider = (*FacebookIdentityProvider)(nil)

type FacebookIdentityProvider struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

func NewFacebookIdentityProvider(
	clientID string,
	clientSecret string,
	redirectURI string,
) FacebookIdentityProvider {
	return FacebookIdentityProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

func (f FacebookIdentityProvider) GetAuthorizationURL() (string, error) {
	u, err := url.Parse("https://www.facebook.com/v9.0/dialog/oauth")
	if err != nil {
		return "", err
	}

	q := url.Values{}
	q.Add("client_id", f.clientID)
	q.Add("scope", "public_profile,email")
	q.Add("redirect_uri", f.redirectURI)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (f FacebookIdentityProvider) GetIdentityToken(code string) (string, error) {
	u, err := url.Parse("https://graph.facebook.com/v9.0/oauth/access_token")
	if err != nil {
		return "", err
	}

	q := url.Values{}
	q.Add("client_id", f.clientID)
	q.Add("client_secret", f.clientSecret)
	q.Add("code", code)
	q.Add("redirect_uri", f.redirectURI)
	u.RawQuery = q.Encode()

	res := struct {
		AccessToken string `json:"access_token"`
	}{}
	if err := DefaultHttpClient.HttpRequestJson(http.MethodGet, u.String(), nil, "", &res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

func (f FacebookIdentityProvider) GetSingleSignOnUser(identityToken string) (SingleSignOnUser, error) {
	u, err := url.Parse("https://graph.facebook.com/v9.0/me")
	if err != nil {
		return SingleSignOnUser{}, err
	}

	q := url.Values{}
	q.Add("access_token", identityToken)
	q.Add("fields", "name,email,picture")
	u.RawQuery = q.Encode()

	res := struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}{}

	if err := DefaultHttpClient.HttpRequestJson(http.MethodGet, u.String(), nil, "", &res); err != nil {
		return SingleSignOnUser{}, err
	}

	return SingleSignOnUser{
		Name:    res.Name,
		Email:   res.Email,
		Picture: res.Picture.Data.URL,
	}, nil
}

var _ IdentityProvider = (*GithubIdentityProvider)(nil)

type GithubIdentityProvider struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

func NewGithubIdentityProvider(
	clientID string,
	clientSecret string,
	redirectURI string,
) GithubIdentityProvider {
	return GithubIdentityProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

func (g GithubIdentityProvider) GetAuthorizationURL() (string, error) {
	u, err := url.Parse("https://github.com/login/oauth/authorize")
	if err != nil {
		return "", err
	}

	q := url.Values{}
	q.Add("client_id", g.clientID)
	q.Add("redirect_uri", g.redirectURI)
	q.Add("scope", "read:user")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (g GithubIdentityProvider) GetIdentityToken(code string) (string, error) {
	u, err := url.Parse("https://github.com/login/oauth/access_token")
	if err != nil {
		return "", err
	}

	body := url.Values{}
	body.Add("client_id", g.clientID)
	body.Add("client_secret", g.clientSecret)
	body.Add("code", code)
	body.Add("redirect_uri", g.redirectURI)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	res := struct {
		AccessToken string `json:"access_token"`
	}{}
	if err := DefaultHttpClient.HttpRequestJson(http.MethodPost, u.String(), headers, body.Encode(), &res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

func (g GithubIdentityProvider) GetSingleSignOnUser(identityToken string) (SingleSignOnUser, error) {
	u, err := url.Parse("https://api.github.com/user")
	if err != nil {
		return SingleSignOnUser{}, err
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("token %s", identityToken),
	}

	res := struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"avatar_url"`
	}{}

	if err := DefaultHttpClient.HttpRequestJson(http.MethodGet, u.String(), headers, "", &res); err != nil {
		return SingleSignOnUser{}, err
	}

	return SingleSignOnUser{
		Name:    res.Name,
		Email:   res.Email,
		Picture: res.Picture,
	}, nil
}

type SingleSignOnFactory struct {
	authenticator Authenticator
	repository    Repository
}

func NewSingleSignOnFactory(
	authenticator Authenticator,
	repository Repository,
) SingleSignOnFactory {
	return SingleSignOnFactory{
		authenticator: authenticator,
		repository:    repository,
	}
}

func (f SingleSignOnFactory) NewSingleSignOn(identityProvider IdentityProvider) SingleSignOn {
	return NewSingleSignOn(
		identityProvider,
		f.authenticator,
		f.repository,
	)
}
