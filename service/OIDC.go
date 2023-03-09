package service

import "golang.org/x/oauth2"

var (
	oauthConfig      *oauth2.Config
	oauthStateString = "random"
)

func init() {
	oauthConfig = &oauth2.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://your-oidc-provider.com/auth",
			TokenURL: "https://your-oidc-provider.com/token",
		},
	}
}
