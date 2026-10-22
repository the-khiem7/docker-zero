// Serves the real frontend with the bound Go services resolved to the gallery
// mocks in gallery/mock/, so every state of every window can be shown in a
// plain browser: /?view=<window>&state=<scenario>.
import { fileURLToPath } from "node:url";
import { defineConfig, mergeConfig, type Plugin } from "vite";
import base from "./vite.config.ts";

const mockDir = fileURLToPath(new URL("./gallery/mock/", import.meta.url));

const MOCKS: [RegExp, (m: RegExpExecArray) => string][] = [
  [
    /\/internal\/ui\/(service|settingsservice|setupservice|adoptservice|maintainservice|uninstallservice)(\.js)?$/,
    (m) => `${m[1]}.ts`,
  ],
];

function galleryMocks(): Plugin {
  return {
    name: "ui-gallery-mocks",
    enforce: "pre",
    resolveId(source, importer) {
      // A mock may import the real module's types; never redirect its own imports.
      if (importer?.replaceAll("\\", "/").includes("/gallery/mock/")) return null;
      for (const [pattern, file] of MOCKS) {
        const m = pattern.exec(source);
        if (m) return mockDir + file(m);
      }
      return null;
    },
  };
}

export default mergeConfig(
  base,
  defineConfig({
    server: { port: Number(process.env.GALLERY_PORT) || 9246, strictPort: true },
    plugins: [galleryMocks()],
  }),
);
