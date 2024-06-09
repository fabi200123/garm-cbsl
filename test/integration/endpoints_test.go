package integration

import (
	"github.com/cloudbase/garm/params"
	"github.com/stretchr/testify/assert"
)

func (suite *GarmSuite) TestGithubEndpointOperations() {
	t := suite.T()
	t.Log("Testing endpoint operations")
	suite.MustDefaultGithubEndpoint()

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

	endpoint, err := suite.CreateGithubEndpoint(endpointParams)
	assert.NoError(t, err)
	assert.Equal(t, endpoint.Name, endpointParams.Name, "Endpoint name mismatch")
	assert.Equal(t, endpoint.Description, endpointParams.Description, "Endpoint description mismatch")
	assert.Equal(t, endpoint.BaseURL, endpointParams.BaseURL, "Endpoint base URL mismatch")
	assert.Equal(t, endpoint.APIBaseURL, endpointParams.APIBaseURL, "Endpoint API base URL mismatch")
	assert.Equal(t, endpoint.UploadBaseURL, endpointParams.UploadBaseURL, "Endpoint upload base URL mismatch")
	assert.Equal(t, string(endpoint.CACertBundle), string(caBundle), "Endpoint CA cert bundle mismatch")

	endpoint2 := suite.GetGithubEndpoint(endpointParams.Name)
	assert.NotNil(t, endpoint, "endpoint is nil")
	assert.NotNil(t, endpoint2, "endpoint2 is nil")

	err = checkEndpointParamsAreEqual(*endpoint, *endpoint2)
	assert.NoError(t, err, "endpoint params are not equal")
	endpoints := suite.ListGithubEndpoints()
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

	err = suite.DeleteGithubEndpoint(endpoint.Name)
	assert.NoError(t, err, "error deleting github endpoint")
}

func (suite *GarmSuite) TestGithubEndpointMustFailToDeleteDefaultGithubEndpoint() {
	t := suite.T()
	t.Log("Testing error when deleting default github.com endpoint")
	err := deleteGithubEndpoint(suite.cli, suite.authToken, "github.com")
	assert.Error(t, err, "expected error when attempting to delete the default github.com endpoint")
}

func (suite *GarmSuite) TestGithubEndpointFailsOnInvalidCABundle() {
	t := suite.T()
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

	_, err = createGithubEndpoint(suite.cli, suite.authToken, endpointParams)
	assert.Error(t, err, "expected error when creating endpoint with invalid CA cert bundle")
}

func (suite *GarmSuite) TestGithubEndpointDeletionFailsWhenCredentialsExist() {
	t := suite.T()
	t.Log("Testing endpoint deletion when credentials exist")
	endpointParams := params.CreateGithubEndpointParams{
		Name:          "dummy",
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
	}

	endpoint, err := suite.CreateGithubEndpoint(endpointParams)
	assert.NoError(t, err, "error creating github endpoint")
	creds, err := suite.createDummyCredentials("test-creds", endpoint.Name)
	assert.NoError(t, err, "error creating dummy credentials")

	err = deleteGithubEndpoint(suite.cli, suite.authToken, endpoint.Name)
	assert.Error(t, err, "expected error when deleting endpoint with credentials")

	err = suite.DeleteGithubCredential(int64(creds.ID))
	assert.NoError(t, err, "error deleting credentials")
	err = suite.DeleteGithubEndpoint(endpoint.Name)
	assert.NoError(t, err, "error deleting endpoint")
}

func (suite *GarmSuite) TestGithubEndpointFailsOnDuplicateName() {
	t := suite.T()
	t.Log("Testing endpoint creation with duplicate name")
	endpointParams := params.CreateGithubEndpointParams{
		Name:          "github.com",
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
	}

	_, err := createGithubEndpoint(suite.cli, suite.authToken, endpointParams)
	assert.Error(t, err, "expected error when creating endpoint with duplicate name")
}

func (suite *GarmSuite) TestGithubEndpointUpdateEndpoint() {
	t := suite.T()
	t.Log("Testing endpoint update")
	endpoint, err := suite.createDummyEndpoint("dummy")
	assert.NoError(t, err, "error creating dummy endpoint")
	defer suite.DeleteGithubEndpoint(endpoint.Name)

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

	updated, err := updateGithubEndpoint(suite.cli, suite.authToken, endpoint.Name, updateParams)
	assert.NoError(t, err, "error updating github endpoint")

	assert.Equal(t, updated.Name, endpoint.Name, "Endpoint name mismatch")
	assert.Equal(t, updated.Description, newDescription, "Endpoint description mismatch")
	assert.Equal(t, updated.BaseURL, newBaseURL, "Endpoint base URL mismatch")
	assert.Equal(t, updated.APIBaseURL, newAPIBaseURL, "Endpoint API base URL mismatch")
	assert.Equal(t, updated.UploadBaseURL, newUploadBaseURL, "Endpoint upload base URL mismatch")
	assert.Equal(t, string(updated.CACertBundle), string(newCABundle), "Endpoint CA cert bundle mismatch")
}

func (suite *GarmSuite) MustDefaultGithubEndpoint() {
	t := suite.T()
	ep := suite.GetGithubEndpoint("github.com")

	assert.NotNil(t, ep, "default GitHub endpoint not found")
	assert.Equal(t, ep.Name, "github.com", "default GitHub endpoint name mismatch")
}

func (suite *GarmSuite) GetGithubEndpoint(name string) *params.GithubEndpoint {
	t := suite.T()
	t.Log("Get GitHub endpoint")
	endpoint, err := getGithubEndpoint(suite.cli, suite.authToken, name)
	assert.NoError(t, err, "error getting GitHub endpoint")

	return endpoint
}

func (suite *GarmSuite) CreateGithubEndpoint(params params.CreateGithubEndpointParams) (*params.GithubEndpoint, error) {
	t := suite.T()
	t.Log("Create GitHub endpoint")
	endpoint, err := createGithubEndpoint(suite.cli, suite.authToken, params)
	assert.NoError(t, err, "error creating GitHub endpoint")

	return endpoint, nil
}

func (suite *GarmSuite) DeleteGithubEndpoint(name string) error {
	t := suite.T()
	t.Log("Delete GitHub endpoint")
	err := deleteGithubEndpoint(suite.cli, suite.authToken, name)
	assert.NoError(t, err, "error deleting GitHub endpoint")

	return nil
}

func (suite *GarmSuite) ListGithubEndpoints() params.GithubEndpoints {
	t := suite.T()
	t.Log("List GitHub endpoints")
	endpoints, err := listGithubEndpoints(suite.cli, suite.authToken)
	assert.NoError(t, err, "error listing GitHub endpoints")

	return endpoints
}

func (suite *GarmSuite) createDummyEndpoint(name string) (*params.GithubEndpoint, error) {
	endpointParams := params.CreateGithubEndpointParams{
		Name:          name,
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
	}

	return suite.CreateGithubEndpoint(endpointParams)
}
