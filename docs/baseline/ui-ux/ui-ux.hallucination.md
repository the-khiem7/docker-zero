---
baseline_schema: "3.0"
pack: "ui-ux"
document: "hallucination"
status: "active"
updated: "2026-11-19"
code_ref: "uncommitted"
---

# ui-ux hallucination journal

## Entry index

| ID | About | Status | Related |
|---|---|---|---|
| UX-D1 | UI inventory: mocked frontend for every state, real exe for the main windows, committed | closed | DZ-Q4, DZ-D27 |
| UX-D2 | App UI redesign direction: brand-first look, ring gauge popup, purposeful motion | closed | DZ-D31, DZ-D35, DZ-D4, DZ-D27 |
| UX-D3 | Redesign is built and approved on a design board in `design/ui-gallery/V2-popup` before the real app; UI only | closed, board folder layout amended by UX-D8 | UX-D2, UX-D1, UX-D9 |
| UX-D4 | UX research joins the redesign: feedback for long steps, no invented progress; Go work for real progress stays out of this round | closed | UX-D3, UX-D2, UX-Q2 |
| UX-D5 | Dela Gothic One is bundled from the fontsource package into the redesign tree, with its licence, no new dependency | closed | UX-D2, UX-D3 |
| UX-D6 | Popup direction B (ring centred) chosen, with the orca of the logo inside the ring; setup window long-step feedback approved | closed | UX-D2, UX-D4 |
| UX-D7 | Settings window: control deck (bento) chosen; same language for the whole lifecycle | closed | UX-D6, UX-D8 |
| UX-D8 | Boards: approval galleries V1, V2-popup, V2-setting, V2-setup; experiments on the frontdoor board | closed | UX-D3, UX-D7 |
| UX-D9 | V2 approved and applied: redesign source moved from `frontend/redesign/` into `frontend/src/`, the parallel tree removed after the built exe was checked, no commit | closed | UX-D3, UX-D7, UX-D8, UX-Q1 |
| UX-D10 | Ubuntu install step keeps UI-only feedback (elapsed time, indeterminate cue, expected duration); no Go change for real progress | closed | UX-D4, UX-D3, UX-Q2, DZ-D19 |

## Open questions

None open.

## Closed decisions

### UX-D1: UI inventory: mocked frontend for every state, real exe for the main windows, committed

Decision: Phase 10 item 2 renders every window and state from the Svelte frontend in a browser with the Wails bindings replaced by per-state mocks, in light and dark; the built exe, driven by UI Automation, supplies reference screenshots of the main windows; the inventory page and its images are committed to the repository.

Closes DZ-Q4.

Why: several states exist only under real conditions (distro not installed, setup running, a failed pull, a stopped distro); mocks make each one repeatable for the before and after compare of item 4, and the real screenshots catch where the mock and WebView2 differ. Committed at the user's choice.

What breaks if ignored: real-exe-only capture needs an uninstall or a blocked Docker Hub per state and per redesign; mock-only capture hides a drift between the mock and the Go services.

Rejected: real exe only, and mocks only - the costs above. The inventory as a gitignored private Artifact, or a local file only - the user chose to commit it. Committed images must carry no account path, user name or machine name (public repository).

### UX-D2: App UI redesign direction: brand-first look, ring gauge popup, purposeful motion

Decision (user, 2026-11-17): the app stops imitating Windows 11 Action Center and takes the Docker Zero brand:

- Look: cream surface, teal `#1e7468` accent, the logo's ring as the state indicator in the tray colours (grey stopped, amber partial, green running, red failure, DZ-D31), a deep-sea dark theme instead of neutral grey; Dela Gothic One only for the app name and the state word, system UI for body text.
- Popup: one large zero ring whose water level rises with the four parts (WSL, distro, Docker Engine, Portainer) replaces the dot, pill and status table; one contextual primary button (Set up, Start or Open Portainer by state); Stop distro is a secondary action; Shut down WSL moves behind a "more" menu.
- Motion: the water rises with the real Start progress; one orca leap when Running is reached, never looped; `prefers-reduced-motion` respected as on the landing page (DZ-D35); CSS and SVG only.

Why: the user wants a UI unique to the app, not a WinUI copy [user]; the logo, README and landing page already carry the brand, so a WinUI-styled app reads as another product [agent]; most popup opens do one thing, and Shut down WSL stops every distro of the account, so it should not sit level with routine actions [agent]; a resident tray app must stay idle on CPU (DZ-D4) [agent].

What breaks if ignored: the app and its front door look like two products; a destructive account-wide action one click from Start; looping animation costing CPU in a resident app.

Rejected: brand colours and title font on the current Action Center layout - still reads as a recoloured WinUI panel [agent]. A neutral minimal style without the brand - unique but detached from the logo and landing page [agent]. Keeping the 2x2 grid of four equal tiles - gives Shut down WSL the weight of Start [agent]. A static UI with colour changes only - loses the brand's character [agent]. Ambient waves and random leaps while the popup is open - CPU cost in a resident app [agent].

### UX-D3: The redesign is built and approved on a design board in `design/ui-gallery/V2-popup` before the real app; UI only

Decision (user, 2026-11-17): Phase 10 item 3 builds the redesigned UI as a design board in `design/ui-gallery/V2-popup/` with the `frontdoor-ui-gallery` workflow (real Svelte views rendered with the Go services mocked, every state, light and dark), reviewed by the user there; item 4 applies it to the real app only after the user approves the board. The redesign changes the UI only: the tech debt rows in the roadmap's Phase 10 table (no way back to setup, missing wrapper not noticed, the corner box error slot, setup form after setup) are deferred to a later round.

Why: the user wants to review and change the design before any app code changes [user]; the board renders every state on demand, so the before board `design/ui-gallery/` and the new one pair state by state [agent]; keeping debts out keeps item 4's no-behaviour-change rule intact, since some fixes need Go changes [agent].

What breaks if ignored: app code rewritten around a design the user has not seen; a redesign mixed with behaviour changes that the inventory compare cannot isolate.

Rejected: redesigning `frontend/src/` directly - the user asked for the board first [user]. Folding the UI-only debt fixes into item 4 with the Go fixes as separate commits - the user chose UI only [user]. Where the board's source lives: settled by UX-D9. Found on this workstation: the user PATH held `docker-zero\bin` while the folder and `docker.cmd` did not exist, so `docker` failed in every shell; the cause is unknown (see the Phase 10 tech debt table in the roadmap).

### UX-D4: UX research joins the redesign: feedback for long steps, no invented progress; Go work for real progress stays out of this round

Decision (user, 2026-11-17): the redesign covers UX as well as looks. The first finding is the setup window: the step "Install the Ubuntu distro" shows only a half-filled circle and one log line for several minutes. The board draws long steps from data the app has: steps done out of total as a real overall bar, elapsed time, an indeterminate activity cue for the step in flight, the expected duration ("several minutes, mostly downloads"), a "still working" notice past a threshold, and the log kept but collapsed. The board shows no percentage the backend cannot report. Real byte progress is UX-Q2; no Go change in this round (UX-D3 stays).

Why: a user cannot tell a long download from a hang when nothing moves [user]; a fake percentage that stalls reads worse than an honest indeterminate state [agent]; the setup flow was verified end to end (DZ-D19), so changing its install route needs its own decision [agent].

What breaks if ignored: the redesign repaints the same dead-looking wizard; a made-up progress bar teaches users to distrust the app.

Rejected: a percentage bar driven by a timer - invents data [agent]. Changing the install route to get byte progress now - touches the verified flow and breaks UI-only scope (UX-Q2 option b) [agent].

### UX-D5: Dela Gothic One is bundled from the fontsource package into the redesign tree, with its licence, no new dependency

Decision (user, 2026-11-17): the font files (woff2 and the SIL Open Font Licence text) are taken from the `@fontsource/dela-gothic-one` npm package, downloaded to a temporary folder, and only the woff2 and licence are copied into `frontend/redesign/`; `frontend/package.json` gains no dependency. Dela Gothic One is used for the app name and state word only (UX-D2).

Why: the app runs offline, so a Google Fonts link cannot ship [agent]; a copied file avoids a runtime dependency and a lock-file change [agent]; the licence file must travel with the font [agent].

What breaks if ignored: the title font fails to load offline and falls back to a system face; a font shipped without its licence.

Rejected: a stylesheet link to Google Fonts - network at run time. A package dependency in `frontend/package.json` - adds a dependency for two files. A system font only - drops the brand character of UX-D2.

### UX-D6: Popup direction B (ring centred) chosen, with the orca of the logo inside the ring; setup window long-step feedback approved

Decision (user, 2026-11-18): round 1 of the board picks direction B for the popup: one large centred ring with the state word below and the four parts as a row. The ring carries the logo's orca leaping out of the ring over the two-tone sea, so the popup shows the mark itself; the ring colour follows the tray colours and the sea level rises with the parts up. The setup window design of UX-D4 (step counter, elapsed time, indeterminate bar, expectation line, still-working notice, collapsed log) is approved.

Why: the orca in the ring makes the popup read as the app's logo, not a generic gauge [user]; B has the most room for the mark [agent].

What breaks if ignored: a popup ring without the orca reads as a generic status gauge unlike the logo, README and landing page.

Rejected: direction A (ring left) and C (sea band) - not chosen [user]. The plain ring without the orca (round 1 first render) - the user asked for the orca [user].

### UX-D7: Settings window design: the control deck, a bento grid of visual controls, and the same bento language for the whole app lifecycle

Decision (user, 2026-11-19): from three layouts explored on the frontdoor board `.frontdoor/product/` (A Harbour rail, B Control deck, C Depth strata) the user chose B: a bento grid where every setting is one large visual control - a power disc for auto start, app-icon tiles for shortcuts, round picks for the browser, a tide slider for the idle mode, a passport stamp for the certificate. The parts the board did not draw (the Portainer card, the reset-password form, the uninstall plan, confirmation, progress and failure, the working notice, setup) follow the same language: a disc or round button as the action, chips for facts, a step track for progress, long explanations behind an `i` button; no paragraphs and no stack of boxes. The parts the user approved are not redesigned again. The popup and the corner box (UX-D6) are kept. The bento language covers the whole lifecycle: setup, daily use, uninstall (setup goes to the board `V2-setup`).

Why: the user found the earlier Settings build still read as Windows 11 [user]; asked for drawn icons instead of text symbols and called walls of text a 90s look [user]; after B was applied only to the parts the board drew, asked for the rest to follow it, without replacing the approved parts and without paragraphs or a bundle of boxes [user]; the user wants the exploration on the frontdoor board and the gallery only for approval (UX-D8) [user].

What breaks if ignored: a window that mixes two languages (some controls bento, the rest boxed rows); paragraphs and stacked boxes bring back the wall of text the user rejected.

Rejected: A Harbour rail - closest to a standard OS settings page, the brand rests on the rail and icons [agent]. C Depth strata - mid-layer contrast is low and the dark theme breaks the depth order [agent]. The porthole bays with tiles (first V2-setting build) - the user wanted another layout to be explored [user]. Hold-to-uninstall - a new interaction not asked for; the plan and confirm steps stay [agent].

### UX-D8: Boards: approval galleries under `design/ui-gallery/`, layout experiments only on the frontdoor board

Decision (user, 2026-11-19): the approval galleries live together in `design/ui-gallery/`: `V1` (the current UI), `V2-popup` (popup and corner box), `V2-setting` (Settings and Maintenance) and `V2-setup` (the Setup window); each is rendered with the `frontdoor-ui-gallery` workflow. A gallery is for approval, never for experiments. Layout options are explored on the frontdoor design board `.frontdoor/product/` (untracked), one round per question, and the chosen option is then built into the gallery.

Why: the user said the gallery is for the approval approach, not for experiments [user]; a window group per gallery keeps each review short [user]; the frontdoor board already compares directions and records the choice [agent].

What breaks if ignored: rough sketches mixed into the board the user approves; one gallery too long to review.

Rejected: three sibling folders `design/ui-gallery`, `design/ui-gallery2`, `design/ui-gallery3` - renamed by the user [user]. The Settings states left in `V2-popup` - the user asked to remove them [user]. Exploring layouts inside the gallery - [user].

### UX-D9: V2 approved and applied into `frontend/src/`; the parallel tree removed after the exe check

Decision (user, 2026-11-19): the boards `V2-popup`, `V2-setting` and `V2-setup` are approved as a whole (`V1` is not part of the approval) and the design is applied to the app. Popup, Corner, Settings, Maintenance, Setup, Ring, Icon, Hint, ErrorNote and `sx.css` move into `frontend/src/`; the redesign tokens replace `frontend/public/style.css`; the font and its licence move to `frontend/public/fonts/`. `frontend/redesign/`, `frontend/vite.gallery2.config.ts`, the `gallery2` script and the `redesign` entries of `frontend/tsconfig.json` stayed until the user ran the built exe and reported that the new UI works (user, 2026-11-19), then the user approved removing them; the `regenerate`, `base` (now the `gallery` port 9246) and `source` fields of the three V2 boards point at `frontend/src/`. `V1` keeps its shots as history and a note in its lede, because rendering it again now shows V2. Nothing is committed.

Closes UX-Q1.

Why: UX-Q1 recommended a parallel tree so that `main` stays shippable and applying the design is a move into `frontend/src/`; the user approved the boards, chose to keep the tree until the exe is checked, then approved removing it [user]. Until then the tree was the only source the V2 boards regenerated from and the comparison point if the exe showed a defect [agent].

What breaks if ignored: deleting `frontend/redesign/` before the exe check would have left the boards without a regenerate path and removed the side-by-side reference; keeping it after the check leaves two copies of every component to drift.

Rejected: removing `frontend/redesign/` in the same pass as the move - the user chose to keep it until the exe was checked [user]. Committing now - the user said no [user].

### UX-D10: The Ubuntu install step keeps UI-only feedback; no Go change for real progress

Decision (user, 2026-11-19): the setup step "Install the Ubuntu distro" keeps the feedback of UX-D4 (elapsed time, an indeterminate cue, the expected duration, a "still working" notice, the log behind Details) and shows no percentage; no round changes Go for byte progress.

Closes UX-Q2.

Why: the step runs `wsl --install --distribution <image> --web-download --no-launch` (`internal/wsl/wsl.go`), which prints no percentage or byte count the app reads, and a step carries only `id`, `title`, `status` and `note` (`StepView` in `internal/ui/setup_service.go`); the honest state is "working, not stuck" [agent]. The user approved the recommendation after the built exe worked [user].

What breaks if ignored: a made-up bar teaches users to distrust the app; real progress through another route touches the verified setup flow.

Rejected: (b) download the rootfs inside the app and import it with `wsl --import` - gives byte progress but changes DZ-D19's install mechanics and the image source [agent]. (c) poll the size of the download cache while the command runs - fragile and undocumented [agent].
