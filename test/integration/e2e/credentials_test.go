package e2e

import (
	"fmt"
	"testing"

	"github.com/cloudbase/garm/params"
)

func EnsureTestCredentials(t *testing.T, name string, oauthToken string, endpointName string) {
	t.Log("Ensuring test credentials exist")
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        name,
		Endpoint:    endpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: oauthToken,
		},
	}
	CreateGithubCredentials(t, createCredsParams)

	createCredsParams.Name = fmt.Sprintf("%s-clone", name)
	CreateGithubCredentials(t, createCredsParams)
}

func createDummyCredentials(t *testing.T, name, endpointName string) *params.GithubCredentials {
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        name,
		Endpoint:    endpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: "dummy",
		},
	}
	return CreateGithubCredentials(t, createCredsParams)
}

func TestGithubCredentialsErrorOnDuplicateCredentialsName(t *testing.T) {
	t.Log("Testing error on duplicate credentials name")
	creds := createDummyCredentials(t, dummyCredentialsName, defaultEndpointName)
	defer DeleteGithubCredential(t, int64(creds.ID))

	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: "dummy",
		},
	}
	if _, err := createGithubCredentials(cli, authToken, createCredsParams); err == nil {
		t.Fatal("expected error when creating credentials with duplicate name")
	}
}

func TestGithubCredentialsFailsToDeleteWhenInUse(t *testing.T) {
	t.Log("Testing error when deleting credentials in use")
	creds := createDummyCredentials(t, dummyCredentialsName, defaultEndpointName)

	repo := CreateRepo(t, "dummy-owner", "dummy-repo", creds.Name, "superSecret@123BlaBla")
	defer func() {
		deleteRepo(cli, authToken, repo.ID)
		deleteGithubCredentials(cli, authToken, int64(creds.ID))
	}()

	if err := deleteGithubCredentials(cli, authToken, int64(creds.ID)); err == nil {
		t.Fatal("expected error when deleting credentials in use")
	}
}

func TestGithubCredentialsFailsOnInvalidAuthType(t *testing.T) {
	t.Log("Testing error on invalid auth type")
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthType("invalid"),
		PAT: params.GithubPAT{
			OAuth2Token: "dummy",
		},
	}
	_, err := createGithubCredentials(cli, authToken, createCredsParams)
	if err == nil {
		t.Fatal("expected error when creating credentials with invalid auth type")
	}
	expectAPIStatusCode(err, 400)
}

func TestGithubCredentialsFailsWhenAuthTypeParamsAreIncorrect(t *testing.T) {
	t.Log("Testing error when auth type params are incorrect")
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		App: params.GithubApp{
			AppID:           123,
			InstallationID:  456,
			PrivateKeyBytes: getTestFileContents(t, "certs/srv-key.pem"),
		},
	}
	_, err := createGithubCredentials(cli, authToken, createCredsParams)
	if err == nil {
		t.Fatal("expected error when creating credentials with invalid auth type params")
	}
	expectAPIStatusCode(err, 400)
}

func TestGithubCredentialsFailsWhenAuthTypeParamsAreMissing(t *testing.T) {
	t.Log("Testing error when auth type params are missing")
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypeApp,
	}
	_, err := createGithubCredentials(cli, authToken, createCredsParams)
	if err == nil {
		t.Fatal("expected error when creating credentials with missing auth type params")
	}
	expectAPIStatusCode(err, 400)
}

func TestGithubCredentialsUpdateFailsWhenBothPATAndAppAreSupplied(t *testing.T) {
	t.Log("Testing error when both PAT and App are supplied")
	creds := createDummyCredentials(t, dummyCredentialsName, defaultEndpointName)
	defer DeleteGithubCredential(t, int64(creds.ID))

	updateCredsParams := params.UpdateGithubCredentialsParams{
		PAT: &params.GithubPAT{
			OAuth2Token: "dummy",
		},
		App: &params.GithubApp{
			AppID:           123,
			InstallationID:  456,
			PrivateKeyBytes: getTestFileContents(t, "certs/srv-key.pem"),
		},
	}
	_, err := updateGithubCredentials(cli, authToken, int64(creds.ID), updateCredsParams)
	if err == nil {
		t.Fatal("expected error when updating credentials with both PAT and App")
	}
	expectAPIStatusCode(err, 400)
}

func TestGithubCredentialsFailWhenAppKeyIsInvalid(t *testing.T) {
	t.Log("Testing error when app key is invalid")
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypeApp,
		App: params.GithubApp{
			AppID:           123,
			InstallationID:  456,
			PrivateKeyBytes: []byte("invalid"),
		},
	}
	_, err := createGithubCredentials(cli, authToken, createCredsParams)
	if err == nil {
		t.Fatal("expected error when creating credentials with invalid app key")
	}
	expectAPIStatusCode(err, 400)
}

func TestGithubCredentialsFailWhenEndpointDoesntExist(t *testing.T) {
	t.Log("Testing error when endpoint doesn't exist")
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    "iDontExist.example.com",
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: "dummy",
		},
	}
	_, err := createGithubCredentials(cli, authToken, createCredsParams)
	if err == nil {
		t.Fatal("expected error when creating credentials with invalid endpoint")
	}
	expectAPIStatusCode(err, 404)
}

func TestGithubCredentialsFailsOnDuplicateName(t *testing.T) {
	t.Log("Testing error on duplicate credentials name")
	creds := createDummyCredentials(t, dummyCredentialsName, defaultEndpointName)
	defer DeleteGithubCredential(t, int64(creds.ID))

	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: "dummy",
		},
	}
	_, err := createGithubCredentials(cli, authToken, createCredsParams)
	if err == nil {
		t.Fatal("expected error when creating credentials with duplicate name")
	}
	expectAPIStatusCode(err, 409)
}
