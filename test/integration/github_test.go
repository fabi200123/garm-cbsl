package integration

// import (
// 	"github.com/cloudbase/garm/params"
// 	"github.com/stretchr/testify/assert"
// )

// func (suite *GarmSuite) TestCredentials() {
// 	t := suite.T()
// 	t.Log("Create credentials")
// 	createParams := params.CreateGithubCredentialsParams{
// 		Name:        "dummy",
// 		Description: "dummy credentials",
// 		AuthType:    params.GithubAuthTypePAT,
// 		PAT: params.GithubPAT{
// 			OAuth2Token: "dummy",
// 		},
// 	}
// 	creds, err := createGithubCredentials(suite.cli, suite.authToken, createParams)
// 	assert.NoError(t, err, "error creating credentials")
// 	assert.NotNil(t, creds, "credentials is nil")

// 	t.Log("List credentials")
// 	credentials, err := listCredentials(suite.cli, suite.authToken)
// 	assert.NoError(t, err, "error listing credentials")
// 	assert.Equal(t, credentials, creds, "credentials not found in list")

// 	t.Log("Update credentials")
// 	updatedName := "updated-dummy"
// 	updatedDescription := "updated dummy credentials"
// 	updateParams := params.UpdateGithubCredentialsParams{
// 		Name:        &updatedName,
// 		Description: &updatedDescription,
// 	}
// 	updatedCreds, err := updateGithubCredentials(suite.cli, suite.authToken, int64(creds.ID), updateParams)
// 	assert.NoError(t, err, "error updating credentials")
// 	assert.Equal(t, updatedName, updatedCreds.Name, "name mismatch")

// 	t.Log("Delete credentials")
// 	err = deleteGithubCredentials(suite.cli, suite.authToken, int64(creds.ID))
// 	assert.NoError(t, err, "error deleting credentials")
// }
