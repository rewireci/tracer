"""Integration tests for the pytest plugin using pytester."""
import pytest

import rewire_tracer


@pytest.fixture(autouse=True)
def reset_rewire_state() -> None:
    # Plugin state is module-global; reset between pytester runs to prevent
    # _shutdown from leaking across tests.
    rewire_tracer._reset()
    yield
    rewire_tracer._reset()


def test_plugin_does_nothing_without_token(pytester: pytest.Pytester) -> None:
    pytester.makepyfile("""
        def test_example():
            assert True
    """)
    # Plugin is already registered via pytest11 entry point; no -p flag needed.
    result = pytester.runpytest()
    result.assert_outcomes(passed=1)


def test_plugin_runs_tests_normally_with_endpoint(
    pytester: pytest.Pytester, monkeypatch: pytest.MonkeyPatch
) -> None:
    monkeypatch.setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
    pytester.makepyfile("""
        def test_pass():
            assert True

        def test_fail():
            assert False

        def test_skip():
            import pytest
            pytest.skip("skipped")
    """)
    result = pytester.runpytest()
    result.assert_outcomes(passed=1, failed=1, skipped=1)


def test_plugin_does_not_suppress_failures(
    pytester: pytest.Pytester, monkeypatch: pytest.MonkeyPatch
) -> None:
    monkeypatch.setenv("REWIRE_TOKEN", "rwt_test")
    pytester.makepyfile("""
        def test_should_fail():
            raise ValueError("boom")
    """)
    result = pytester.runpytest()
    result.assert_outcomes(failed=1)
