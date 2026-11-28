import assert from "node:assert/strict";
import { describe, it } from "node:test";

import worker from "./index.ts";

const ASSET_BODY = "static asset";

const assetRequests: string[] = [];

const env = {
  ASSETS: {
    fetch: (request: Request) => {
      assetRequests.push(new URL(request.url).pathname);
      return Promise.resolve(new Response(ASSET_BODY));
    },
  },
} as unknown as Env;

const get = (path: string): Promise<Response> =>
  worker.fetch(new Request(`https://manael.org${path}`), env, {} as ExecutionContext);

const metaContent = (html: string, name: string): string | undefined =>
  html.match(new RegExp(`<meta name="${name}" content="([^"]*)" />`, "u"))?.[1];

const goImport = (modulePrefix: string, repository: string): string =>
  `${modulePrefix} git https://github.com/${repository}.git`;

const goSource = (modulePrefix: string, repository: string): string => {
  const repositoryUrl = `https://github.com/${repository}`;
  return `${modulePrefix} ${repositoryUrl} ${repositoryUrl}/tree/main{/dir} ${repositoryUrl}/blob/main{/dir}/{file}#L{line}`;
};

interface Expectation {
  path: string;
  modulePrefix: string;
  repository: string;
  packagePath: string;
}

const expectations: Expectation[] = [
  {
    modulePrefix: "manael.org/x/manael",
    packagePath: "manael.org/x/manael",
    path: "/x/manael",
    repository: "manaelproxy/manael",
  },
  {
    modulePrefix: "manael.org/x/manael",
    packagePath: "manael.org/x/manael/cmd/manael",
    path: "/x/manael/cmd/manael",
    repository: "manaelproxy/manael",
  },
  {
    modulePrefix: "manael.org/x/manael/v3",
    packagePath: "manael.org/x/manael/v3",
    path: "/x/manael/v3",
    repository: "manaelproxy/manael",
  },
  {
    modulePrefix: "manael.org/x/manael/v3",
    packagePath: "manael.org/x/manael/v3/internal/foo",
    path: "/x/manael/v3/internal/foo",
    repository: "manaelproxy/manael",
  },
  {
    modulePrefix: "manael.org/x/image",
    packagePath: "manael.org/x/image",
    path: "/x/image",
    repository: "manaelproxy/image",
  },
  {
    modulePrefix: "manael.org/x/image",
    packagePath: "manael.org/x/image/webp",
    path: "/x/image/webp",
    repository: "manaelproxy/image",
  },
  {
    modulePrefix: "manael.org/x/image",
    packagePath: "manael.org/x/image/avif",
    path: "/x/image/avif",
    repository: "manaelproxy/image",
  },
  {
    modulePrefix: "manael.org/x/image/v2",
    packagePath: "manael.org/x/image/v2",
    path: "/x/image/v2",
    repository: "manaelproxy/image",
  },
  {
    modulePrefix: "manael.org/x/image/v2",
    packagePath: "manael.org/x/image/v2/webp",
    path: "/x/image/v2/webp",
    repository: "manaelproxy/image",
  },
  {
    // "v0" and "v01" are not major version suffixes, so they are packages
    // of the unversioned module.
    modulePrefix: "manael.org/x/image",
    packagePath: "manael.org/x/image/v0",
    path: "/x/image/v0",
    repository: "manaelproxy/image",
  },
  {
    modulePrefix: "manael.org/x/image",
    packagePath: "manael.org/x/image/v01/webp",
    path: "/x/image/v01/webp",
    repository: "manaelproxy/image",
  },
];

describe("vanity import paths", () => {
  for (const { path, modulePrefix, repository, packagePath } of expectations) {
    for (const query of ["", "?go-get=1"]) {
      it(`serves ${path}${query}`, async () => {
        const response = await get(`${path}${query}`);
        const html = await response.text();
        const packageUrl = `https://pkg.go.dev/${packagePath}`;

        assert.equal(response.status, 200);
        assert.equal(response.headers.get("Content-Type"), "text/html; charset=utf-8");
        assert.equal(metaContent(html, "go-import"), goImport(modulePrefix, repository));
        assert.equal(metaContent(html, "go-source"), goSource(modulePrefix, repository));
        assert.ok(html.includes(`<meta http-equiv="refresh" content="0; url=${packageUrl}" />`));
        assert.ok(html.includes(`<a href="${packageUrl}">${packageUrl}</a>`));
        assert.ok(html.includes(`<title>${packagePath}</title>`));
      });
    }
  }

  it("escapes HTML in the requested path", async () => {
    const response = await get("/x/image/a'b&c");
    const html = await response.text();

    assert.ok(html.includes("<title>manael.org/x/image/a&#39;b&amp;c</title>"));
    assert.ok(html.includes("url=https://pkg.go.dev/manael.org/x/image/a&#39;b&amp;c"));
    assert.ok(!html.includes("a'b&c"));
  });
});

describe("unrelated paths", () => {
  for (const path of [
    "/",
    "/docs/",
    "/x",
    "/x/",
    "/x/manaelfoo",
    "/x/imagefoo",
    "/x/images/webp",
    "/y/x/image",
    "/X/image",
  ]) {
    it(`falls back to static assets for ${path}`, async () => {
      assetRequests.length = 0;
      const response = await get(path);

      assert.equal(await response.text(), ASSET_BODY);
      assert.deepEqual(assetRequests, [path]);
    });
  }
});
