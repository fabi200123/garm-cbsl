package integration

import (
	"log/slog"
	"net/url"
	"os"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/client"
	"github.com/cloudbase/garm/params"
	"github.com/go-openapi/runtime"
	openapiRuntimeClient "github.com/go-openapi/runtime/client"
)

var (
	adminPassword = os.Getenv("GARM_PASSWORD")
	adminUsername = os.Getenv("GARM_ADMIN_USERNAME")
	adminFullName = "GARM Admin"
	adminEmail    = "admin@example.com"

	baseURL         = os.Getenv("GARM_BASE_URL")
	credentialsName = os.Getenv("CREDENTIALS_NAME")

	repoName          = os.Getenv("REPO_NAME")
	repoWebhookSecret = os.Getenv("REPO_WEBHOOK_SECRET")
	repoPoolParams    = params.CreatePoolParams{
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
	repoPoolParams2 = params.CreatePoolParams{
		MaxRunners:     2,
		MinIdleRunners: 0,
		Flavor:         "default",
		Image:          "ubuntu:22.04",
		OSType:         commonParams.Linux,
		OSArch:         commonParams.Amd64,
		ProviderName:   "test_external",
		Tags:           []string{"repo-runner-2"},
		Enabled:        true,
	}

	orgName          = os.Getenv("ORG_NAME")
	orgWebhookSecret = os.Getenv("ORG_WEBHOOK_SECRET")
	orgPoolParams    = params.CreatePoolParams{
		MaxRunners:     2,
		MinIdleRunners: 0,
		Flavor:         "default",
		Image:          "ubuntu:22.04",
		OSType:         commonParams.Linux,
		OSArch:         commonParams.Amd64,
		ProviderName:   "lxd_local",
		Tags:           []string{"org-runner"},
		Enabled:        true,
	}

	cli              *client.GarmAPI
	authToken        runtime.ClientAuthInfoWriter
	ghToken          = os.Getenv("GH_TOKEN")
	workflowFileName = os.Getenv("WORKFLOW_FILE_NAME")
	repo             *params.Repository
	repoPool         *params.Pool
	repoPool2        *params.Pool
	orgPool          *params.Pool
	org              *params.Organization
	orgHookInfo      *params.HookInfo
)

func StartClient() {
	// Start the client
	garmURL, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Failed to get GARM_BASE_URL", err)
	}

	apiPath, err := url.JoinPath(garmURL.Path, client.DefaultBasePath)
	if err != nil {
		slog.Error("Failed to join path", err)
	}

	transportCfg := client.DefaultTransportConfig().
		WithHost(garmURL.Host).
		WithBasePath(apiPath).
		WithSchemes([]string{garmURL.Scheme})
	cli = client.NewHTTPClientWithConfig(nil, transportCfg)

	slog.Info("First run")
	newUser := params.NewUserParams{
		Username: adminUsername,
		Password: adminPassword,
		FullName: adminFullName,
		Email:    adminEmail,
	}
	_, err = firstRun(cli, newUser)
	if err != nil {
		slog.Error("Error at first run", err)
	}

	slog.Info("Login")
	loginParams := params.PasswordLoginParams{
		Username: adminUsername,
		Password: adminPassword,
	}
	token, err := login(cli, loginParams)
	if err != nil {
		slog.Error("Error at login", err)
	}
	authToken = openapiRuntimeClient.BearerToken(token)
	slog.Info("Log in succesful")

	EnsureTestCredentials(credentialsName, ghToken, "github.com")
	slog.Info("Test credentials created")

	// Get Controller Info
	_, err = getControllerInfo(cli, authToken)

	if err != nil {
		slog.Error("Error getting controller info", err)
	}
}
