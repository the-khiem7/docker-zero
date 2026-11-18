---
baseline_schema: "3.0"
pack: "ui-ux"
document: "introduction"
status: "active"
updated: "2026-11-19"
code_ref: "uncommitted"
---

# ui-ux - the application UI and UX design of Docker Zero

## Scope

Design and apply the user interface and the UX of the Docker Zero app (tray popup, corner box, Settings with Maintenance, Setup window): inventory of today's UI, a redesign that follows the brand instead of imitating Windows 11, review on design boards, then the change in `frontend/src/`. Pack `docker-zero` owns the product, the Go code, the brand logo, the README and the landing page; this pack took over Phase 10 items 2 to 4 of that pack.

In scope: window layouts, controls, icons, motion, error and progress feedback, light and dark themes, the boards that show them.

Out of scope: Go behaviour (real progress for the Ubuntu install step stays out, UX-D10), the deferred tech debt rows (no way back to setup, a missing wrapper not noticed), the README and the landing page (they follow this design afterwards).

## Current truth

| Item | State |
|---|---|
| Boards | `design/ui-gallery/V1` (current UI), `V2-popup`, `V2-setting`; `V2-setup` to build (UX-D8) |
| Redesign source | `frontend/redesign/`, never imported by the production build; `frontend/src/` is the shipped app and is unchanged |
| Chosen design | Popup: ring centred with the logo's orca (UX-D6). Settings: control deck, a bento grid of visual controls (UX-D7) |
| Toolchain | Boards render the real Svelte views with the Go services mocked (`frontend/gallery/`), served on port 9247 by `npm run gallery2`; `shoot.py` from the `frontdoor-ui-gallery` skill |

## Target

- One visual language across every window: the brand's ring and two-tone sea, drawn icons, chips and bento controls instead of text symbols, paragraphs and stacked boxes.
- The user approves each board before any app code changes (UX-D3).

## Constraints

| Constraint | Reason |
|---|---|
| UI only, no Go change in the redesign | Some fixes need Go work and the verified setup flow stays untouched (UX-D3, UX-D4) |
| `frontend/src/` changes only after the user approves the boards | UX-D3 |
| Galleries are for approval; layout experiments go on the frontdoor board | UX-D8 |
| No invented progress values | A made-up bar teaches users to distrust the app (UX-D4) |
| Public repository: no usernames, machine names or home paths in files or images | Audience is anyone (DZ-D1 in pack `docker-zero`) |
| Standard-user rules and the pinned Wails version of `AGENTS.md` | Unchanged by the redesign |
