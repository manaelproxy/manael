# Manael Agent Guide

Repository-specific conventions for agents. This file is the source of truth for implementation and review. Claude Code reads it directly, as other coding agents do, so every line here is loaded in every session; keep it to norms and the reasons behind them. Do not add a `CLAUDE.md` or `CLAUDE.local.md`: while either exists, Claude Code reads that file instead of this one.

## Output language

Respond to the user in the language they wrote in, even though this guide and the rest of this repository's documentation are written in English. A Japanese message is answered in Japanese, an English one in English.

Judge the language from the user's own prose, not from what the prompt quotes: a log line, an error message, or a snippet of code pasted into an English question does not make the question Japanese. Respond in **English** when no user prose has settled the language, such as a run started from CI or a scheduled run.

Code, identifiers, commit messages, and quoted technical terms stay as-is in either direction; only the explanations, summaries, and questions to the user follow the user's language.

## Documentation, test names, and GitHub: English

Every Markdown document this repository owns (`README.md`, `CONTRIBUTING.md`, this file, the pull request and issue templates) is written in English, and so are Go test names (`t.Run` labels and `Test*` functions). There are no translated copies, because a translation nothing checks drifts from its original.

What is filed on GitHub is English too: Issue titles and bodies, pull request titles and bodies, and commit messages, whatever language the session is held in. They outlive the session and are read by contributors who never saw it.

## Issue references in documentation and code comments

Long-lived documentation and code comments state the rule or decision and the reason for it, and do not cite the Issue it came out of; provenance belongs in the commit message and the pull request body (`Closes #NNN`). A reference to an **open** Issue describing something not done yet may stay, like a TODO, and is removed when that Issue closes.

## Git commits and pull requests

- Commit subjects follow [Conventional Commits](https://www.conventionalcommits.org/); `commitlint` checks them in the `commit-msg` hook.
- Disclose AI assistance with a single trailer in the form `Assisted-by: <tool>:<model>`, for example `Assisted-by: Claude Code:claude-opus-5-5`. Do not add `Co-Authored-By:` lines for an agent, or a "Generated with" footer to a pull request; `.claude/settings.json` turns Claude Code's own attribution off for this reason.
- Fill in `.github/PULL_REQUEST_TEMPLATE.md` (Description / Related Issue / Checklist) and link the Issue the pull request closes.
- Pull requests are squash-merged and the merge commit message is the pull request body, so the per-commit trailers do not survive the merge. End the body with the same `Assisted-by:` trailer as its last paragraph. Do not write `@coderabbitai summary` in the body: `.coderabbit.yaml` turns the summary off, but CodeRabbit still writes one wherever that placeholder appears.
- Commit or push only when asked, and branch off `main` first.

### Verification before a push

Run the checks for what changed, from the repository root, and do not describe an unrun command as passing.

| Changed area                          | Command                                                       |
| ------------------------------------- | ------------------------------------------------------------- |
| Go code (`*.go`, `cmd/`, `internal/`) | `go build ./...`, `go test ./...`, `golangci-lint run`        |
| Image processing behaviour end to end | `go test -tags=integration -timeout 10m ./...` (needs Docker) |
| Markdown, JSON, YAML, TypeScript      | `pnpm format-check` (`pnpm format` to fix)                    |
| `website/`                            | `pnpm build`                                                  |

The integration tests need Docker; a pull request says plainly when they were skipped rather than implying they passed.

## Build requirements

Manael links against libvips through cgo (`github.com/h2non/bimg`), so building and testing need `libvips-dev` and `CGO_ENABLED=1`. The devcontainer installs it; outside it, install libvips before concluding that a build failure is caused by the change.
