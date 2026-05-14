package rewire

import "os"

// CIContext holds the detected CI platform and run identifier.
type CIContext struct {
	Platform string
	RunID    string
}

// DetectCI inspects well-known CI environment variables and returns the active
// platform and run identifier. Platforms are checked in priority order:
// GitHub Actions, CircleCI, GitLab, then a generic REWIRE_RUN_ID fallback.
// If no known CI environment is detected, Platform is "unknown" and RunID is empty.
func DetectCI() CIContext {
	if runID := os.Getenv("GITHUB_RUN_ID"); runID != "" {
		return CIContext{Platform: "github_actions", RunID: runID}
	}
	if runID := os.Getenv("CIRCLE_WORKFLOW_ID"); runID != "" {
		return CIContext{Platform: "circleci", RunID: runID}
	}
	if runID := os.Getenv("CI_PIPELINE_ID"); runID != "" {
		return CIContext{Platform: "gitlab", RunID: runID}
	}
	if runID := os.Getenv("REWIRE_RUN_ID"); runID != "" {
		return CIContext{Platform: "generic", RunID: runID}
	}
	return CIContext{Platform: "unknown"}
}
