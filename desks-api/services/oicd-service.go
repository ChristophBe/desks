package services

import (
	"context"
	"fmt"
	"github.com/ChristophBe/desks/desks-api/util/configuration"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
)

type OICDService interface {
	GetRedirectUrl(state string) string
	VerifyToken(ctx context.Context, tokenString string) (*oidc.IDToken, error)
	GetIdToken(ctx context.Context, code string) (*oidc.IDToken, error)
}

func NewOICDService(ctx context.Context) OICDService {
	service := new(oicdServiceImpl)
	var err error
	service.provider, err = oidc.NewProvider(ctx, configuration.OICDIssuerUrl.GetValue())
	if err != nil {
		log.Fatal("can not resolve OICD Issuer", err)
	}
	service.oauthConfig = &oauth2.Config{
		ClientID:     configuration.OauthClientId.GetValue(),
		ClientSecret: configuration.OauthClientSecret.GetValue(),
		RedirectURL:  fmt.Sprintf("%s/auth/token/", configuration.BaseUrl.GetValue()),
		Endpoint:     service.provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups", "family_name", "given_name"},
	}

	service.idTokenVerifier = service.provider.Verifier(&oidc.Config{ClientID: configuration.OauthClientId.GetValue()})

	return service
}

type oicdServiceImpl struct {
	provider        *oidc.Provider
	oauthConfig     *oauth2.Config
	idTokenVerifier *oidc.IDTokenVerifier
}

func (o oicdServiceImpl) GetIdToken(ctx context.Context, code string) (*oidc.IDToken, error) {
	oauth2Token, err := o.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange oauth token")
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("failed to get idToken from oauth token")
	}

	return o.VerifyToken(ctx, rawIDToken)

}

func (o oicdServiceImpl) GetRedirectUrl(state string) string {
	return o.oauthConfig.AuthCodeURL(state)
}

func (o oicdServiceImpl) VerifyToken(ctx context.Context, tokenString string) (token *oidc.IDToken, err error) {
	idToken, err := o.idTokenVerifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("could not verify bearer token: %v", err)
	}
	return idToken, nil
}
