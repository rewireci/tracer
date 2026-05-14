import { defineConfig } from "tsup";

export default defineConfig([
  {
    entry: { index: "src/index.ts" },
    format: ["esm"],
    dts: true,
    clean: true,
    target: "node20",
    shims: true,
  },
  {
    entry: { register: "src/register.ts" },
    format: ["cjs"],
    dts: false,
    target: "node20",
    shims: true,
    outExtension: () => ({ js: ".cjs" }),
  },
]);
