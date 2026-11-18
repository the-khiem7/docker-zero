---
baseline_schema: "3.0"
pack: "ui-ux"
document: "roadmap"
status: "active"
updated: "2026-11-19"
code_ref: "uncommitted"
---

# ui-ux roadmap

| Phase | Depends on | Status | Acceptance |
|---|---|---|---|
| Phase 1 | - | active | The inventory page lists every window and state and the user confirms the list is complete |
| Phase 2 | Phase 1 | complete | The user approves the boards `V2-popup`, `V2-setting` and `V2-setup`; nothing under `frontend/src/` changes before that (UX-D3) |
| Phase 3 | Phase 2 | complete | The approved design is in `frontend/src/`; `wails3 task test`, `wails3 task lint` and `wails3 task build` pass; every inventory state re-shot from the built exe and matched to its approved mockup |

Origin: these phases were items 2, 3 and 4 of Phase 10 in pack `docker-zero`; their rows, evidence and tech debt moved here. The original rows:

| 2. UI inventory | One HTML page listing every window and state (popup, corner box, settings, maintenance, setup with the take-over card, confirmations, errors), each with a screenshot: mocked frontend for every state, real exe for the main windows, committed (UX-D1) | User confirms the list is complete |
| 3. UI evaluation and design | Redesign board `design/ui-gallery/V2-popup/` built with the `frontdoor-ui-gallery` workflow, every state in light and dark, paired with the before board `design/ui-gallery/`; direction UX-D2, board-first and UI only UX-D3, UX research included UX-D4; round 1 the popup in two or three ring directions plus the running setup window with long-step feedback, round 2 every view | User approves the board; nothing under `frontend/src/` changes before that |
| 4. Apply the design | Svelte and CSS changes under `frontend/`, the window sizes in `internal/ui` if a design needs them; no behaviour change | `wails3 task test`, `wails3 task lint`, `wails3 task build` pass; every inventory state re-shot from the built exe in light and dark and matched to its approved mockup |

## Phase 1: UI inventory

Status: built, awaiting the user's review. Board `design/ui-gallery/V1/` (UX-D1).

Item 2 evidence, awaiting the user's review. Since `3ed9a76` the inventory uses the `frontdoor-ui-gallery` layout: `design/ui-inventory/` became `design/ui-gallery/` (state list in `gallery.json`, 5 views, 50 states, 4 real references; `shoot.py` and `real.ps1` ship with the skill), `frontend/inventory/` became `frontend/gallery/`, `vite.inventory.config.ts` became `vite.gallery.config.ts` and the script `npm run gallery`; the paths in the table below are the names at the time:

| Check | Result |
|---|---|
| Mock harness | `frontend/vite.inventory.config.ts` (`npm run inventory`, port 9246) swaps the six bound service modules for `frontend/inventory/mock/*`; scenarios in `frontend/inventory/fixtures.ts`; each mock `satisfies` the generated binding's signatures, so a new or changed Go method fails `svelte-check` |
| Render | `design/ui-inventory/shoot.py` (Playwright on the installed Edge): 50 states over popup, corner box, settings, maintenance and setup, light and dark, WebP at 2x, 4.4 MB |
| Page | `design/ui-inventory/index.html`: filters by window, screenshot theme and review state, Keep / Revise / Discuss with a note per state kept in the browser, Copy review as Markdown; no horizontal scroll at 390 px; no page errors |
| Real references | `design/ui-inventory/real.ps1` captured the popup, corner box, settings and setup windows of the installed exe; account paths painted over |
| Observations on the page | 10; among them a failed Start in the corner box shows nothing (the error line is squeezed to zero height), and Setup opened after setup shows Existing setup found / Nothing left to replace above a form that would run setup again |
| `wails3 task lint`, `wails3 task test`, `wails3 task build` | 0 issues, 0 svelte-check errors; tests ok; build ok, the production bundle does not import the mocks |

Not verified: the mocks against the Go services at run time beyond the type check; the page in browsers other than Edge.

## Phase 2: UI evaluation and design

Status: active. Direction UX-D2, board-first and UI only UX-D3, UX research UX-D4, popup direction UX-D6, Settings control deck UX-D7, board layout UX-D8.

| Board | Content | State |
|---|---|---|
| `design/ui-gallery/V1/` | The current UI, 50 states, 4 real references | built |
| `design/ui-gallery/V2-popup/` | Popup and corner box, 25 states; the Settings and Setup states moved out | built, user review pending |
| `design/ui-gallery/V2-setting/` | Settings and Maintenance as the control deck, 27 states in light and dark | built, user review pending |
| `design/ui-gallery/V2-setup/` | The Setup window in the same bento language (hero card with the ring and chips, a track of eight steps, a key disc as the start button, take-over tiles), 11 states in light and dark | built, user review pending |
| `.frontdoor/product/` | Untracked exploration board; round 1 compared three Settings layouts, the user chose B (UX-D7) | round 1 done |

Evidence of the Settings and Setup rebuilds: `npm run check` 0 errors, `npm run build` ok with no redesign code in `frontend/dist`, `shoot.py check` 0 FAIL on `V2-popup`, `V2-setting` and `V2-setup`; some states were viewed, the rest only checked. Not verified: the mocks against the live backend, the redesign in the built exe, motion played in a real window.

Source: built in `frontend/redesign/` (Popup, Corner, Settings, Maintenance, Setup, Ring, Icon, Hint, ErrorNote, `sx.css`, tokens, font) and served by `frontend/vite.gallery2.config.ts`; moved into `frontend/src/` and the tree removed on 2026-11-19 (UX-D9); the boards now render from `frontend/src/` with `npm run gallery`.

UX research input for item 3 (user, 2026-11-17): the setup step "Install the Ubuntu distro" gives no progress for minutes (UX-D4, UX-Q2). Skills found with `npx skills find`, not installed; installing any is the user's call, and each is third-party text to read before use:

| Skill | Installs | Use for |
|---|---|---|
| `uxcel-lab/product-skills@ux-loaders-and-progress` | 123 | feedback for long-running steps |
| `davidortinau/maui-skills@ux-desktop` | 35 | desktop window and tray conventions |
| `vercel-labs/agent-skills@web-design-guidelines` | 868.9K | accessibility and interface review of the board |
| `nextlevelbuilder/ui-ux-pro-max-skill@ui-ux-pro-max` | 455K | UX heuristics and palette checks |

Item 3 evidence (round 2, built and viewed, not yet reviewed by the user):

| Check | Result |
|---|---|
| Board | `design/ui-gallery/V2-popup/`: 59 states in light and dark, ids paired with `design/ui-gallery/` plus motion frames (`popup-starting-1`, `popup-starting-2`, `popup-leap-0`, `popup-leap-50`, `popup-leap-100`), the long-step setup states (`setup-distro`, `setup-distro-slow`) and `popup-more`; `shots/` 4.9 MB |
| Source | `frontend/redesign/` (Popup, Corner, Settings, Maintenance, Setup, Ring, tokens, Dela Gothic One woff2 with its licence, UX-D5); served by `frontend/vite.gallery2.config.ts` on port 9247, script `gallery2`; Settings and Maintenance keep the current markup with brand styles |
| Checks | `npm run check` 0 errors, `npm run build` ok, `frontend/dist` holds no redesign code or font; `shoot.py check design/ui-gallery/V2-popup` 0 FAIL; no user name, machine name or home path in `gallery.json` or `frontend/redesign/` |
| Toolchain on this machine | Go `1.27.4`; the installed `wails3` is `v3.1.0`, not the pinned beta.30, so bindings were generated with beta.30 built into a temporary folder; `uv` absent, so `shoot.py` ran from a Python virtual environment in a temporary folder |

Defects found while viewing the board and fixed: Svelte renames scoped `@keyframes`, so the board's frozen leap frames showed no motion until the keyframes were made global; the popup pushed its buttons out of the 420 px window when the confirmation or an error showed, so the ring shrinks and the hint hides then; the Portainer part read `running` while Portainer was not answering.

Not verified: the mocks against the live backend; the redesign in the built exe; the leap and the water rising played in a real window (the board shows frames only); a 0 to 100 percent download bar (the backend reports none, UX-Q2).

Tech debt deferred to a round after item 4 (research first, then design); the redesign is UI only (UX-D3):

| Debt | Evidence | Research to do |
|---|---|---|
| The setup window has no visible way back once setup is done | Reachable only from the tray menu entry "Set up Docker Zero..." and the small "Set up..." link at the foot of the popup; the Settings window has no entry; the window opens by itself only while the distro is not installed (DZ-D24) | Which name and place suit a returning user (a "Repair setup" action in Settings, a banner in the popup when a part is missing); what the window should show on a finished setup (the "Nothing left to replace" observation of item 2) |
| Setup does not notice a missing or stale part | On this workstation `docker-zero\bin` was on the user PATH while the folder and `docker.cmd` were gone; every shell reported `docker` not found and the popup showed Running (DZ-D36) | Cause (uninstall, a manual delete, or a skipped step 8); whether the popup should check the wrapper files and PATH and offer Repair |

## Phase 3: Apply the design

Status: complete (user, 2026-11-19: closed on the exe check already done, the per-state re-shot from the built exe waived). The user approved `V2-popup`, `V2-setting` and `V2-setup` on 2026-11-19 (UX-D9). The components, `sx.css`, the tokens (as `public/style.css`) and the font (`public/fonts/`) are in `frontend/src/` and `frontend/public/`; `frontend/redesign/`, `frontend/vite.gallery2.config.ts` and the `gallery2` script are removed.

Evidence: `npm run check` 0 errors; `npm run build` passes and the font is copied to `dist/fonts/`; `wails3 task build` with the pinned `v3.0.0-beta.30` produced `bin/docker-zero.exe`; the user ran the exe and reported the new UI works (user, 2026-11-19); `npm run gallery` on port 9246 re-rendered 25 + 27 + 11 states of the three V2 boards from `frontend/src/` and `shoot.py check` gave 0 FAIL on each; `wails3 task test` and `wails3 task lint` pass.

Lint: `wails3 task lint` first stopped on `gofumpt` for `internal/adopt/adopt.go`, `internal/target/target.go` and `main.go`; the user chose to fix it (2026-11-19) with `golangci-lint fmt ./...`, which rewrote the working-tree line endings of the Go files from CRLF to LF (the repository index already holds LF, `git diff` shows no content change) and `wails3 task lint` now passes.

Not verified: the inventory states re-shot from the built exe and matched one by one to the mockups (waived by the user; the user confirmed the exe in general); motion played in a real window beyond what the user saw; the mocks against the live backend.

## Next action

1. Phase 3 is complete. Commit the work in groups once the user authorizes; `design/ui-gallery/` stays out of the commits and in the working tree until the user says otherwise.
2. The README and the landing page follow the new design (out of this pack's scope, see the introduction).
