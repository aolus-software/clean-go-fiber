# GitHub Copilot Instructions — clean-go-fiber

> You are working on a Go REST API boilerplate. Before writing any code, read the agent documentation below.

## Start Here

1. Read `.agent/AGENTS.md` — project overview, architecture, file map, task patterns
2. Read `.agent/rules/general.md` — coding rules you must always follow
3. Read `.agent/skills/**/SKILL.md-fiber.md` — skill reference for all tech in this stack

## Your Responsibilities

- Follow clean architecture strictly: handler → service → repository → model
- Never skip layers, never put logic in handlers, never access DB outside repositories
- Always validate input, always type errors, always use the response helpers
- Every new handler must have complete swaggo annotations
- Run `make docs` after changing handlers

## When in Doubt

Check `.agent/AGENTS.md` section 7 (Common Task Patterns) and section 8 (File Reference Map) before creating any new file.
