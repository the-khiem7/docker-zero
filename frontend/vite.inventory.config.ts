import { fileURLToPath } from "node:url";
import { defineConfig, mergeConfig, type Plugin } from "vite";
import base from "./vite.config.ts";

// Serves the real frontend with the bound Go services replaced by the
// scenario mocks in inventory/, so every state of every window can be shown
// in a plain browser: /?view=<window>&state=<scenario>.
const mockDir = fileURLToPath(new URL("./inventory/mock/", import.meta.url));
const bound = /\/internal\/ui\/(service|settingsservice|setupservice|adoptservice|maintainservice|uninstallservice)(\.js)?$/;

function mockServices(): Plugin {
  return {
    name: "docker-zero-inventory-mocks",
    enforce: "pre",
    resolveId(source, importer) {
      const m = bound.exec(source);
      if (!m || importer?.replaceAll("\\", "/").includes("/inventory/mock/")) return null;
      return `${mockDir}${m[1]}.ts`;
    },
  };
}

export default mergeConfig(
  base,
  defineConfig({
    server: { port: Number(process.env.INVENTORY_PORT) || 9246 },
    plugins: [mockServices()],
  }),
);
