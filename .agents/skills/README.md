# Skills

Engineering skills available to Claude Code in this repo.

`.claude/skills` is a **symlink** to this directory (`.agents/skills/`), so the same bundle is
visible to any agent that reads either path. Do not replace the symlink with a real directory — the
two would drift.

Skills here are **vendored** — the files in this directory are the source of truth. There is no
provenance lock: `skills-lock.json` is not tracked (it is in `.gitignore`). The universal set was
copied from the sibling `clean-elysia` bundle; `go-fiber` and `golang-pro` were carried over from
this repo's former `.agent/skills/`.

> Because nothing records each skill's upstream source or content hash, **drift is only detectable by
> manual diff.** If you suspect a skill is out of date, compare it against upstream directly.

## Installed skills

### Workflow and process

| Skill | Use it for |
| --- | --- |
| `brainstorming` | Explore intent, requirements, and design before any implementation work |
| `grill-me` / `grill-with-docs` | Stress-test a plan, optionally against the repo's docs |
| `writing-plans` | Turn a feature request into a reviewable, step-by-step plan |
| `executing-plans` | Work through an approved plan one step at a time |
| `test-driven-development` | Write the failing `_test.go` first, implement to green, then refactor |
| `systematic-debugging` | Reproduce, isolate, hypothesise, verify — instead of guessing |
| `requesting-code-review` / `receiving-code-review` | Ask for and act on a review of a diff |
| `verification-before-completion` | Confirm the work actually runs before calling it done |
| `finishing-a-development-branch` | Wrap up, clean up, and land a branch |
| `using-git-worktrees` | Run parallel work trees without stepping on the main checkout |
| `dispatching-parallel-agents` / `subagent-driven-development` | Fan work out to subagents, then integrate |
| `handoff` | Compact the current session into a handoff document for the next agent |
| `writing-skills` | Author a new skill correctly |
| `using-superpowers` / `find-skills` | Discover and pull in skills that are not installed yet |

### Engineering craft

| Skill | Use it for |
| --- | --- |
| `go-fiber` | Fiber v2 routing, middleware, context, error handler |
| `golang-pro` | Idiomatic Go — interfaces, errors, concurrency, generics |
| `clean-code` | Turn "code that works" into code that reads well — naming, functions, structure |
| `postgres` | PostgreSQL best practices, query optimisation, connection troubleshooting |
| `postgresql-optimization` | PostgreSQL-specific features — JSONB, arrays, ranges, full-text search, extensions |
| `caveman` | Ultra-compressed output mode when token efficiency is explicitly requested |

The TypeScript-only skills of the sibling repos (`typescript-pro`, `typescript-advanced-types`) and
framework skills for other stacks are deliberately not vendored here.

## Relationship to `.claude/rules/`

Skills are **general** engineering technique. The standards that are specific to *this* codebase —
layering, response envelope, RBAC gating, validation, repository allow-lists — live in
`.claude/rules/` and take precedence wherever the two overlap. Read the relevant rule before writing
code; reach for a skill for the technique.

## Adding a skill

Create `<skill-name>/SKILL.md` here and commit it — the vendored copy is what the repo keeps:

```markdown
---
name: <skill-name>
description: <one line — when this skill should trigger>
---

# <Skill Name>

<instructions / checklist / examples>
```
