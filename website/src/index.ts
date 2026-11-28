const HOST = "manael.org";
const MAJOR_VERSION_PATH_RE = /^v[1-9]\d*$/u;

interface ModuleRoot {
  // Path of the module root on this host, without a trailing slash.
  path: string;
  // GitHub repository that hosts the module, as "owner/name".
  repository: string;
}

// Every vanity import root this site serves. A request is answered by the
// root it equals or descends from at a slash boundary, so "/x/manael" and
// "/x/manael/foo" match "/x/manael" while "/x/manaelfoo" does not. Each path
// here also has to be listed in `run_worker_first` in wrangler.toml, or the
// request is served from the static assets before this Worker sees it.
const MODULE_ROOTS: readonly ModuleRoot[] = [
  { path: "/x/manael", repository: "manaelproxy/manael" },
  { path: "/x/image", repository: "manaelproxy/image" },
];

const escapeHtml = (value: string): string =>
  value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");

const findModuleRoot = (pathname: string): ModuleRoot | undefined =>
  MODULE_ROOTS.find((root) => pathname === root.path || pathname.startsWith(`${root.path}/`));

const renderGoImport = (root: ModuleRoot, pathname: string): string => {
  const relativePath = pathname.slice(root.path.length);
  const version = relativePath.startsWith("/") ? relativePath.slice(1).split("/", 1)[0] : "";
  const hasMajorVersion = !!version && MAJOR_VERSION_PATH_RE.test(version);
  const modulePrefix = hasMajorVersion ? `${HOST}${root.path}/${version}` : `${HOST}${root.path}`;
  const packagePath = `${HOST}${pathname}`;
  const packageUrl = `https://pkg.go.dev/${packagePath}`;
  const repositoryUrl = `https://github.com/${root.repository}`;
  const escapedModulePrefix = escapeHtml(modulePrefix);
  const escapedPackagePath = escapeHtml(packagePath);
  const escapedPackageUrl = escapeHtml(packageUrl);
  const escapedRepositoryUrl = escapeHtml(repositoryUrl);

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="go-import" content="${escapedModulePrefix} git ${escapedRepositoryUrl}.git" />
  <meta name="go-source" content="${escapedModulePrefix} ${escapedRepositoryUrl} ${escapedRepositoryUrl}/tree/main{/dir} ${escapedRepositoryUrl}/blob/main{/dir}/{file}#L{line}" />
  <meta http-equiv="refresh" content="0; url=${escapedPackageUrl}" />
  <title>${escapedPackagePath}</title>
</head>
<body>
  <p aria-live="polite">Redirecting to <a href="${escapedPackageUrl}">${escapedPackageUrl}</a>...</p>
</body>
</html>`;
};

export default {
  fetch(request, env): Promise<Response> {
    const url = new URL(request.url);
    const root = findModuleRoot(url.pathname);

    if (root) {
      return Promise.resolve(
        new Response(renderGoImport(root, url.pathname), {
          headers: {
            "Content-Type": "text/html; charset=utf-8",
          },
        }),
      );
    }

    return Promise.resolve(env.ASSETS.fetch(request));
  },
} satisfies ExportedHandler<Env>;
