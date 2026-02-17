---
name: todo-add
description: Add tasks, study topics, or notes to the project TODO.md file, categorizing them appropriately based on the user's input
---

## What I do

Add new entries to the `TODO.md` file at the project root. Entries fall into two
categories:

1. **Project tasks** -- bugs to fix, features to build, refactors, infra work,
   tests to write, etc.
2. **Study topics** -- backend concepts the user wants to learn or deepen,
   either as general backend knowledge or as something to explore hands-on in
   Go.

The user is a mid-level fullstack developer transitioning to backend-focused
work, so study topics should be scoped to backend engineering (not frontend) and
should target mid-level knowledge gaps (not beginner basics, not staff-level
distributed-systems theory).

## When to use me

Use this skill when the user asks to:

- Add something to the TODO list / backlog
- Track a task, idea, or improvement for the project
- Note a concept or topic they want to study
- Any variation of "remind me to...", "add to TODO...", "I should learn about..."

## Steps

1. **Read the current TODO.md** -- Read `TODO.md` at the project root to
   understand the existing sections, formatting, and what is already tracked.
   This prevents duplicates and ensures consistent style.

2. **Classify the entry** -- Determine whether the user's input is:
   - A **project task** (code change, bug, feature, refactor, infra, test, etc.)
   - A **study topic** (concept to learn, technology to explore, pattern to
     understand)
   - **Both** (e.g., "implement graceful shutdown" is a task AND a learning
     opportunity)

3. **Pick the right section** -- Place the entry in the most appropriate
   existing section of `TODO.md`. The current sections are:

   | Section                              | Use for                                          |
   | ------------------------------------ | ------------------------------------------------ |
   | Critical Bugs                        | Blocking / data-loss bugs                        |
   | Architecture                         | Design improvements, domain modeling              |
   | Database                             | Schema, indexes, constraints, migrations          |
   | Error Handling                        | Error handling and resilience improvements        |
   | Security                             | Input validation, secrets, TLS, auth              |
   | Production Readiness (for new main.go) | HTTP server hardening, observability            |
   | Tests                                | New tests, test improvements                      |
   | Build & Deployment                   | CI/CD, Docker, Makefile, dependencies             |
   | Research / Design Questions           | Open questions that need investigation            |
   | Feature Ideas                        | New feature proposals                             |
   | Study Backlog                        | Backend concepts and Go topics to learn           |
   | Completed                            | Done items (never add new items here)             |

   If no existing section fits, create a new section following the same
   formatting conventions (## heading, `- [ ]` checkboxes).

   For **study topics**, always use the **Study Backlog** section. If this
   section does not exist yet, create it right before the **Completed** section
   with this format:

   ```markdown
   ## Study Backlog

   > Backend concepts and Go-specific topics to study. Ordered roughly by
   > relevance to current project work.

   - [ ] Topic -- Brief description of what to learn and why it matters
   ```

4. **Format the entry** -- Follow these formatting rules:
   - Use `- [ ]` checkbox syntax (never `- [x]` for new items)
   - Keep entries concise -- one line when possible
   - For project tasks: describe _what_ to do, not _how_
   - For study topics: include the topic name and a short note on scope or
     relevance (e.g., `- [ ] Connection pooling -- understand how database/sql
     pool works, when to tune MaxOpenConns/MaxIdleConns`)
   - For study topics: add a very short description (1-2 sentences) about what
     the topic is and why it matters, as indented lines right after the main
     entry. This gives the user immediate context without needing to open links.
   - For study topics: add 2-4 reference links as indented lines below the
     description. Prefer official docs, Go blog posts, and well-known tutorial
     sites. Format each link as `- <url> (<short label>)`. Only include links
     you are confident are real and correct (official Go docs, pkg.go.dev,
     go.dev/blog, go.dev/doc, DigitalOcean tutorials, etc.). Do NOT fabricate
     or guess URLs.
   - Study topic example:
     ```markdown
     - [ ] Connection pooling -- understand how database/sql pool works, when to tune MaxOpenConns/MaxIdleConns
       Manages a pool of reusable DB connections to avoid the overhead of opening a new connection per query.
       - https://go.dev/doc/database/manage-connections (official Go docs on managing connections)
       - https://pkg.go.dev/database/sql (database/sql package reference)
     ```
   - Match the wording style of existing entries (imperative verbs for tasks:
     "Add...", "Fix...", "Implement...")

5. **Check for duplicates** -- Before adding, scan the file for entries that
   already cover the same topic. If an existing entry is close but not
   identical:
   - If the new request adds meaningful detail, update the existing entry
     instead of creating a duplicate
   - If it is truly the same thing, inform the user it already exists and show
     them the existing entry

6. **Apply the edit** -- Use the Edit tool to insert the new entry into the
   correct section. Place it at the end of the section's item list (before the
   next `##` heading or end of file), preserving the blank line between
   sections.

7. **Confirm to the user** -- After editing, tell the user:
   - What was added
   - Which section it was placed in
   - If it was a study topic, optionally suggest a brief learning path or
     resource angle (e.g., "You could start by looking at how `database/sql`
     handles this internally" or "The Go blog has a good post on this")

## Study topic guidelines

When adding study topics, apply these filters:

- **Backend-only**: No frontend topics (React, CSS, browser APIs, etc.)
- **Mid-level scope**: Skip basics (what is HTTP, what is a database) and
  overly advanced topics (consensus algorithms, CRDTs) unless the user
  explicitly asks
- **Practical**: Prefer topics the user can explore hands-on in this project or
  in Go generally
- **Examples of good study topics**:
  - Context propagation and cancellation patterns in Go
  - Database connection pooling (`database/sql` internals)
  - Graceful shutdown and signal handling
  - Structured logging with `log/slog`
  - Database transaction isolation levels
  - Index design and query planning in PostgreSQL
  - Rate limiting algorithms (token bucket, sliding window)
  - HTTP middleware patterns in Go
  - Error wrapping and sentinel errors in Go
  - Database migration strategies (zero-downtime migrations)
  - Observability: metrics, tracing, health checks
  - Testing patterns: table-driven tests, test fixtures, integration vs unit

## Important rules

- NEVER remove or modify existing entries unless updating a near-duplicate
- NEVER mark items as completed (`[x]`) -- only the user decides when something
  is done
- NEVER add entries to the **Completed** section
- ALWAYS read TODO.md before editing to avoid duplicates and stay consistent
- ALWAYS preserve the existing file structure and formatting
- If the user's request is ambiguous, ask for clarification before adding
