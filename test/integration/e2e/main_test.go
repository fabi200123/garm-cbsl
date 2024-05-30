package e2e

import (
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"testing"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/client"
	"github.com/cloudbase/garm/params"
	"github.com/go-openapi/runtime"
	openapiRuntimeClient "github.com/go-openapi/runtime/client"
)

var (
	baseURL       string
	adminPassword string
	adminUsername string
	adminFullName string
	adminEmail    string

	credentialsName   string
	repoName          string
	repoWebhookSecret string
	repoPoolParams    params.CreatePoolParams
	repoPoolParams2   params.CreatePoolParams

	orgName          string
	orgWebhookSecret string
	orgPoolParams    params.CreatePoolParams

	ghToken          string
	workflowFileName string

	cli         *client.GarmAPI
	authToken   runtime.ClientAuthInfoWriter
	repo        *params.Repository
	webhookInfo *params.HookInfo
	pool        *params.Pool
)

func initVars() error {
	adminPassword = os.Getenv("GARM_PASSWORD")
	adminUsername = os.Getenv("GARM_ADMIN_USERNAME")
	adminFullName = "GARM Admin"
	adminEmail = "admin@example.com"

	baseURL = os.Getenv("GARM_BASE_URL")
	credentialsName = os.Getenv("CREDENTIALS_NAME")
	repoWebhookSecret = os.Getenv("REPO_WEBHOOK_SECRET")

	repoPoolParams = params.CreatePoolParams{
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

	orgWebhookSecret = os.Getenv("ORG_WEBHOOK_SECRET")
	orgPoolParams = params.CreatePoolParams{
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

	orgName = os.Getenv("ORG_NAME")
	repoName = os.Getenv("REPO_NAME")

	ghToken = os.Getenv("GH_TOKEN")

	return nil
}

func initClient() error {
	garmURL, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	apiPath, err := url.JoinPath(garmURL.Path, client.DefaultBasePath)
	if err != nil {
		return err
	}
	transportCfg := client.DefaultTransportConfig().
		WithHost(garmURL.Host).
		WithBasePath(apiPath).
		WithSchemes([]string{garmURL.Scheme})
	cli = client.NewHTTPClientWithConfig(nil, transportCfg)

	return nil
}

func createUserAndLogin() error {
	newUser := params.NewUserParams{
		Username: adminUsername,
		Password: adminPassword,
		FullName: adminFullName,
		Email:    adminEmail,
	}
	_, err := firstRun(cli, newUser)
	if err != nil {
		return err
	}

	loginParams := params.PasswordLoginParams{
		Username: adminUsername,
		Password: adminPassword,
	}
	token, err := login(cli, loginParams)
	if err != nil {
		return err
	}
	authToken = openapiRuntimeClient.BearerToken(token)

	return nil
}

func initRepo(orgName, repoName, credentialsName, repoWebhookSecret string) error {
	createParams := params.CreateRepoParams{
		Owner:           orgName,
		Name:            repoName,
		CredentialsName: credentialsName,
		WebhookSecret:   repoWebhookSecret,
	}
	repository, err := createRepo(cli, authToken, createParams)
	if err != nil {
		return err
	}
	repo = repository
	return nil
}

func TestMain(m *testing.M) {
	// Clean up
	defer GracefulCleanup(&testing.T{})
	defer GHCleanupTest()

	// Initialize variables
	err := initVars()
	if err != nil {
		log.Fatalf("failed to initialize variables: %s", err)
	}

	// Initialize API client
	err = initClient()
	if err != nil {
		log.Fatalf("failed to initialize client: %s", err)
	}

	// Create user and login
	err = createUserAndLogin()
	if err != nil {
		log.Fatalf("failed to create user and login: %s", err)
	}

	// Initialize repo
	err = initRepo(orgName, repoName, credentialsName, repoWebhookSecret)
	if err != nil {
		log.Fatalf("failed to initialize repo: %s", err)
	}

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func GHCleanupTest() {
	controllerID, ctrlIDFound := os.LookupEnv("GARM_CONTROLLER_ID")
	if ctrlIDFound {
		_ = GhOrgRunnersCleanup(ghToken, orgName, controllerID)
		_ = GhRepoRunnersCleanup(ghToken, orgName, repoName, controllerID)
	} else {
		slog.Warn("Env variable GARM_CONTROLLER_ID is not set, skipping GitHub runners cleanup")
	}

	baseURL, baseURLFound := os.LookupEnv("GARM_BASE_URL")
	if ctrlIDFound && baseURLFound {
		webhookURL := fmt.Sprintf("%s/webhooks/%s", baseURL, controllerID)
		_ = GhOrgWebhookCleanup(ghToken, webhookURL, orgName)
		_ = GhRepoWebhookCleanup(ghToken, webhookURL, orgName, repoName)
	} else {
		slog.Warn("Env variables GARM_CONTROLLER_ID & GARM_BASE_URL are not set, skipping webhooks cleanup")
	}
}
