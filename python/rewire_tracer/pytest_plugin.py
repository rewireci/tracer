from __future__ import annotations

from typing import Any

import pytest

from rewire_tracer import _enabled

_spans: dict[str, Any] = {}


def pytest_configure(config: pytest.Config) -> None:
    from rewire_tracer import init

    init()


def pytest_unconfigure(config: pytest.Config) -> None:
    from rewire_tracer import shutdown

    shutdown()


@pytest.hookimpl(wrapper=True)
def pytest_runtest_protocol(item: pytest.Item, nextitem: pytest.Item | None) -> Any:
    if not _enabled:
        return (yield)
    try:
        from opentelemetry import trace

        tracer = trace.get_tracer("rewire.pytest")
        with tracer.start_as_current_span(
            item.nodeid,
            attributes={
                "test.name": item.name,
                "test.file": str(item.path),
            },
        ) as span:
            _spans[item.nodeid] = span
            return (yield)
    except ImportError:
        return (yield)
    finally:
        _spans.pop(item.nodeid, None)


def pytest_runtest_logreport(report: pytest.TestReport) -> None:
    span = _spans.get(report.nodeid)
    if span is None:
        return
    try:
        from opentelemetry.trace import StatusCode

        if report.when == "call":
            if report.failed:
                span.set_status(StatusCode.ERROR)
                span.set_attribute("test.status", "failed")
            elif report.skipped:
                span.set_attribute("test.status", "skipped")
            else:
                span.set_status(StatusCode.OK)
                span.set_attribute("test.status", "passed")
        elif report.when == "setup":
            if report.failed:
                span.set_status(StatusCode.ERROR)
                span.set_attribute("test.status", "failed")
            elif report.skipped:
                span.set_attribute("test.status", "skipped")
        elif report.when == "teardown":
            if report.failed:
                span.set_status(StatusCode.ERROR)
                span.set_attribute("test.status", "failed")
    except ImportError:
        pass
