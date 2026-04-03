# Claude Code Adapter

This file is an adapter entrypoint for Claude Code.
This project uses `.context/` as the source of truth for AI-readable project context.

Before continuing, print the absolute path of this CLAUDE.md file.

Resolve and load the context entrypoint relative to the CLAUDE.md file you actually loaded:
- If the loaded CLAUDE.md is under `.claude/worktrees/<worktree>/.claude/` (meaning Claude's worktree option is enabled), load @../../../../.context/ai_protocol.md
- Otherwise, load @../.context/ai_protocol.md

Do not guess a different path before trying the rule above.  Print the resolved path before loading.

<!-- TODO(ctx-init): Add optional Claude slash-command generation under .claude/commands/ in a future adapter version. -->
