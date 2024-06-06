package integration

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/cloudbase/garm/params"
)

func ListCredentials() (params.Credentials, error) {
	slog.Info("List credentials")
	credentials, err := listCredentials(cli, authToken)
	if err != nil {
		return params.Credentials{}, err
	}
	return credentials, nil
}

func GetGithubCredential(id int64) (*params.GithubCredentials, error) {
	slog.Info("Get GitHub credential")
	credentials, err := getGithubCredential(cli, authToken, id)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func ListGithubEndpoints() (params.GithubEndpoints, error) {
	slog.Info("List GitHub endpoints")
	endpoints, err := listGithubEndpoints(cli, authToken)
	if err != nil {
		return params.GithubEndpoints{}, err
	}
	return endpoints, nil
}

func GetGithubEndpoint(name string) (*params.GithubEndpoint, error) {
	slog.Info("Get GitHub endpoint")
	endpoint, err := getGithubEndpoint(cli, authToken, name)
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func DeleteGithubEndpoint(name string) error {
	slog.Info("Delete GitHub endpoint")
	if err := deleteGithubEndpoint(cli, authToken, name); err != nil {
		return err
	}

	return nil
}

func UpdateGithubEndpoint(name string, updateParams params.UpdateGithubEndpointParams) (*params.GithubEndpoint, error) {
	slog.Info("Update GitHub endpoint")
	updated, err := updateGithubEndpoint(cli, authToken, name, updateParams)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func ListProviders() (params.Providers, error) {
	slog.Info("List providers")
	providers, err := listProviders(cli, authToken)
	if err != nil {
		return params.Providers{}, err
	}
	return providers, nil
}

func GetMetricsToken() error {
	slog.Info("Get metrics token")
	_, err := getMetricsToken(cli, authToken)
	if err != nil {
		return err
	}
	return nil
}

func GetControllerInfo() (*params.ControllerInfo, error) {
	slog.Info("Get controller info")
	controllerInfo, err := getControllerInfo(cli, authToken)
	if err != nil {
		return nil, err
	}
	if err := appendCtrlInfoToGitHubEnv(&controllerInfo); err != nil {
		return nil, err
	}
	if err := printJSONResponse(controllerInfo); err != nil {
		return nil, err
	}
	return &controllerInfo, nil
}

func appendCtrlInfoToGitHubEnv(controllerInfo *params.ControllerInfo) error {
	envFile, found := os.LookupEnv("GITHUB_ENV")
	if !found {
		slog.Info("GITHUB_ENV not set, skipping appending controller info")
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
