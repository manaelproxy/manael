export const onRequest = (context) => {
  const { request } = context;
  const url = new URL(request.url);

  // Example: url.pathname = "/x/manael/v3/cmd/manael"
  const path = url.pathname;

  const html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="go-import" content="manael.org/x/manael git https://github.com/manaelproxy/manael.git" />
  <meta name="go-source" content="manael.org/x/manael https://github.com/manaelproxy/manael https://github.com/manaelproxy/manael/tree/main{/dir} https://github.com/manaelproxy/manael/blob/main{/dir}/{file}#L{line}" />
  <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/manael.org${path}" />
  <title>manael</title>
</head>
<body>
  <p aria-live="polite">Redirecting to <a href="https://pkg.go.dev/manael.org${path}">pkg.go.dev/manael.org${path}</a>...</p>
</body>
</html>`;

  return new Response(html, {
    headers: { "Content-Type": "text/html; charset=utf-8" },
  });
};
