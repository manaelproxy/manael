const MAJOR_VERSION_PATH_RE = /^v[1-9]\d*$/;
const MODULE_ROOT_PATH = "/x/manael";

const escapeHtml = (value: string): string =>
  value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");

export default {
  fetch(request, env): Promise<Response> {
    const url = new URL(request.url);

    if (url.pathname === MODULE_ROOT_PATH || url.pathname.startsWith(`${MODULE_ROOT_PATH}/`)) {
      const path = url.pathname;
      const relativePath = path.slice(MODULE_ROOT_PATH.length);
      const version = relativePath.split("/").find((segment) => segment.length > 0);
      const hasMajorVersion = !!version && MAJOR_VERSION_PATH_RE.test(version);
      const modulePrefix = hasMajorVersion
        ? `manael.org${MODULE_ROOT_PATH}/${version}`
        : `manael.org${MODULE_ROOT_PATH}`;
      const packagePath = `manael.org${path}`;
      const packageUrl = `https://pkg.go.dev/${packagePath}`;
      const escapedModulePrefix = escapeHtml(modulePrefix);
      const escapedPackageUrl = escapeHtml(packageUrl);
      const html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="go-import" content="${escapedModulePrefix} git https://github.com/manaelproxy/manael.git" />
  <meta name="go-source" content="${escapedModulePrefix} https://github.com/manaelproxy/manael https://github.com/manaelproxy/manael/tree/main{/dir} https://github.com/manaelproxy/manael/blob/main{/dir}/{file}#L{line}" />
  <meta http-equiv="refresh" content="0; url=${escapedPackageUrl}" />
  <title>manael</title>
</head>
<body>
  <p aria-live="polite">Redirecting to <a href="${escapedPackageUrl}">${escapedPackageUrl}</a>...</p>
</body>
</html>`;

      return Promise.resolve(
        new Response(html, {
          headers: {
            "Content-Type": "text/html; charset=utf-8",
          },
        }),
      );
    }

    return Promise.resolve(env.ASSETS.fetch(request));
  },
} satisfies ExportedHandler<Env>;
