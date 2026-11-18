package github

import (
	"encoding/json"
	"testing"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/stretchr/testify/require"
)

func TestWorkflowJobStepsFrame(t *testing.T) {
	var job googlegithub.WorkflowJob
	require.NoError(t, json.Unmarshal([]byte(`{"id":45,"name":"Build (linux)","run_attempt":2,"steps":[{"name":"Build","number":1,"status":"completed","conclusion":"success","started_at":"2026-11-17T08:00:00Z","completed_at":"2026-11-17T08:02:30Z"},{"name":"Publish","number":2,"status":"in_progress","started_at":"2026-11-17T08:03:00Z"},{"name":"Optional","number":3,"status":"completed","conclusion":"skipped","started_at":"2026-11-17T08:03:00Z","completed_at":"2026-11-17T08:03:00Z"}]}`), &job))
	run := &googlegithub.WorkflowRun{ID: googlegithub.Ptr(int64(12)), WorkflowID: googlegithub.Ptr(int64(34)), Name: googlegithub.Ptr("CI")}
	steps := workflowJobSteps{{run: run, job: &job}}
	frame := steps.Frames()[0]
	require.Equal(t, "workflow_job_steps", frame.Name)
	require.Equal(t, 3, frame.Rows())
	fields := map[string]int{}
	for i, field := range frame.Fields {
		fields[field.Name] = i
	}
	require.Equal(t, int64(34), frame.Fields[fields["workflow_id"]].At(0))
	require.Equal(t, int64(12), frame.Fields[fields["run_id"]].At(0))
	require.Equal(t, int64(2), frame.Fields[fields["run_attempt"]].At(0))
	require.Equal(t, "Build (linux)", frame.Fields[fields["job_name"]].At(0))
	duration := frame.Fields[fields["duration_seconds"]]
	require.Equal(t, "s", duration.Config.Unit)
	require.Equal(t, 150.0, *duration.At(0).(*float64))
	require.Nil(t, duration.At(1))
	require.Nil(t, duration.At(2))
}
