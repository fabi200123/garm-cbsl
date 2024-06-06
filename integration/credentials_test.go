package integration

import (
	"log/slog"
	"testing"

	"github.com/cloudbase/garm/params"
	"github.com/stretchr/testify/assert"
)

func TestGithubCredentialsErrorOnDuplicateCredentialsName(t *testing.T) {
	slog.Info("Testing error on duplicate credentials name")
	creds, err := createDummyCredentials(dummyCredentialsName, defaultEndpointName)
	assert.NoError(t, err)
	defer DeleteGithubCredential(int64(creds.ID))

	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: "dummy",
		},
	}
	_, err = createGithubCredentials(cli, authToken, createCredsParams)
	assert.Error(t, err, "expected error when creating credentials with duplicate name")
}

func TestGithubCredentialsFailsToDeleteWhenInUse(t *testing.T) {
	slog.Info("Testing error when deleting credentials in use")
	creds, err := createDummyCredentials(dummyCredentialsName, defaultEndpointName)
	assert.NoError(t, err)

	repo, err := CreateRepo("dummy-owner", "dummy-repo", creds.Name, "superSecret@123BlaBla")
	assert.NoError(t, err)
	defer func() {
		deleteRepo(cli, authToken, repo.ID)
		deleteGithubCredentials(cli, authToken, int64(creds.ID))
	}()

	err = deleteGithubCredentials(cli, authToken, int64(creds.ID))
	assert.Error(t, err, "expected error when deleting credentials in use")
}

func TestGithubCredentialsFailsOnInvalidAuthType(t *testing.T) {
	t.Logf("Testing error on invalid auth type")
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
	assert.Error(t, err, "expected error when creating credentials with invalid auth type")
	expectAPIStatusCode(err, 400)
}

func TestGithubCredentialsFailsWhenAuthTypeParamsAreIncorrect(t *testing.T) {
	t.Logf("Testing error when auth type params are incorrect")
	privateKeyBytes, err := getTestFileContents("certs/srv-key.pem")
	assert.NoError(t, err)
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		App: params.GithubApp{
			AppID:           123,
			InstallationID:  456,
			PrivateKeyBytes: privateKeyBytes,
		},
	}
	_, err = createGithubCredentials(cli, authToken, createCredsParams)
	assert.Error(t, err, "expected error when creating credentials with invalid auth type params")

	expectAPIStatusCode(err, 400)
}

func TestGithubCredentialsFailsWhenAuthTypeParamsAreMissing(t *testing.T) {
	slog.Info("Testing error when auth type params are missing")
	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypeApp,
	}
	_, err := createGithubCredentials(cli, authToken, createCredsParams)
	assert.Error(t, err, "expected error when creating credentials with missing auth type params")
	expectAPIStatusCode(err, 400)
}

func TestGithubCredentialsUpdateFailsWhenBothPATAndAppAreSupplied(t *testing.T) {
	slog.Info("Testing error when both PAT and App are supplied")
	creds, err := createDummyCredentials(dummyCredentialsName, defaultEndpointName)
	assert.NoError(t, err)
	defer DeleteGithubCredential(int64(creds.ID))

	privateKeyBytes, err := getTestFileContents("certs/srv-key.pem")
	assert.NoError(t, err)
	updateCredsParams := params.UpdateGithubCredentialsParams{
		PAT: &params.GithubPAT{
			OAuth2Token: "dummy",
		},
		App: &params.GithubApp{
			AppID:           123,
			InstallationID:  456,
			PrivateKeyBytes: privateKeyBytes,
		},
	}
	_, err = updateGithubCredentials(cli, authToken, int64(creds.ID), updateCredsParams)
	assert.Error(t, err, "expected error when updating credentials with both PAT and App")
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
	assert.Error(t, err, "expected error when creating credentials with invalid app key")
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
	assert.Error(t, err, "expected error when creating credentials with invalid endpoint")
	expectAPIStatusCode(err, 404)
}

func TestGithubCredentialsFailsOnDuplicateName(t *testing.T) {
	slog.Info("Testing error on duplicate credentials name")
	creds, err := createDummyCredentials(dummyCredentialsName, defaultEndpointName)
	assert.NoError(t, err)
	defer DeleteGithubCredential(int64(creds.ID))

	createCredsParams := params.CreateGithubCredentialsParams{
		Name:        dummyCredentialsName,
		Endpoint:    defaultEndpointName,
		Description: "GARM test credentials",
		AuthType:    params.GithubAuthTypePAT,
		PAT: params.GithubPAT{
			OAuth2Token: "dummy",
		},
	}
	_, err = createGithubCredentials(cli, authToken, createCredsParams)
	assert.Error(t, err, "expected error when creating credentials with duplicate name")
	expectAPIStatusCode(err, 409)
}
