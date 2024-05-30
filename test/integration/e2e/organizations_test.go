package e2e

import (
	"testing"
	"time"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/params"
)

func CreateOrg(t *testing.T, orgName, credentialsName, orgWebhookSecret string) *params.Organization {
	t.Log("Create org", "org_name", orgName)
	orgParams := params.CreateOrgParams{
		Name:            orgName,
		CredentialsName: credentialsName,
		WebhookSecret:   orgWebhookSecret,
	}
	org, err := createOrg(cli, authToken, orgParams)
	if err != nil {
		t.Fatal(err)
	}
	return org
}

func UpdateOrg(t *testing.T, id, credentialsName string) *params.Organization {
	t.Log("Update org", "org_id", id)
	updateParams := params.UpdateEntityParams{
		CredentialsName: credentialsName,
	}
	org, err := updateOrg(cli, authToken, id, updateParams)
	if err != nil {
		t.Fatal(err)
	}
	return org
}

func InstallOrgWebhook(t *testing.T, id string) *params.HookInfo {
	t.Log("Install org webhook", "org_id", id)
	webhookParams := params.InstallWebhookParams{
		WebhookEndpointType: params.WebhookEndpointDirect,
	}
	_, err := installOrgWebhook(cli, authToken, id, webhookParams)
	if err != nil {
		t.Fatal(err)
	}
	webhookInfo, err := getOrgWebhook(cli, authToken, id)
	if err != nil {
		t.Fatal(err)
	}
	return webhookInfo
}

func UninstallOrgWebhook(t *testing.T, id string) {
	t.Log("Uninstall org webhook", "org_id", id)
	if err := uninstallOrgWebhook(cli, authToken, id); err != nil {
		t.Fatal(err)
	}
}

func CreateOrgPool(t *testing.T, orgID string, poolParams params.CreatePoolParams) *params.Pool {
	t.Log("Create org pool", "org_id", orgID)
	pool, err := createOrgPool(cli, authToken, orgID, poolParams)
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func GetOrgPool(t *testing.T, orgID, orgPoolID string) *params.Pool {
	t.Log("Get org pool", "org_id", orgID, "pool_id", orgPoolID)
	pool, err := getOrgPool(cli, authToken, orgID, orgPoolID)
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func UpdateOrgPool(t *testing.T, orgID, orgPoolID string, maxRunners, minIdleRunners uint) *params.Pool {
	t.Log("Update org pool", "org_id", orgID, "pool_id", orgPoolID)
	poolParams := params.UpdatePoolParams{
		MinIdleRunners: &minIdleRunners,
		MaxRunners:     &maxRunners,
	}
	pool, err := updateOrgPool(cli, authToken, orgID, orgPoolID, poolParams)
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func DeleteOrgPool(t *testing.T, orgID, orgPoolID string) {
	t.Log("Delete org pool", "org_id", orgID, "pool_id", orgPoolID)
	if err := deleteOrgPool(cli, authToken, orgID, orgPoolID); err != nil {
		t.Fatal(err)
	}
}

func WaitOrgRunningIdleInstances(t *testing.T, orgID string, timeout time.Duration) {
	orgPools, err := listOrgPools(cli, authToken, orgID)
	if err != nil {
		t.Fatal(err)
	}
	for _, pool := range orgPools {
		err := WaitPoolInstances(t, pool.ID, commonParams.InstanceRunning, params.RunnerIdle, timeout)
		if err != nil {
			_ = dumpOrgInstancesDetails(t, orgID)
			t.Fatal(err)
		}
	}
}

func dumpOrgInstancesDetails(t *testing.T, orgID string) error {
	// print org details
	t.Log("Dumping org details", "org_id", orgID)
	org, err := getOrg(cli, authToken, orgID)
	if err != nil {
		return err
	}
	if err := printJSONResponse(org); err != nil {
		return err
	}

	// print org instances details
	t.Log("Dumping org instances details", "org_id", orgID)
	instances, err := listOrgInstances(cli, authToken, orgID)
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
