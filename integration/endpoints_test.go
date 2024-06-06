package integration

import (
	"log/slog"
	"testing"

	"github.com/cloudbase/garm/params"
	"github.com/stretchr/testify/assert"
)

func TestGithubEndpointOperations(t *testing.T) {
	slog.Info("Testing endpoint operations")
	MustDefaultGithubEndpoint()

	caBundle, err := getTestFileContents("certs/srv-pub.pem")
	assert.NoError(t, err)

	endpointParams := params.CreateGithubEndpointParams{
		Name:          "test-endpoint",
		Description:   "Test endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
		CACertBundle:  caBundle,
	}

	endpoint, err := CreateGithubEndpoint(endpointParams)
	assert.NoError(t, err)
	assert.Equal(t, endpoint.Name, endpointParams.Name, "Endpoint name mismatch")
	assert.Equal(t, endpoint.Description, endpointParams.Description, "Endpoint description mismatch")
	assert.Equal(t, endpoint.BaseURL, endpointParams.BaseURL, "Endpoint base URL mismatch")
	assert.Equal(t, endpoint.APIBaseURL, endpointParams.APIBaseURL, "Endpoint API base URL mismatch")
	assert.Equal(t, endpoint.UploadBaseURL, endpointParams.UploadBaseURL, "Endpoint upload base URL mismatch")
	assert.Equal(t, string(endpoint.CACertBundle), string(caBundle), "Endpoint CA cert bundle mismatch")

	endpoint2, err := GetGithubEndpoint(endpointParams.Name)
	assert.NoError(t, err)
	assert.NotNil(t, endpoint, "endpoint is nil")
	assert.NotNil(t, endpoint2, "endpoint2 is nil")

	err = checkEndpointParamsAreEqual(*endpoint, *endpoint2)
	assert.NoError(t, err, "endpoint params are not equal")
	endpoints, err := ListGithubEndpoints()
	assert.NoError(t, err, "error listing github endpoints")
	var found bool
	for _, ep := range endpoints {
		if ep.Name == endpointParams.Name {
			checkEndpointParamsAreEqual(*endpoint, ep)
			found = true
			break
		}
	}
	assert.Equal(t, found, true, "endpoint not found in list")

	err = DeleteGithubEndpoint(endpoint.Name)
	assert.NoError(t, err, "error deleting github endpoint")
}

func TestGithubEndpointMustFailToDeleteDefaultGithubEndpoint(t *testing.T) {
	t.Log("Testing error when deleting default github.com endpoint")
	err := deleteGithubEndpoint(cli, authToken, "github.com")
	assert.Error(t, err, "expected error when attempting to delete the default github.com endpoint")
}

func TestGithubEndpointFailsOnInvalidCABundle(t *testing.T) {
	t.Log("Testing endpoint creation with invalid CA cert bundle")
	badCABundle, err := getTestFileContents("certs/srv-key.pem")
	assert.NoError(t, err, "error reading CA cert bundle")

	endpointParams := params.CreateGithubEndpointParams{
		Name:          "dummy",
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
		CACertBundle:  badCABundle,
	}

	_, err = createGithubEndpoint(cli, authToken, endpointParams)
	assert.Error(t, err, "expected error when creating endpoint with invalid CA cert bundle")
}

func TestGithubEndpointDeletionFailsWhenCredentialsExist(t *testing.T) {
	slog.Info("Testing endpoint deletion when credentials exist")
	endpointParams := params.CreateGithubEndpointParams{
		Name:          "dummy",
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
	}

	endpoint, err := CreateGithubEndpoint(endpointParams)
	assert.NoError(t, err, "error creating github endpoint")
	creds, err := createDummyCredentials("test-creds", endpoint.Name)
	assert.NoError(t, err, "error creating dummy credentials")

	err = deleteGithubEndpoint(cli, authToken, endpoint.Name)
	assert.Error(t, err, "expected error when deleting endpoint with credentials")

	err = DeleteGithubCredential(int64(creds.ID))
	assert.NoError(t, err, "error deleting credentials")
	err = DeleteGithubEndpoint(endpoint.Name)
	assert.NoError(t, err, "error deleting endpoint")
}

func TestGithubEndpointFailsOnDuplicateName(t *testing.T) {
	t.Log("Testing endpoint creation with duplicate name")
	endpointParams := params.CreateGithubEndpointParams{
		Name:          "github.com",
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
	}

	_, err := createGithubEndpoint(cli, authToken, endpointParams)
	assert.Error(t, err, "expected error when creating endpoint with duplicate name")
}

func TestGithubEndpointUpdateEndpoint(t *testing.T) {
	slog.Info("Testing endpoint update")
	endpoint, err := createDummyEndpoint("dummy")
	assert.NoError(t, err, "error creating dummy endpoint")
	defer DeleteGithubEndpoint(endpoint.Name)

	newDescription := "Updated description"
	newBaseURL := "https://ghes2.example.com"
	newAPIBaseURL := "https://api.ghes2.example.com/"
	newUploadBaseURL := "https://uploads.ghes2.example.com/"
	newCABundle, err := getTestFileContents("certs/srv-pub.pem")
	assert.NoError(t, err, "error reading CA cert bundle")

	updateParams := params.UpdateGithubEndpointParams{
		Description:   &newDescription,
		BaseURL:       &newBaseURL,
		APIBaseURL:    &newAPIBaseURL,
		UploadBaseURL: &newUploadBaseURL,
		CACertBundle:  newCABundle,
	}

	updated, err := updateGithubEndpoint(cli, authToken, endpoint.Name, updateParams)
	assert.NoError(t, err, "error updating github endpoint")

	assert.Equal(t, updated.Name, endpoint.Name, "Endpoint name mismatch")
	assert.Equal(t, updated.Description, newDescription, "Endpoint description mismatch")
	assert.Equal(t, updated.BaseURL, newBaseURL, "Endpoint base URL mismatch")
	assert.Equal(t, updated.APIBaseURL, newAPIBaseURL, "Endpoint API base URL mismatch")
	assert.Equal(t, updated.UploadBaseURL, newUploadBaseURL, "Endpoint upload base URL mismatch")
	assert.Equal(t, string(updated.CACertBundle), string(newCABundle), "Endpoint CA cert bundle mismatch")
}
