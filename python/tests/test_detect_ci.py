import pytest

from rewire_tracer.detect_ci import CiContext, detect_ci


@pytest.fixture(autouse=True)
def clean_env(monkeypatch: pytest.MonkeyPatch) -> None:
    for var in ("GITHUB_RUN_ID", "CIRCLE_WORKFLOW_ID", "CI_PIPELINE_ID", "REWIRE_RUN_ID"):
        monkeypatch.delenv(var, raising=False)


def test_detects_github_actions(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("GITHUB_RUN_ID", "12345")
    assert detect_ci() == CiContext(platform="github_actions", run_id="12345")


def test_detects_circleci(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("CIRCLE_WORKFLOW_ID", "abc-workflow")
    assert detect_ci() == CiContext(platform="circleci", run_id="abc-workflow")


def test_detects_gitlab(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("CI_PIPELINE_ID", "99")
    assert detect_ci() == CiContext(platform="gitlab", run_id="99")


def test_detects_generic_rewire_run_id(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("REWIRE_RUN_ID", "custom-run")
    assert detect_ci() == CiContext(platform="generic", run_id="custom-run")


def test_unknown_when_no_ci_vars_set() -> None:
    assert detect_ci() == CiContext(platform="unknown", run_id=None)


def test_prefers_github_actions_over_other_platforms(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("GITHUB_RUN_ID", "gh-1")
    monkeypatch.setenv("CIRCLE_WORKFLOW_ID", "ci-1")
    assert detect_ci() == CiContext(platform="github_actions", run_id="gh-1")
