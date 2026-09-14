# mwanachama-backend-insights — open tasks

| Task | Title | Notes |
|------|-------|-------|
| I7 | Wire an automatic trigger (a Stop hook or slash command) so a session calls `insight_create` at the end, instead of relying on manually calling the tool | Deliberately deferred at I1-I6 so the tool could be exercised manually first — see `todo_done.md`'s board context. **Still needs fine-tuning before implementing**: which mechanism (Stop hook vs. slash command vs. something else), and how it decides what counts as "the session's insight" — not yet designed. |
