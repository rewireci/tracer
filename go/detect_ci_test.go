package rewire

import (
	"os"
	"testing"
)

var ciEnvVars = []string{"GITHUB_RUN_ID", "CIRCLE_WORKFLOW_ID", "CI_PIPELINE_ID", "REWIRE_RUN_ID"}

func clearCIVars(t *testing.T) {
	t.Helper()
	saved := make(map[string]string)
	for _, k := range ciEnvVars {
		saved[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for k, v := range saved {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})
}

func TestDetectCI_GitHub(t *testing.T) {
	clearCIVars(t)
	os.Setenv("GITHUB_RUN_ID", "12345")
	ci := DetectCI()
	if ci.Platform != "github_actions" || ci.RunID != "12345" {
		t.Errorf("got {%s %s}, want {github_actions 12345}", ci.Platform, ci.RunID)
	}
}

func TestDetectCI_CircleCI(t *testing.T) {
	clearCIVars(t)
	os.Setenv("CIRCLE_WORKFLOW_ID", "abc-workflow")
	ci := DetectCI()
	if ci.Platform != "circleci" || ci.RunID != "abc-workflow" {
		t.Errorf("got {%s %s}, want {circleci abc-workflow}", ci.Platform, ci.RunID)
	}
}

func TestDetectCI_GitLab(t *testing.T) {
	clearCIVars(t)
	os.Setenv("CI_PIPELINE_ID", "99")
	ci := DetectCI()
	if ci.Platform != "gitlab" || ci.RunID != "99" {
		t.Errorf("got {%s %s}, want {gitlab 99}", ci.Platform, ci.RunID)
	}
}

func TestDetectCI_Generic(t *testing.T) {
	clearCIVars(t)
	os.Setenv("REWIRE_RUN_ID", "custom-run")
	ci := DetectCI()
	if ci.Platform != "generic" || ci.RunID != "custom-run" {
		t.Errorf("got {%s %s}, want {generic custom-run}", ci.Platform, ci.RunID)
	}
}

func TestDetectCI_Unknown(t *testing.T) {
	clearCIVars(t)
	ci := DetectCI()
	if ci.Platform != "unknown" || ci.RunID != "" {
		t.Errorf("got {%s %q}, want {unknown \"\"}", ci.Platform, ci.RunID)
	}
}

func TestDetectCI_PrefersGitHub(t *testing.T) {
	clearCIVars(t)
	os.Setenv("GITHUB_RUN_ID", "gh-1")
	os.Setenv("CIRCLE_WORKFLOW_ID", "ci-1")
	ci := DetectCI()
	if ci.Platform != "github_actions" || ci.RunID != "gh-1" {
		t.Errorf("got {%s %s}, want {github_actions gh-1}", ci.Platform, ci.RunID)
	}
}
