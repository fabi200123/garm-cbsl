package e2e

import (
	"log/slog"
	"testing"
	"time"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/params"
)

func CreateRepo(t *testing.T, orgName, repoName, credentialsName, repoWebhookSecret string) *params.Repository {
	t.Log("Create repository", "owner_name", orgName, "repo_name", repoName)
	createParams := params.CreateRepoParams{
		Owner:           orgName,
		Name:            repoName,
		CredentialsName: credentialsName,
		WebhookSecret:   repoWebhookSecret,
	}
	var err error
	repo, err := createRepo(cli, authToken, createParams)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

func UpdateRepo(t *testing.T, id, credentialsName string) *params.Repository {
	t.Log("Update repo", "repo_id", id)
	updateParams := params.UpdateEntityParams{
		CredentialsName: credentialsName,
	}
	repo, err := updateRepo(cli, authToken, id, updateParams)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

func InstallRepoWebhook(t *testing.T, id string) {
	t.Log("Install repo webhook", "repo_id", id)
	webhookParams := params.InstallWebhookParams{
		WebhookEndpointType: params.WebhookEndpointDirect,
	}
	_, err := installRepoWebhook(cli, authToken, id, webhookParams)
	if err != nil {
		slog.Error("Failed to install repo webhook", "error", err)
		t.Fatal(err)
	}
	webhookInfo, err = getRepoWebhook(cli, authToken, id)
	if err != nil {
		t.Fatal(err)
	}
}

func UninstallRepoWebhook(t *testing.T, id string) {
	t.Log("Uninstall repo webhook", "repo_id", id)
	if err := uninstallRepoWebhook(cli, authToken, id); err != nil {
		t.Fatal(err)
	}
}

func CreateRepoPool(t *testing.T, repoID string, poolParams params.CreatePoolParams) *params.Pool {
	t.Log("Create repo pool", "repo_id", repoID, "pool_params", poolParams)
	repoPool, err := createRepoPool(cli, authToken, repoID, poolParams)
	if err != nil {
		slog.Error("Failed to create repo pool", "error", err)
		t.Fatal(err)
	}
	pool = repoPool
	return repoPool
}

func GetRepoPool(t *testing.T, repoID, repoPoolID string) *params.Pool {
	t.Log("Get repo pool", "repo_id", repoID, "pool_id", repoPoolID)
	pool, err := getRepoPool(cli, authToken, repoID, repoPoolID)
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func UpdateRepoPool(t *testing.T, repoID, repoPoolID string, maxRunners, minIdleRunners uint) *params.Pool {
	t.Log("Update repo pool", "repo_id", repoID, "pool_id", repoPoolID)
	poolParams := params.UpdatePoolParams{
		MinIdleRunners: &minIdleRunners,
		MaxRunners:     &maxRunners,
	}
	pool, err := updateRepoPool(cli, authToken, repoID, repoPoolID, poolParams)
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func DeleteRepoPool(t *testing.T, repoID, repoPoolID string) {
	t.Log("Delete repo pool", "repo_id", repoID, "pool_id", repoPoolID)
	if err := deleteRepoPool(cli, authToken, repoID, repoPoolID); err != nil {
		t.Fatal(err)
	}
}

func DisableRepoPool(t *testing.T, repoID, repoPoolID string) {
	t.Log("Disable repo pool", "repo_id", repoID, "pool_id", repoPoolID)
	enabled := false
	poolParams := params.UpdatePoolParams{Enabled: &enabled}
	if _, err := updateRepoPool(cli, authToken, repoID, repoPoolID, poolParams); err != nil {
		t.Fatal(err)
	}
}

func WaitRepoRunningIdleInstances(t *testing.T, repoID string, timeout time.Duration) {
	repoPools, err := listRepoPools(cli, authToken, repoID)
	if err != nil {
		t.Fatal(err)
	}
	for _, pool := range repoPools {
		err := WaitPoolInstances(t, pool.ID, commonParams.InstanceRunning, params.RunnerIdle, timeout)
		if err != nil {
			_ = dumpRepoInstancesDetails(t, repoID)
			t.Fatal(err)
		}
	}
}

func dumpRepoInstancesDetails(t *testing.T, repoID string) error {
	// print repo details
	t.Log("Dumping repo details", "repo_id", repoID)
	repo, err := getRepo(cli, authToken, repoID)
	if err != nil {
		return err
	}
	if err := printJSONResponse(repo); err != nil {
		return err
	}

	// print repo instances details
	t.Log("Dumping repo instances details", "repo_id", repoID)
	instances, err := listRepoInstances(cli, authToken, repoID)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		instance, err := getInstance(cli, authToken, instance.Name)
		if err != nil {
			return err
		}
		t.Log("Instance info", "instance_name", instance.Name)
		if err := printJSONResponse(instance); err != nil {
			return err
		}
	}
	return nil
}
