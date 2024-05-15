package keycloak

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

type KeycloakService struct {
	baseUrl             string
	realm               string
	restApiClientId     string
	restApiClientSecret string
	client              *gocloak.GoCloak
}

func NewKeycloakService() *KeycloakService {
	// baseUrl := "http://host.docker.internal:8081"
	// realm := "Element3"
	// restApiClientId := "e3-challenge-server"
	// restApiClientSecret := "Q4xgljWZtLQGj50FTZqRC4rkEkBUfS0u"
	// client := gocloak.NewClient(baseUrl)
	baseUrl := "http://host.docker.internal:8081"
	realm := "master"
	restApiClientId := "e3-keycloak-master"
	restApiClientSecret := "4yC1ofIorC38oIhCgX8x7CXhshbk2jnd"
	client := gocloak.NewClient(baseUrl)

	return &KeycloakService{
		baseUrl:             baseUrl,
		realm:               realm,
		restApiClientId:     restApiClientId,
		restApiClientSecret: restApiClientSecret,
		client:              client,
	}
}

func (ks *KeycloakService) loginRestApiClient(ctx context.Context) (*gocloak.JWT, error) {
	token, err := ks.client.LoginClient(ctx, ks.restApiClientId, ks.restApiClientSecret, ks.realm)
	if err != nil {
		return nil, err
	}
	return token, nil
}


func (ks *KeycloakService) CreateNewClient(ctx context.Context, realmName string, clientName string) (*string, error) {
	token, err := ks.loginRestApiClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain token: %w", err)
	}

	newClient := gocloak.Client{
		ClientID: &clientName,
	}

	client, err := ks.client.CreateClient(ctx, token.AccessToken, realmName, newClient)

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &client, nil

}

func (ks *KeycloakService) CreateNewRealm(ctx context.Context, realmName string) (*string, error) {
	token, err := ks.loginRestApiClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain token: %w", err)
	}

	realmDetails := gocloak.RealmRepresentation{
		Realm: &realmName,
	}

	realm, err := ks.client.CreateRealm(ctx, token.AccessToken, realmDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to create realm: %w", err)
	}

	return &realm, nil
}

func (ks *KeycloakService) CreateUser(ctx context.Context, user gocloak.User) (*gocloak.User, error) {

	token, err := ks.loginRestApiClient(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := ks.client.CreateUser(ctx, token.AccessToken, ks.realm, user)
	if err != nil {
		return nil, err
	}

	userKeycloak, err := ks.client.GetUserByID(ctx, token.AccessToken, ks.realm, userId)
	if err != nil {
		return nil, err
	}

	actions := []string{"VERIFY_EMAIL", "UPDATE_PASSWORD", "UPDATE_PROFILE"} 
	redirectURI := "http://localhost:5173"
	clientID := "e3-challenge-client"
    executeActionsEmail := gocloak.ExecuteActionsEmail{
        UserID:   &userId,
        Lifespan: gocloak.IntP(12 * 60 * 60),
        Actions:  &actions,
		RedirectURI: &redirectURI,
		ClientID: &clientID,
    }
    err = ks.client.ExecuteActionsEmail(ctx, token.AccessToken, ks.realm, executeActionsEmail)
    if err != nil {
        return nil, err 
    }

	return userKeycloak, nil
}

func (ks *KeycloakService) RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {

	rptResult, err := ks.client.RetrospectToken(ctx, accessToken, ks.restApiClientId, ks.restApiClientSecret, ks.realm)
	if err != nil {
		return nil, err
	}
	if !*rptResult.Active {
		return nil, errors.New("token no longer valid")
	}
	return rptResult, nil
}

