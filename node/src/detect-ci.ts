export interface CiContext {
  platform: string;
  runId: string | undefined;
}

export function detectCi(): CiContext {
  if (process.env.GITHUB_RUN_ID) {
    return { platform: "github_actions", runId: process.env.GITHUB_RUN_ID };
  }
  if (process.env.CIRCLE_WORKFLOW_ID) {
    return { platform: "circleci", runId: process.env.CIRCLE_WORKFLOW_ID };
  }
  if (process.env.CI_PIPELINE_ID) {
    return { platform: "gitlab", runId: process.env.CI_PIPELINE_ID };
  }
  if (process.env.REWIRE_RUN_ID) {
    return { platform: "generic", runId: process.env.REWIRE_RUN_ID };
  }
  return { platform: "unknown", runId: undefined };
}
