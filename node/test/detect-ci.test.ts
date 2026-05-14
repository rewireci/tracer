import { afterEach, describe, expect, it } from "vitest";

import { detectCi } from "../src/detect-ci.js";

afterEach(() => {
  delete process.env.GITHUB_RUN_ID;
  delete process.env.CIRCLE_WORKFLOW_ID;
  delete process.env.CI_PIPELINE_ID;
  delete process.env.REWIRE_RUN_ID;
});

describe("detectCi", () => {
  it("detects GitHub Actions", () => {
    process.env.GITHUB_RUN_ID = "12345";
    expect(detectCi()).toEqual({ platform: "github_actions", runId: "12345" });
  });

  it("detects CircleCI", () => {
    process.env.CIRCLE_WORKFLOW_ID = "abc-workflow";
    expect(detectCi()).toEqual({
      platform: "circleci",
      runId: "abc-workflow",
    });
  });

  it("detects GitLab", () => {
    process.env.CI_PIPELINE_ID = "99";
    expect(detectCi()).toEqual({ platform: "gitlab", runId: "99" });
  });

  it("detects generic REWIRE_RUN_ID", () => {
    process.env.REWIRE_RUN_ID = "custom-run";
    expect(detectCi()).toEqual({ platform: "generic", runId: "custom-run" });
  });

  it("returns unknown when no CI env vars are set", () => {
    expect(detectCi()).toEqual({ platform: "unknown", runId: undefined });
  });

  it("prefers GitHub Actions over other platforms", () => {
    process.env.GITHUB_RUN_ID = "gh-1";
    process.env.CIRCLE_WORKFLOW_ID = "ci-1";
    expect(detectCi()).toEqual({ platform: "github_actions", runId: "gh-1" });
  });
});
