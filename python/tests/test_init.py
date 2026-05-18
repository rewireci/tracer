import warnings

import pytest

from rewire_tracer import _reset, init, shutdown


@pytest.fixture(autouse=True)
def clean_state(monkeypatch: pytest.MonkeyPatch) -> None:
    _reset()
    for var in (
        "OTEL_EXPORTER_OTLP_ENDPOINT",
        "REWIRE_TOKEN",
        "GITHUB_RUN_ID",
        "OTEL_SERVICE_NAME",
        "GITHUB_REPOSITORY",
    ):
        monkeypatch.delenv(var, raising=False)
    yield
    _reset()


class TestInitNoConfiguration:
    def test_warns_and_returns_no_op(self) -> None:
        with warnings.catch_warnings(record=True) as caught:
            warnings.simplefilter("always")
            stop = init()
        messages = [str(w.message) for w in caught]
        assert any("OTEL_EXPORTER_OTLP_ENDPOINT" in m for m in messages)
        assert callable(stop)

    def test_no_op_shutdown_does_not_raise(self) -> None:
        with warnings.catch_warnings(record=True):
            warnings.simplefilter("always")
            stop = init()
        stop()  # must not raise

    def test_called_twice_returns_same_function(self) -> None:
        with warnings.catch_warnings(record=True):
            warnings.simplefilter("always")
            a = init()
            b = init()
        assert a is b


class TestShutdownBeforeInit:
    def test_resolves_without_error(self) -> None:
        shutdown()  # must not raise


class TestInitWithOtelEndpoint:
    def test_initializes_without_error(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
        monkeypatch.setenv("GITHUB_RUN_ID", "42")

        with warnings.catch_warnings(record=True) as caught:
            warnings.simplefilter("always")
            stop = init()

        assert callable(stop)
        assert not any("OTEL_EXPORTER_OTLP_ENDPOINT" in str(w.message) for w in caught), (
            "Should not warn when endpoint is set"
        )
        stop()

    def test_called_twice_returns_same_function(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
        a = init()
        b = init()
        assert a is b
        a()


class TestInitWithRewireToken:
    def test_initializes_without_error(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv("REWIRE_TOKEN", "rwt_test")
        with warnings.catch_warnings(record=True):
            warnings.simplefilter("always")
            stop = init()
        assert callable(stop)
        stop()

    def test_configures_exporter_with_correct_traces_path(
        self, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        monkeypatch.setenv("REWIRE_TOKEN", "rwt_test")
        captured: list[str] = []

        from unittest.mock import patch

        from opentelemetry.exporter.otlp.proto.http.trace_exporter import (
            OTLPSpanExporter,
        )

        original_init = OTLPSpanExporter.__init__

        def capture_init(self: OTLPSpanExporter, *args: object, **kwargs: object) -> None:
            captured.append(str(kwargs.get("endpoint", "")))
            original_init(self, *args, **kwargs)

        with patch.object(OTLPSpanExporter, "__init__", capture_init):
            stop = init()
        stop()
        assert len(captured) == 1
        assert captured[0] == "https://rewireci.com/otlp/v1/traces"
