from __future__ import annotations

import os
import warnings
from typing import Callable

from .detect_ci import detect_ci

REWIRE_ENDPOINT = "https://rewireci.com/otlp/v1"

_shutdown: Callable[[], None] | None = None
_enabled: bool = False


def _noop() -> None:
    pass


def init() -> Callable[[], None]:
    global _shutdown, _enabled
    if _shutdown is not None:
        return _shutdown

    endpoint = os.environ.get("OTEL_EXPORTER_OTLP_ENDPOINT")
    token = os.environ.get("REWIRE_TOKEN")

    if not endpoint and not token:
        warnings.warn(
            "[rewire] Neither OTEL_EXPORTER_OTLP_ENDPOINT nor REWIRE_TOKEN is set"
            " — tracing disabled",
            stacklevel=2,
        )
        _shutdown = _noop
        return _shutdown

    try:
        from opentelemetry import trace
        from opentelemetry.exporter.otlp.proto.http.trace_exporter import (
            OTLPSpanExporter,
        )
        from opentelemetry.sdk.resources import Resource
        from opentelemetry.sdk.trace import TracerProvider
        from opentelemetry.sdk.trace.export import BatchSpanProcessor
    except ImportError:
        warnings.warn(
            "[rewire] opentelemetry-sdk is not installed — tracing disabled",
            stacklevel=2,
        )
        _shutdown = _noop
        return _shutdown

    try:
        ci = detect_ci()
        base_endpoint = endpoint or REWIRE_ENDPOINT
        headers: dict[str, str] = {}
        if not endpoint and token:
            headers["Authorization"] = f"Bearer {token}"

        resource_attrs: dict[str, str] = {
            "ci.platform": ci.platform,
            "service.name": (
                os.environ.get("OTEL_SERVICE_NAME")
                or os.environ.get("GITHUB_REPOSITORY")
                or "unknown"
            ),
        }
        if ci.run_id:
            resource_attrs["run.id"] = ci.run_id

        resource = Resource(attributes=resource_attrs)
        provider = TracerProvider(resource=resource)
        if endpoint:
            trace_endpoint = f"{endpoint.rstrip('/')}/v1/traces"
        else:
            trace_endpoint = f"{REWIRE_ENDPOINT}/traces"
        exporter = OTLPSpanExporter(
            endpoint=trace_endpoint,
            headers=headers,
        )
        provider.add_span_processor(BatchSpanProcessor(exporter))
        trace.set_tracer_provider(provider)

        def _do_shutdown() -> None:
            provider.shutdown()

        _shutdown = _do_shutdown
        _enabled = True
        return _shutdown
    except Exception as exc:
        warnings.warn(f"[rewire] Failed to initialize OTel SDK: {exc}", stacklevel=2)
        _shutdown = _noop
        return _shutdown


def shutdown() -> None:
    if _shutdown is not None:
        _shutdown()


def _reset() -> None:
    global _shutdown, _enabled
    _shutdown = None
    _enabled = False
