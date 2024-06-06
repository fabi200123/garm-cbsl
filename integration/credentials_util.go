package integration

import (
	"fmt"
	"log/slog"

	"github.com/cloudbase/garm/params"
)

func DeleteGithubCredential(id int64) error {
	slog.Info("Delete GitHub credential")
	if err := deleteGithubCredentials(cli, authToken, id); err != nil {
		return err
	}
	return nil
}

func EnsureTestCredentials(name string, oauthToken string, endpointName string) {
	slog.Info("Ensuring test credentials exist")
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        name,
		Endpoint:    endpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: oauthToken,
		},
	}
	CreateGithubCredentials(createCredsParams)

	createCredsParams.Name = fmt.Sprintf("%s-clone", name)
	CreateGithubCredentials(createCredsParams)
}

func createDummyCredentials(name, endpointName string) (*params.GithubCredentials, error) {
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        name,
		Endpoint:    endpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: "dummy",
		},
	}
	return CreateGithubCredentials(createCredsParams)
}

func CreateGithubCredentials(credentialsParams params.CreateGithubCredentialsParams) (*params.GithubCredentials, error) {
	slog.Info("Create GitHub credentials")
	credentials, err := createGithubCredentials(cli, authToken, credentialsParams)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}
