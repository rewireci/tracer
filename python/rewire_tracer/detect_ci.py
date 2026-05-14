from __future__ import annotations

import os
from typing import NamedTuple


class CiContext(NamedTuple):
    platform: str
    run_id: str | None


def detect_ci() -> CiContext:
    if run_id := os.environ.get("GITHUB_RUN_ID"):
        return CiContext(platform="github_actions", run_id=run_id)
    if run_id := os.environ.get("CIRCLE_WORKFLOW_ID"):
        return CiContext(platform="circleci", run_id=run_id)
    if run_id := os.environ.get("CI_PIPELINE_ID"):
        return CiContext(platform="gitlab", run_id=run_id)
    if run_id := os.environ.get("REWIRE_RUN_ID"):
        return CiContext(platform="generic", run_id=run_id)
    return CiContext(platform="unknown", run_id=None)
