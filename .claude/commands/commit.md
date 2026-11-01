---
name: commit
description: Stage and commit all current changes following the repo's commit style
risk: safe
source: local
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git add:*), Bash(git commit:*), Bash(make lint:*), Bash(go build:*), Bash(go vet:*), Bash(gofmt:*)
---

Commit the current working-tree changes.

Run these in parallel first to gather context:

- `git status` (no `-uall`)
- `git diff` (staged + unstaged)
- `git log --oneline -20` to match the existing commit-message style (prefixes like `chore:`, `refactor:`, `fix:`, `remove:`, `feat:`)

Then run the gates. There is no pre-commit hook in this repo, so these are the only check before a commit:

- `make lint` (`gofmt -l` must print nothing, then `go vet ./...`)
- `go build ./...`

If a gate fails, stop and fix the cause before committing.

Then:

1. Review the diff and decide the right prefix based on intent (`feat`, `fix`, `chore`, `refactor`, `docs`, `remove`, `test`).
2. Stage files explicitly by path — **never** `git add -A` / `git add .`. Skip `.env`, credentials, `tmp/`, compiled binaries, and any other untracked files that look sensitive or build-output.
3. Write a concise 1–2 sentence message focused on **why**, not a file list. Match the lowercased, prefix-style of recent commits.
4. Create the commit with a HEREDOC body (no `--amend`).
5. Run `git status` after to confirm the working tree is clean.

Do not push unless I explicitly ask.
