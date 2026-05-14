import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { init, shutdown, _reset } from "../src/index.js";

beforeEach(() => {
  _reset();
  delete process.env.OTEL_EXPORTER_OTLP_ENDPOINT;
  delete process.env.REWIRE_TOKEN;
  delete process.env.GITHUB_RUN_ID;
  delete process.env.OTEL_SERVICE_NAME;
  delete process.env.GITHUB_REPOSITORY;
});

afterEach(() => {
  _reset();
});

describe("init — no configuration", () => {
  it("warns and returns a no-op when no env vars are set", () => {
    const warn = vi.spyOn(console, "warn").mockImplementation(() => {});
    const stop = init();

    expect(warn).toHaveBeenCalledWith(
      expect.stringContaining("OTEL_EXPORTER_OTLP_ENDPOINT"),
    );
    expect(stop).toBeTypeOf("function");
    warn.mockRestore();
  });

  it("no-op shutdown resolves without error", async () => {
    vi.spyOn(console, "warn").mockImplementation(() => {});
    const stop = init();
    await expect(stop()).resolves.toBeUndefined();
  });

  it("returns the same shutdown function when called twice", () => {
    vi.spyOn(console, "warn").mockImplementation(() => {});
    const a = init();
    const b = init();
    expect(a).toBe(b);
  });
});

describe("shutdown — before init", () => {
  it("resolves without error", async () => {
    await expect(shutdown()).resolves.toBeUndefined();
  });
});

describe("init — with REWIRE_TOKEN", () => {
  it("initializes and exports traces when token is set", () => {
    process.env.REWIRE_TOKEN = "rwt_test";

    const warn = vi.spyOn(console, "warn").mockImplementation(() => {});
    const stop = init();

    expect(stop).toBeTypeOf("function");
    expect(warn).not.toHaveBeenCalled();
    warn.mockRestore();
  });
});

describe("init — with OTEL_EXPORTER_OTLP_ENDPOINT", () => {
  it("initializes without error and returns a shutdown function", async () => {
    process.env.OTEL_EXPORTER_OTLP_ENDPOINT = "http://localhost:4318";
    process.env.GITHUB_RUN_ID = "42";

    const warn = vi.spyOn(console, "warn").mockImplementation(() => {});
    const stop = init();

    expect(stop).toBeTypeOf("function");
    // No warning should be emitted — we have a valid endpoint
    expect(warn).not.toHaveBeenCalled();

    await stop();
    warn.mockRestore();
  });
});
