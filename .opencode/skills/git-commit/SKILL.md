---
name: git-commit
description: Read git status and create well-organized commits, splitting into separate commits when changes are logically distinct or combining into a single commit when appropriate
---

## What I do

Analyze the current git working tree and create well-structured commits. I determine whether changes should be grouped into a single commit or split into multiple logically separated commits.

## When to use me

Use this skill when the user asks to commit their changes, prepare commits, or review and commit what they've done.

## Steps

1. **Gather context** - Run these commands in parallel:
   - `git status` to see all modified, deleted, and untracked files
   - `git diff` to see unstaged changes
   - `git diff --cached` to see already-staged changes
   - `git log --oneline -10` to understand the recent commit message style

2. **Read and understand all changes** - Read every modified, deleted, and new file to fully understand what was changed and why. Do NOT skip any files. You need full context to make good grouping decisions.

3. **Check for files that should NOT be committed**:
   - Environment files containing secrets (`.env`, `.env.local`, `.env.production`, etc.)
   - Credentials, API keys, tokens
   - If such files are found, warn the user and suggest adding them to `.gitignore` instead

4. **Ask for help if unclear** - If you cannot determine the intent or logical grouping of the changes, ask the user for clarification before proceeding. Do not guess when uncertain.

5. **Decide: single commit or multiple commits** - Apply these rules:
   - **Single commit** when all changes serve one logical purpose (e.g., a single feature, a single refactor, a single bug fix)
   - **Multiple commits** when changes are logically distinct and independent. Common split boundaries:
     - Infrastructure/tooling changes vs application code changes
     - Refactoring existing code vs adding new features
     - Test changes that go with a specific feature vs unrelated test fixes
     - Documentation updates vs code changes
     - Dependency updates vs code that uses those dependencies (unless tightly coupled)

6. **Present the plan to the user** - Before committing, show:
   - How many commits will be created
   - For each commit: the files included and the proposed commit message
   - Ask the user to confirm or adjust

7. **Create the commits** - For each commit:
   - Stage only the relevant files using `git add <specific files>`
   - Write a commit message following conventional commits format (`feat:`, `fix:`, `refactor:`, `docs:`, `chore:`, `test:`, etc.)
   - The first line should be concise (under 72 chars)
   - Add a blank line and a body paragraph explaining the "why" when the change is non-trivial
   - Run `git status` after the final commit to verify a clean working tree

## Commit message guidelines

- ALWAYS use Conventional Commits format: `type: short description`
- Types: `feat`, `fix`, `refactor`, `docs`, `chore`, `test`, `build`, `ci`, `style`, `perf`
- Optional scope: `type(scope): short description`
- Body should explain motivation, not repeat the diff
- Do NOT include file lists in the commit message body

## Important rules

- NEVER commit files that contain secrets
- NEVER use `git add .` or `git add -A` when making multiple commits - always stage specific files
- NEVER force push or amend existing commits unless explicitly asked
- ALWAYS ask the user for confirmation before creating the commits
- ALWAYS verify with `git status` after the last commit
