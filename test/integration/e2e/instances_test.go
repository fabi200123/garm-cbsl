package e2e

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/params"
)

func waitInstanceStatus(t *testing.T, name string, status commonParams.InstanceStatus, runnerStatus params.RunnerStatus, timeout time.Duration) (*params.Instance, error) {
	var timeWaited time.Duration // default is 0
	var instance *params.Instance
	var err error

	t.Log("Waiting for instance to reach desired status", "instance", name, "desired_status", status, "desired_runner_status", runnerStatus)
	for timeWaited < timeout {
		instance, err = getInstance(cli, authToken, name)
		if err != nil {
			return nil, err
		}
		t.Log("Instance status", "instance_name", name, "status", instance.Status, "runner_status", instance.RunnerStatus)
		if instance.Status == status && instance.RunnerStatus == runnerStatus {
			return instance, nil
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}

	if err := printJSONResponse(*instance); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("timeout waiting for instance %s status to reach status %s and runner status %s", name, status, runnerStatus)
}

func DeleteInstance(t *testing.T, name string, forceRemove, bypassGHUnauthorized bool) {
	t.Log("Delete instance", "instance_name", name, "force_remove", forceRemove)
	if err := deleteInstance(cli, authToken, name, forceRemove, bypassGHUnauthorized); err != nil {
		slog.Error("Failed to delete instance", "instance_name", name, "error", err)
		t.Fatal(err)
	}
	t.Log("Instance deletion initiated", "instance_name", name)
}

func WaitInstanceToBeRemoved(t *testing.T, name string, timeout time.Duration) error {
	var timeWaited time.Duration // default is 0
	var instance *params.Instance

	t.Log("Waiting for instance to be removed", "instance_name", name)
	for timeWaited < timeout {
		instances, err := listInstances(cli, authToken)
		if err != nil {
			return err
		}

		instance = nil
		for k, v := range instances {
			if v.Name == name {
				instance = &instances[k]
				break
			}
		}
		if instance == nil {
			// The instance is not found in the list. We can safely assume
			// that it is removed
			return nil
		}

		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}

	if err := printJSONResponse(*instance); err != nil {
		return err
	}
	return fmt.Errorf("instance %s was not removed within the timeout", name)
}

func WaitPoolInstances(t *testing.T, poolID string, status commonParams.InstanceStatus, runnerStatus params.RunnerStatus, timeout time.Duration) error {
	var timeWaited time.Duration // default is 0

	pool, err := getPool(cli, authToken, poolID)
	if err != nil {
		return err
	}

	t.Log("Waiting for pool instances to reach desired status", "pool_id", poolID, "desired_status", status, "desired_runner_status", runnerStatus)
	for timeWaited < timeout {
		poolInstances, err := listPoolInstances(cli, authToken, poolID)
		if err != nil {
			return err
		}

		instancesCount := 0
		for _, instance := range poolInstances {
			if instance.Status == status && instance.RunnerStatus == runnerStatus {
				instancesCount++
			}
		}

		t.Log(
			"Pool instance reached status",
			"pool_id", poolID,
			"status", status,
			"runner_status", runnerStatus,
			"desired_instance_count", instancesCount,
			"pool_instance_count", len(poolInstances))
		if int(pool.MinIdleRunners) == instancesCount {
			return nil
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}

	_ = dumpPoolInstancesDetails(t, pool.ID)

	return fmt.Errorf("timeout waiting for pool %s instances to reach status: %s and runner status: %s", poolID, status, runnerStatus)
}
