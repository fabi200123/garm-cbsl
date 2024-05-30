package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudbase/garm/params"
)

const (
	defaultEndpointName  string = "github.com"
	dummyCredentialsName string = "dummy"
)

func checkEndpointParamsAreEqual(t *testing.T, a, b params.GithubEndpoint) error {
	if a.Name != b.Name {
		t.Fatal("Endpoint name mismatch")
	}

	if a.Description != b.Description {
		t.Fatal("Endpoint description mismatch")
	}

	if a.BaseURL != b.BaseURL {
		t.Fatal("Endpoint base URL mismatch")
	}

	if a.APIBaseURL != b.APIBaseURL {
		t.Fatal("Endpoint API base URL mismatch")
	}

	if a.UploadBaseURL != b.UploadBaseURL {
		t.Fatal("Endpoint upload base URL mismatch")
	}

	if string(a.CACertBundle) != string(b.CACertBundle) {
		t.Fatal("Endpoint CA cert bundle mismatch")
	}

	return nil
}

func MustDefaultGithubEndpoint(t *testing.T) {
	ep := GetGithubEndpoint(t, "github.com")
	if ep == nil {
		t.Fatal("Default GitHub endpoint not found")
	}

	if ep.Name != "github.com" {
		t.Fatal("Default GitHub endpoint name mismatch")
	}
}

func getTestFileContents(t *testing.T, relPath string) []byte {
	baseDir := os.Getenv("GARM_CHECKOUT_DIR")
	if baseDir == "" {
		t.Fatal("GARM_CHECKOUT_DIR not set")
	}
	contents, err := os.ReadFile(filepath.Join(baseDir, "testdata", relPath))
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func TestGithubEndpointOperations(t *testing.T) {
	t.Log("Testing endpoint operations")
	MustDefaultGithubEndpoint(t)

	caBundle := getTestFileContents(t, "certs/srv-pub.pem")

	endpointParams := params.CreateGithubEndpointParams{
		Name:          "test-endpoint",
		Description:   "Test endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
		CACertBundle:  caBundle,
	}

	endpoint := CreateGithubEndpoint(t, endpointParams)
	if endpoint.Name != endpointParams.Name {
		t.Fatal("Endpoint name mismatch")
	}

	if endpoint.Description != endpointParams.Description {
		t.Fatal("Endpoint description mismatch")
	}

	if endpoint.BaseURL != endpointParams.BaseURL {
		t.Fatal("Endpoint base URL mismatch")
	}

	if endpoint.APIBaseURL != endpointParams.APIBaseURL {
		t.Fatal("Endpoint API base URL mismatch")
	}

	if endpoint.UploadBaseURL != endpointParams.UploadBaseURL {
		t.Fatal("Endpoint upload base URL mismatch")
	}

	if string(endpoint.CACertBundle) != string(caBundle) {
		t.Fatal("Endpoint CA cert bundle mismatch")
	}

	endpoint2 := GetGithubEndpoint(t, endpointParams.Name)
	if endpoint == nil || endpoint2 == nil {
		t.Fatal("endpoint is nil")
	}
	checkEndpointParamsAreEqual(t, *endpoint, *endpoint2)

	endpoints := ListGithubEndpoints(t)
	var found bool
	for _, ep := range endpoints {
		if ep.Name == endpointParams.Name {
			checkEndpointParamsAreEqual(t, *endpoint, ep)
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Endpoint not found in list")
	}

	DeleteGithubEndpoint(t, endpoint.Name)
}

func TestGithubEndpointMustFailToDeleteDefaultGithubEndpoint(t *testing.T) {
	t.Log("Testing error when deleting default github.com endpoint")
	if err := deleteGithubEndpoint(cli, authToken, "github.com"); err == nil {
		t.Fatal("expected error when attempting to delete the default github.com endpoint")
	}
}

func TestGithubEndpointFailsOnInvalidCABundle(t *testing.T) {
	t.Log("Testing endpoint creation with invalid CA cert bundle")
	badCABundle := getTestFileContents(t, "certs/srv-key.pem")

	endpointParams := params.CreateGithubEndpointParams{
		Name:          "dummy",
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
		CACertBundle:  badCABundle,
	}

	if _, err := createGithubEndpoint(cli, authToken, endpointParams); err == nil {
		t.Fatal("expected error when creating endpoint with invalid CA cert bundle")
	}
}

func TestGithubEndpointDeletionFailsWhenCredentialsExist(t *testing.T) {
	t.Log("Testing endpoint deletion when credentials exist")
	endpointParams := params.CreateGithubEndpointParams{
		Name:          "dummy",
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
	}

	endpoint := CreateGithubEndpoint(t, endpointParams)
	creds := createDummyCredentials(t, "test-creds", endpoint.Name)

	if err := deleteGithubEndpoint(cli, authToken, endpoint.Name); err == nil {
		t.Fatal("expected error when deleting endpoint with credentials")
	}

	DeleteGithubCredential(t, int64(creds.ID))
	DeleteGithubEndpoint(t, endpoint.Name)
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

	if _, err := createGithubEndpoint(cli, authToken, endpointParams); err == nil {
		t.Fatal("expected error when creating endpoint with duplicate name")
	}
}

func TestGithubEndpointUpdateEndpoint(t *testing.T) {
	t.Log("Testing endpoint update")
	endpoint := createDummyEndpoint(t, "dummy")
	defer DeleteGithubEndpoint(t, endpoint.Name)

	newDescription := "Updated description"
	newBaseURL := "https://ghes2.example.com"
	newAPIBaseURL := "https://api.ghes2.example.com/"
	newUploadBaseURL := "https://uploads.ghes2.example.com/"
	newCABundle := getTestFileContents(t, "certs/srv-pub.pem")

	updateParams := params.UpdateGithubEndpointParams{
		Description:   &newDescription,
		BaseURL:       &newBaseURL,
		APIBaseURL:    &newAPIBaseURL,
		UploadBaseURL: &newUploadBaseURL,
		CACertBundle:  newCABundle,
	}

	updated, err := updateGithubEndpoint(cli, authToken, endpoint.Name, updateParams)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Name != endpoint.Name {
		t.Fatal("Endpoint name mismatch")
	}

	if updated.Description != newDescription {
		t.Fatal("Endpoint description mismatch")
	}

	if updated.BaseURL != newBaseURL {
		t.Fatal("Endpoint base URL mismatch")
	}

	if updated.APIBaseURL != newAPIBaseURL {
		t.Fatal("Endpoint API base URL mismatch")
	}

	if updated.UploadBaseURL != newUploadBaseURL {
		t.Fatal("Endpoint upload base URL mismatch")
	}

	if string(updated.CACertBundle) != string(newCABundle) {
		t.Fatal("Endpoint CA cert bundle mismatch")
	}
}

func createDummyEndpoint(t *testing.T, name string) *params.GithubEndpoint {
	endpointParams := params.CreateGithubEndpointParams{
		Name:          name,
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
	}

	return CreateGithubEndpoint(t, endpointParams)
}
