package integration

import (
	"context"
	"fmt"
	"time"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/params"
	"github.com/google/go-github/v57/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func (suite *GarmSuite) EnsureTestCredentials(name string, oauthToken string, endpointName string) {
	t := suite.T()
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
	suite.CreateGithubCredentials(createCredsParams)

	createCredsParams.Name = fmt.Sprintf("%s-clone", name)
	suite.CreateGithubCredentials(createCredsParams)
}

func (suite *GarmSuite) TestRepositories() {
	t := suite.T()

	t.Logf("Update repo with repo_id %s", suite.repo.ID)
	updateParams := params.UpdateEntityParams{
		CredentialsName: fmt.Sprintf("%s-clone", suite.credentialsName),
	}
	repo, err := updateRepo(suite.cli, suite.authToken, suite.repo.ID, updateParams)
	assert.NoError(t, err, "error updating repository")
	assert.Equal(t, fmt.Sprintf("%s-clone", suite.credentialsName), repo.CredentialsName, "credentials name mismatch")
	suite.repo = repo

	hookRepoInfo := suite.InstallRepoWebhook(suite.repo.ID)
	suite.ValidateRepoWebhookInstalled(suite.ghToken, hookRepoInfo.URL, orgName, repoName)
	suite.UninstallRepoWebhook(suite.repo.ID)
	suite.ValidateRepoWebhookUninstalled(suite.ghToken, hookRepoInfo.URL, orgName, repoName)

	suite.InstallRepoWebhook(suite.repo.ID)
	suite.ValidateRepoWebhookInstalled(suite.ghToken, hookRepoInfo.URL, orgName, repoName)

	repoPoolParams := params.CreatePoolParams{
		MaxRunners:     2,
		MinIdleRunners: 0,
		Flavor:         "default",
		Image:          "ubuntu:22.04",
		OSType:         commonParams.Linux,
		OSArch:         commonParams.Amd64,
		ProviderName:   "lxd_local",
		Tags:           []string{"repo-runner"},
		Enabled:        true,
	}

	repoPool := suite.CreateRepoPool(suite.repo.ID, repoPoolParams)
	assert.Equal(t, repoPool.MaxRunners, repoPoolParams.MaxRunners, "max runners mismatch")
	assert.Equal(t, repoPool.MinIdleRunners, repoPoolParams.MinIdleRunners, "min idle runners mismatch")

	repoPoolGet := suite.GetRepoPool(suite.repo.ID, repoPool.ID)
	assert.Equal(t, *repoPool, *repoPoolGet, "pool get mismatch")

	suite.DeleteRepoPool(suite.repo.ID, repoPool.ID)

	repoPool = suite.CreateRepoPool(suite.repo.ID, repoPoolParams)
	updatedRepoPool := suite.UpdateRepoPool(suite.repo.ID, repoPool.ID, repoPoolParams.MaxRunners, 1)
	assert.NotEqual(t, updatedRepoPool.MinIdleRunners, repoPool.MinIdleRunners, "min idle runners mismatch")

	suite.WaitRepoRunningIdleInstances(suite.repo.ID, 6*time.Minute)
}

func (suite *GarmSuite) InstallRepoWebhook(id string) *params.HookInfo {
	t := suite.T()
	t.Logf("Install repo webhook with repo_id %s", id)
	webhookParams := params.InstallWebhookParams{
		WebhookEndpointType: params.WebhookEndpointDirect,
	}
	_, err := installRepoWebhook(suite.cli, suite.authToken, id, webhookParams)
	assert.NoError(t, err, "error installing repository webhook")

	webhookInfo, err := getRepoWebhook(suite.cli, suite.authToken, id)
	assert.NoError(t, err, "error getting repository webhook")
	return webhookInfo
}

func (suite *GarmSuite) ValidateRepoWebhookInstalled(ghToken, url, orgName, repoName string) {
	t := suite.T()
	hook, err := getGhRepoWebhook(url, ghToken, orgName, repoName)
	assert.NoError(t, err, "error getting github webhook")
	assert.NotNil(t, hook, "github webhook with url %s, for repo %s/%s was not properly installed", url, orgName, repoName)
}

func getGhRepoWebhook(url, ghToken, orgName, repoName string) (*github.Hook, error) {
	client := getGithubClient(ghToken)
	ghRepoHooks, _, err := client.Repositories.ListHooks(context.Background(), orgName, repoName, nil)
	if err != nil {
		return nil, err
	}

	for _, hook := range ghRepoHooks {
		hookURL, ok := hook.Config["url"].(string)
		if ok && hookURL == url {
			return hook, nil
		}
	}

	return nil, nil
}

func getGithubClient(oauthToken string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: oauthToken})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func (suite *GarmSuite) UninstallRepoWebhook(id string) {
	t := suite.T()
	t.Logf("Uninstall repo webhook with repo_id %s", id)
	err := uninstallRepoWebhook(suite.cli, suite.authToken, id)
	assert.NoError(t, err, "error uninstalling repository webhook")
}

func (suite *GarmSuite) ValidateRepoWebhookUninstalled(ghToken, url, orgName, repoName string) {
	t := suite.T()
	hook, err := getGhRepoWebhook(url, ghToken, orgName, repoName)
	assert.NoError(t, err, "error getting github webhook")
	assert.Nil(t, hook, "github webhook with url %s, for repo %s/%s was not properly uninstalled", url, orgName, repoName)
}

func (suite *GarmSuite) CreateRepoPool(repoID string, poolParams params.CreatePoolParams) *params.Pool {
	t := suite.T()
	t.Logf("Create repo pool with repo_id %s and pool_params %+v", repoID, poolParams)
	pool, err := createRepoPool(suite.cli, suite.authToken, repoID, poolParams)
	assert.NoError(t, err, "error creating repository pool")
	return pool
}

func (suite *GarmSuite) GetRepoPool(repoID, repoPoolID string) *params.Pool {
	t := suite.T()
	t.Logf("Get repo pool repo_id %s and pool_id %s", repoID, repoPoolID)
	pool, err := getRepoPool(suite.cli, suite.authToken, repoID, repoPoolID)
	assert.NoError(t, err, "error getting repository pool")
	return pool
}

func (suite *GarmSuite) DeleteRepoPool(repoID, repoPoolID string) {
	t := suite.T()
	t.Logf("Delete repo pool with repo_id %s and pool_id %s", repoID, repoPoolID)
	err := deleteRepoPool(suite.cli, suite.authToken, repoID, repoPoolID)
	assert.NoError(t, err, "error deleting repository pool")
}

func (suite *GarmSuite) UpdateRepoPool(repoID, repoPoolID string, maxRunners, minIdleRunners uint) *params.Pool {
	t := suite.T()
	t.Logf("Update repo pool with repo_id %s and pool_id %s", repoID, repoPoolID)
	poolParams := params.UpdatePoolParams{
		MinIdleRunners: &minIdleRunners,
		MaxRunners:     &maxRunners,
	}
	pool, err := updateRepoPool(suite.cli, suite.authToken, repoID, repoPoolID, poolParams)
	assert.NoError(t, err, "error updating repository pool")
	return pool
}

func (suite *GarmSuite) WaitRepoRunningIdleInstances(repoID string, timeout time.Duration) {
	t := suite.T()
	repoPools, err := listRepoPools(suite.cli, suite.authToken, repoID)
	assert.NoError(t, err, "error listing repo pools")
	for _, pool := range repoPools {
		err := suite.WaitPoolInstances(pool.ID, commonParams.InstanceRunning, params.RunnerIdle, timeout)
		if err != nil {
			suite.dumpRepoInstancesDetails(repoID)
			t.Errorf("error waiting for pool instances to be running idle: %v", err)
		}
	}
}

func (suite *GarmSuite) dumpRepoInstancesDetails(repoID string) {
	t := suite.T()
	// print repo details
	t.Logf("Dumping repo details for repo %s", repoID)
	repo, err := getRepo(suite.cli, suite.authToken, repoID)
	assert.NoError(t, err, "error getting repo")
	err = printJSONResponse(repo)
	assert.NoError(t, err, "error printing repo")

	// print repo instances details
	t.Logf("Dumping repo instances details for repo %s", repoID)
	instances, err := listRepoInstances(suite.cli, suite.authToken, repoID)
	assert.NoError(t, err, "error listing repo instances")
	for _, instance := range instances {
		instance, err := getInstance(suite.cli, suite.authToken, instance.Name)
		assert.NoError(t, err, "error getting instance")
		t.Logf("Instance info for instance %s", instance.Name)
		err = printJSONResponse(instance)
		assert.NoError(t, err, "error printing instance")
	}
}
