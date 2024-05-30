package e2e

import (
	"fmt"
	"testing"
	"time"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/params"
)

func ValidateJobLifecycle(t *testing.T, label string) {
	t.Log("Validate GARM job lifecycle", "label", label)

	TriggerWorkflow(ghToken, orgName, repoName, workflowFileName, label)
	// wait for job list to be updated
	job, err := waitLabelledJob(t, label, 4*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	// check expected job status
	job, err = waitJobStatus(t, job.ID, params.JobStatusQueued, 4*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	job, err = waitJobStatus(t, job.ID, params.JobStatusInProgress, 4*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	// check expected instance status
	instance, err := waitInstanceStatus(t, job.RunnerName, commonParams.InstanceRunning, params.RunnerActive, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	// wait for job to be completed
	_, err = waitJobStatus(t, job.ID, params.JobStatusCompleted, 4*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	// wait for instance to be removed
	err = WaitInstanceToBeRemoved(t, instance.Name, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	// wait for GARM to rebuild the pool running idle instances
	err = WaitPoolInstances(t, instance.PoolID, commonParams.InstanceRunning, params.RunnerIdle, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
}

func waitLabelledJob(t *testing.T, label string, timeout time.Duration) (*params.Job, error) {
	var timeWaited time.Duration // default is 0
	var jobs params.Jobs
	var err error

	t.Log("Waiting for job", "label", label)
	for timeWaited < timeout {
		jobs, err = listJobs(cli, authToken)
		if err != nil {
			return nil, err
		}
		for _, job := range jobs {
			for _, jobLabel := range job.Labels {
				if jobLabel == label {
					return &job, err
				}
			}
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}

	if err := printJSONResponse(jobs); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("failed to wait job with label %s", label)
}

func waitJobStatus(t *testing.T, id int64, status params.JobStatus, timeout time.Duration) (*params.Job, error) {
	var timeWaited time.Duration // default is 0
	var job *params.Job

	t.Log("Waiting for job to reach status", "job_id", id, "status", status)
	for timeWaited < timeout {
		jobs, err := listJobs(cli, authToken)
		if err != nil {
			return nil, err
		}

		job = nil
		for k, v := range jobs {
			if v.ID == id {
				job = &jobs[k]
				break
			}
		}

		if job == nil {
			if status == params.JobStatusCompleted {
				// The job is not found in the list. We can safely assume
				// that it is completed
				return nil, nil
			}
			// if the job is not found, and expected status is not "completed",
			// we need to error out.
			return nil, fmt.Errorf("job %d not found, expected to be found in status %s", id, status)
		} else if job.Status == string(status) {
			return job, nil
		}
		time.Sleep(5 * time.Second)
		timeWaited += 5 * time.Second
	}

	if err := printJSONResponse(*job); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("timeout waiting for job %d to reach status %s", id, status)
}
