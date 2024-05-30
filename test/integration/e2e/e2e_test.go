package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudbase/garm/params"
)

func ListCredentials(t *testing.T) params.Credentials {
	t.Log("List credentials")
	credentials, err := listCredentials(cli, authToken)
	if err != nil {
		t.Fatal(err)
	}
	return credentials
}

func CreateGithubCredentials(t *testing.T, credentialsParams params.CreateGithubCredentialsParams) *params.GithubCredentials {
	t.Logf("Create GitHub credentials")
	credentials, err := createGithubCredentials(cli, authToken, credentialsParams)
	if err != nil {
		t.Fatal(err)
	}
	return credentials
}

func GetGithubCredential(t *testing.T, id int64) *params.GithubCredentials {
	t.Logf("Get GitHub credential")
	credentials, err := getGithubCredential(cli, authToken, id)
	if err != nil {
		t.Fatal(err)
	}
	return credentials
}

func DeleteGithubCredential(t *testing.T, id int64) {
	t.Log("Delete GitHub credential")
	if err := deleteGithubCredentials(cli, authToken, id); err != nil {
		t.Fatal(err)
	}
}

func CreateGithubEndpoint(t *testing.T, endpointParams params.CreateGithubEndpointParams) *params.GithubEndpoint {
	t.Log("Create GitHub endpoint")
	endpoint, err := createGithubEndpoint(cli, authToken, endpointParams)
	if err != nil {
		t.Fatal(err)
	}
	return endpoint
}

func ListGithubEndpoints(t *testing.T) params.GithubEndpoints {
	t.Log("List GitHub endpoints")
	endpoints, err := listGithubEndpoints(cli, authToken)
	if err != nil {
		t.Fatal(err)
	}
	return endpoints
}

func GetGithubEndpoint(t *testing.T, name string) *params.GithubEndpoint {
	t.Log("Get GitHub endpoint")
	endpoint, err := getGithubEndpoint(cli, authToken, name)
	if err != nil {
		t.Fatal(err)
	}
	return endpoint
}

func DeleteGithubEndpoint(t *testing.T, name string) {
	t.Log("Delete GitHub endpoint")
	if err := deleteGithubEndpoint(cli, authToken, name); err != nil {
		t.Fatal(err)
	}
}

func UpdateGithubEndpoint(t *testing.T, name string, updateParams params.UpdateGithubEndpointParams) *params.GithubEndpoint {
	t.Log("Update GitHub endpoint")
	updated, err := updateGithubEndpoint(cli, authToken, name, updateParams)
	if err != nil {
		t.Fatal(err)
	}
	return updated
}

func ListProviders(t *testing.T) params.Providers {
	t.Log("List providers")
	providers, err := listProviders(cli, authToken)
	if err != nil {
		t.Fatal(err)
	}
	return providers
}

func GetMetricsToken(t *testing.T) {
	t.Log("Get metrics token")
	_, err := getMetricsToken(cli, authToken)
	if err != nil {
		t.Fatal(err)
	}
}

func GetControllerInfo(t *testing.T) *params.ControllerInfo {
	t.Log("Get controller info")
	controllerInfo, err := getControllerInfo(cli, authToken)
	if err != nil {
		t.Fatal(err)
	}
	if err := appendCtrlInfoToGitHubEnv(t, &controllerInfo); err != nil {
		t.Fatal(err)
	}
	if err := printJSONResponse(controllerInfo); err != nil {
		t.Fatal(err)
	}
	return &controllerInfo
}

func GracefulCleanup(t *testing.T) {
	t.Log("Graceful cleanup")
	// disable all the pools
	pools, err := listPools(cli, authToken)
	if err != nil {
		t.Fatal(err)
	}
	enabled := false
	poolParams := params.UpdatePoolParams{Enabled: &enabled}
	for _, pool := range pools {
		if _, err := updatePool(cli, authToken, pool.ID, poolParams); err != nil {
			t.Fatal(err)
		}
		t.Log("Pool disabled", "pool_id", pool.ID, "stage", "graceful_cleanup")
	}

	// delete all the instances
	for _, pool := range pools {
		poolInstances, err := listPoolInstances(cli, authToken, pool.ID)
		if err != nil {
			t.Fatal(err)
		}
		for _, instance := range poolInstances {
			if err := deleteInstance(cli, authToken, instance.Name, false, false); err != nil {
				t.Fatal(err)
			}
			t.Log("Instance deletion initiated", "instance", instance.Name, "stage", "graceful_cleanup")
		}
	}

	// wait for all instances to be deleted
	for _, pool := range pools {
		if err := waitPoolNoInstances(t, pool.ID, 3*time.Minute); err != nil {
			t.Fatal(err)
		}
	}

	// delete all the pools
	for _, pool := range pools {
		if err := deletePool(cli, authToken, pool.ID); err != nil {
			t.Fatal(err)
		}
		t.Log("Pool deleted", "pool_id", pool.ID, "stage", "graceful_cleanup")
	}

	// delete all the repositories
	repos, err := listRepos(cli, authToken)
	if err != nil {
		t.Fatal(err)
	}
	for _, repo := range repos {
		if err := deleteRepo(cli, authToken, repo.ID); err != nil {
			t.Fatal(err)
		}
		t.Log("Repo deleted", "repo_id", repo.ID, "stage", "graceful_cleanup")
	}

	// delete all the organizations
	orgs, err := listOrgs(cli, authToken)
	if err != nil {
		t.Fatal(err)
	}
	for _, org := range orgs {
		if err := deleteOrg(cli, authToken, org.ID); err != nil {
			t.Fatal(err)
		}
		t.Log("Org deleted", "org_id", org.ID, "stage", "graceful_cleanup")
	}
}

func appendCtrlInfoToGitHubEnv(t *testing.T, controllerInfo *params.ControllerInfo) error {
	envFile, found := os.LookupEnv("GITHUB_ENV")
	if !found {
		t.Log("GITHUB_ENV not set, skipping appending controller info")
		return nil
	}
	file, err := os.OpenFile(envFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(fmt.Sprintf("export GARM_CONTROLLER_ID=%s\n", controllerInfo.ControllerID)); err != nil {
		return err
	}
	return nil
}
