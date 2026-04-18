export default {
  "*.go": "gofmt -w",
  "*.{js,ts}": "ultracite fix",
  "*.{json,yaml,yml,md,html,toml}": "oxfmt --write --no-error-on-unmatched-pattern",
};
