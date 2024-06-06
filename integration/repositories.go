package integration

import (
	"log/slog"
	"time"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/params"
)

func CreateRepo(orgName, repoName, credentialsName, repoWebhookSecret string) (*params.Repository, error) {
	slog.Info("Create repository", "owner_name", orgName, "repo_name", repoName)
	createParams := params.CreateRepoParams{
		Owner:           orgName,
		Name:            repoName,
		CredentialsName: credentialsName,
		WebhookSecret:   repoWebhookSecret,
	}
	repo, err := createRepo(cli, authToken, createParams)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func UpdateRepo(id, credentialsName string) (*params.Repository, error) {
	slog.Info("Update repo", "repo_id", id)
	updateParams := params.UpdateEntityParams{
		CredentialsName: credentialsName,
	}
	repo, err := updateRepo(cli, authToken, id, updateParams)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func InstallRepoWebhook(id string) (*params.HookInfo, error) {
	slog.Info("Install repo webhook", "repo_id", id)
	webhookParams := params.InstallWebhookParams{
		WebhookEndpointType: params.WebhookEndpointDirect,
	}
	_, err := installRepoWebhook(cli, authToken, id, webhookParams)
	if err != nil {
		slog.Error("Failed to install repo webhook", "error", err)
		return nil, err
	}
	webhookInfo, err := getRepoWebhook(cli, authToken, id)
	if err != nil {
		return nil, err
	}
	return webhookInfo, nil
}

func UninstallRepoWebhook(id string) error {
	slog.Info("Uninstall repo webhook", "repo_id", id)
	if err := uninstallRepoWebhook(cli, authToken, id); err != nil {
		return err
	}
	return nil
}

func CreateRepoPool(repoID string, poolParams params.CreatePoolParams) (*params.Pool, error) {
	slog.Info("Create repo pool", "repo_id", repoID, "pool_params", poolParams)
	pool, err := createRepoPool(cli, authToken, repoID, poolParams)
	if err != nil {
		slog.Error("Failed to create repo pool", "error", err)
		return nil, err
	}
	return pool, nil
}

func GetRepoPool(repoID, repoPoolID string) (*params.Pool, error) {
	slog.Info("Get repo pool", "repo_id", repoID, "pool_id", repoPoolID)
	pool, err := getRepoPool(cli, authToken, repoID, repoPoolID)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func UpdateRepoPool(repoID, repoPoolID string, maxRunners, minIdleRunners uint) (*params.Pool, error) {
	slog.Info("Update repo pool", "repo_id", repoID, "pool_id", repoPoolID)
	poolParams := params.UpdatePoolParams{
		MinIdleRunners: &minIdleRunners,
		MaxRunners:     &maxRunners,
	}
	pool, err := updateRepoPool(cli, authToken, repoID, repoPoolID, poolParams)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func DeleteRepoPool(repoID, repoPoolID string) error {
	slog.Info("Delete repo pool", "repo_id", repoID, "pool_id", repoPoolID)
	if err := deleteRepoPool(cli, authToken, repoID, repoPoolID); err != nil {
		return err
	}
	return nil
}

func DisableRepoPool(repoID, repoPoolID string) error {
	slog.Info("Disable repo pool", "repo_id", repoID, "pool_id", repoPoolID)
	enabled := false
	poolParams := params.UpdatePoolParams{Enabled: &enabled}
	if _, err := updateRepoPool(cli, authToken, repoID, repoPoolID, poolParams); err != nil {
		return err
	}
	return nil
}

func WaitRepoRunningIdleInstances(repoID string, timeout time.Duration) error {
	repoPools, err := listRepoPools(cli, authToken, repoID)
	if err != nil {
		return err
	}
	for _, pool := range repoPools {
		err := WaitPoolInstances(pool.ID, commonParams.InstanceRunning, params.RunnerIdle, timeout)
		if err != nil {
			_ = dumpRepoInstancesDetails(repoID)
			return err
		}
	}
	return nil
}

func dumpRepoInstancesDetails(repoID string) error {
	// print repo details
	slog.Info("Dumping repo details", "repo_id", repoID)
	repo, err := getRepo(cli, authToken, repoID)
	if err != nil {
		return err
	}
	if err := printJSONResponse(repo); err != nil {
		return err
	}

	// print repo instances details
	slog.Info("Dumping repo instances details", "repo_id", repoID)
	instances, err := listRepoInstances(cli, authToken, repoID)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		instance, err := getInstance(cli, authToken, instance.Name)
		if err != nil {
			return err
		}
		slog.Info("Instance info", "instance_name", instance.Name)
		if err := printJSONResponse(instance); err != nil {
			return err
		}
	}
	return nil
}
