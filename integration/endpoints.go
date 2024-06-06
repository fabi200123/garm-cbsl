package integration

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/cloudbase/garm/params"
)

const (
	defaultEndpointName  string = "github.com"
	dummyCredentialsName string = "dummy"
)

func MustDefaultGithubEndpoint() error {
	ep, err := GetGithubEndpoint("github.com")
	if err != nil {
		return err
	}
	if ep == nil {
		return fmt.Errorf("default GitHub endpoint not found")
	}

	if ep.Name != "github.com" {
		return fmt.Errorf("default GitHub endpoint name mismatch")
	}

	return nil
}

func checkEndpointParamsAreEqual(a, b params.GithubEndpoint) error {
	if a.Name != b.Name {
		return fmt.Errorf("endpoint name mismatch")
	}

	if a.Description != b.Description {
		return fmt.Errorf("endpoint description mismatch")
	}

	if a.BaseURL != b.BaseURL {
		return fmt.Errorf("endpoint base URL mismatch")
	}

	if a.APIBaseURL != b.APIBaseURL {
		return fmt.Errorf("endpoint API base URL mismatch")
	}

	if a.UploadBaseURL != b.UploadBaseURL {
		return fmt.Errorf("endpoint upload base URL mismatch")
	}

	if string(a.CACertBundle) != string(b.CACertBundle) {
		return fmt.Errorf("endpoint CA cert bundle mismatch")
	}
	return nil
}

func getTestFileContents(relPath string) ([]byte, error) {
	baseDir := os.Getenv("GARM_CHECKOUT_DIR")
	if baseDir == "" {
		return nil, fmt.Errorf("ariable GARM_CHECKOUT_DIR not set")
	}
	contents, err := os.ReadFile(filepath.Join(baseDir, "testdata", relPath))
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func createDummyEndpoint(name string) (*params.GithubEndpoint, error) {
	endpointParams := params.CreateGithubEndpointParams{
		Name:          name,
		Description:   "Dummy endpoint",
		BaseURL:       "https://ghes.example.com",
		APIBaseURL:    "https://api.ghes.example.com/",
		UploadBaseURL: "https://uploads.ghes.example.com/",
	}

	return CreateGithubEndpoint(endpointParams)
}

func CreateGithubEndpoint(endpointParams params.CreateGithubEndpointParams) (*params.GithubEndpoint, error) {
	slog.Info("Create GitHub endpoint")
	endpoint, err := createGithubEndpoint(cli, authToken, endpointParams)
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}
