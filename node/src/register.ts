// --require entrypoint. Must be synchronous — NodeSDK.start() is sync so this is safe.
// Compiled as CJS so it works with Node.js --require flag.

import { init, shutdown } from "./index.js";

init();

// Flush spans when the event loop drains (e.g. after jest/vitest finishes all tests).
// NodeSDK registers its own SIGTERM handler; this covers normal process exit.
process.once("beforeExit", () => {
  shutdown().catch(console.error);
});
