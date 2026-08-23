import { defineConfig } from "oxfmt";
import ultracite from "ultracite/oxfmt";

export default defineConfig({
  extends: [ultracite],
  ignorePatterns: [
    ...ultracite.ignorePatterns,
    ".devcontainer/devcontainer-lock.json",
    "website/worker-configuration.d.ts",
  ],
});
