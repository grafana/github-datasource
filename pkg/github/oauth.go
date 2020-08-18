package github

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func GetOAuthConfig(clientID string, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"repo", "read:org"},
		Endpoint:     github.Endpoint,
	}
}
