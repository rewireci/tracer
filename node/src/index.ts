import { createRequire } from "node:module";

import { detectCi } from "./detect-ci.js";

const _require = createRequire(import.meta.url);

const REWIRE_ENDPOINT = "https://rewireci.com";

export type ShutdownFn = () => Promise<void>;

let _shutdown: ShutdownFn | undefined;

function tryRequire(id: string): unknown {
  try {
    return _require(id);
  } catch {
    return null;
  }
}

export function init(): ShutdownFn {
  if (_shutdown) return _shutdown;

  const endpoint = process.env.OTEL_EXPORTER_OTLP_ENDPOINT;
  const token = process.env.REWIRE_TOKEN;

  if (!endpoint && !token) {
    console.warn(
      "[rewire] Neither OTEL_EXPORTER_OTLP_ENDPOINT nor REWIRE_TOKEN is set — tracing disabled",
    );
    _shutdown = () => Promise.resolve();
    return _shutdown;
  }

  const sdkModule = tryRequire("@opentelemetry/sdk-node") as {
    NodeSDK: new (opts: unknown) => {
      start(): void;
      shutdown(): Promise<void>;
    };
  } | null;

  if (!sdkModule) {
    console.warn(
      "[rewire] @opentelemetry/sdk-node is not installed — tracing disabled",
    );
    _shutdown = () => Promise.resolve();
    return _shutdown;
  }

  try {
    const { NodeSDK } = sdkModule;

    const exporterModule = tryRequire(
      "@opentelemetry/exporter-trace-otlp-http",
    ) as { OTLPTraceExporter: new (opts: unknown) => unknown } | null;

    if (!exporterModule) {
      console.warn(
        "[rewire] @opentelemetry/exporter-trace-otlp-http is not installed — tracing disabled",
      );
      _shutdown = () => Promise.resolve();
      return _shutdown;
    }

    const resourcesModule = tryRequire("@opentelemetry/resources") as {
      Resource: new (attrs: Record<string, string>) => unknown;
    } | null;

    if (!resourcesModule) {
      console.warn(
        "[rewire] @opentelemetry/resources is not installed — tracing disabled",
      );
      _shutdown = () => Promise.resolve();
      return _shutdown;
    }

    const { OTLPTraceExporter } = exporterModule;
    const { Resource } = resourcesModule;

    const autoInstrModule = tryRequire(
      "@opentelemetry/auto-instrumentations-node",
    ) as { getNodeAutoInstrumentations: () => unknown[] } | null;

    const ci = detectCi();
    const baseUrl = endpoint ?? REWIRE_ENDPOINT;
    const headers: Record<string, string> = {};
    if (!endpoint && token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const resourceAttrs: Record<string, string> = {
      "ci.platform": ci.platform,
      "service.name":
        process.env.OTEL_SERVICE_NAME ??
        process.env.GITHUB_REPOSITORY ??
        "unknown",
    };
    if (ci.runId) resourceAttrs["run.id"] = ci.runId;

    const traceUrl = endpoint
      ? `${endpoint.replace(/\/$/, "")}/v1/traces`
      : `${baseUrl}/otlp/v1/traces`;

    const sdk = new NodeSDK({
      traceExporter: new OTLPTraceExporter({
        url: traceUrl,
        headers,
      }),
      instrumentations: autoInstrModule
        ? autoInstrModule.getNodeAutoInstrumentations()
        : [],
      resource: new Resource(resourceAttrs),
    });

    sdk.start();

    _shutdown = () => sdk.shutdown();
    return _shutdown;
  } catch (err) {
    console.warn("[rewire] Failed to initialize OTel SDK:", err);
    _shutdown = () => Promise.resolve();
    return _shutdown;
  }
}

export function shutdown(): Promise<void> {
  return _shutdown ? _shutdown() : Promise.resolve();
}

/**
 * @internal — for testing only. Do not use in application code.
 */
export function _reset(): void {
  _shutdown = undefined;
}
