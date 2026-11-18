<script lang="ts">
  import { live, run, service } from "./lib/status.svelte";
  import Ring from "./Ring.svelte";
  import ErrorNote from "./ErrorNote.svelte";

  interface Props {
    /** Board only: freeze the orca leap at this point, 0 to 1. */
    frame?: number;
  }
  let { frame }: Props = $props();

  let menu = $state(false);
  let confirmShutdown = $state(false);
  let errOpen = $state(false);

  const v = $derived(live.view);
  const busy = $derived(v?.busy ?? "");
  const level = $derived(v?.level ?? "stopped");
  const running = $derived(v?.distro === "running");
  const notInstalled = $derived(v?.distro === "not-installed");
  const noWsl = $derived(!!v && !v.wslAvailable);

  const parts = $derived([
    { label: "WSL", ok: !!v?.wslAvailable, text: v ? (v.wslAvailable ? "Available" : "Missing") : "..." },
    { label: "Ubuntu", ok: running, text: distroText(v?.distro) },
    { label: "Docker", ok: v?.docker === "active", text: v?.docker === "active" ? "Active" : (v?.docker ?? "...") },
    {
      label: "Portainer",
      ok: !!v?.portainerUp,
      text: v?.portainerUp ? v.portainerVersion : busy ? "Waiting" : portainerText(v?.portainer),
    },
  ]);
  // An error or the confirmation needs room in the 420 px window, so the ring shrinks.
  const compact = $derived(!!(live.error || v?.error) || confirmShutdown);
  const upCount = $derived(parts.filter((p) => p.ok).length);

  const word = $derived.by(() => {
    if (busy) return { start: "Starting", stop: "Stopping", shutdown: "Shutting down" }[busy] ?? "Working";
    if (noWsl) return "No WSL";
    if (notInstalled) return "Not set up";
    return { running: "Running", partial: "Not ready", stopped: "Stopped", failed: "Needs setup" }[level] ?? "Stopped";
  });

  const hint = $derived.by(() => {
    if (busy) return `${upCount} of 4 parts up`;
    if (noWsl) return "WSL is not enabled on this PC. Turning it on needs an administrator.";
    if (notInstalled) return "Ubuntu, Docker and Portainer are not installed yet.";
    if (level === "running") return `Portainer ${v?.portainerVersion ?? ""} is ready${v?.holding ? ", kept running" : ""}`;
    if (level === "partial") return "Ubuntu is up but Portainer is not answering.";
    return "Ubuntu is stopped.";
  });

  function portainerText(s?: string): string {
    return { running: "Not answering", exited: "Stopped", missing: "Not created", unknown: "Unknown" }[s ?? ""] ?? "...";
  }
  function distroText(d?: string): string {
    return { running: "Running", stopped: "Stopped", "not-installed": "Not installed" }[d ?? ""] ?? "...";
  }

  const primary = $derived.by(() => {
    if (busy) return { label: word + "...", act: () => {}, off: true };
    if (notInstalled) return { label: "Set up", act: openSetup, off: false };
    if (noWsl) return { label: "Start", act: start, off: true };
    if (level === "running") return { label: "Open Portainer", act: open, off: !v?.portainerUp };
    return { label: "Start", act: start, off: false };
  });

  async function start() {
    await run(() => service.Start());
    if (!live.error) await service.HidePopup();
  }
  async function openSetup() {
    menu = false;
    service.OpenSetup();
    await service.HidePopup();
  }
  async function open() {
    await run(() => service.OpenPortainer());
    if (!live.error) await service.HidePopup();
  }
  async function shutdown() {
    confirmShutdown = false;
    await run(() => service.ShutdownWSL());
  }
  function pick(fn: () => void) {
    menu = false;
    fn();
  }
</script>

<main>
  <header>
    <h1>Docker Zero</h1>
    <div class="menuwrap">
      <button class="more" aria-label="More" aria-expanded={menu} onclick={() => (menu = !menu)}>
        <svg viewBox="0 0 16 16" width="16" height="16" aria-hidden="true"
          ><circle cx="3" cy="8" r="1.4" /><circle cx="8" cy="8" r="1.4" /><circle cx="13" cy="8" r="1.4" /></svg
        >
      </button>
      {#if menu}
        <div class="menu" role="menu">
          <button role="menuitem" onclick={() => pick(() => service.OpenSettings())}>Settings</button>
          <label class="item">
            <input
              type="checkbox"
              checked={v?.corner ?? false}
              onchange={(e) => service.SetCornerBox(e.currentTarget.checked)}
            />
            Corner box
          </label>
          <button role="menuitem" onclick={() => pick(openSetup)}>Set up...</button>
          <hr />
          <button
            role="menuitem"
            class="danger"
            disabled={!!busy || !v?.wslAvailable}
            onclick={() => pick(() => (confirmShutdown = true))}>Shut down WSL</button
          >
        </div>
      {/if}
    </div>
  </header>

  <section class="hero">
    <Ring id="popup" orca {frame} parts={upCount} {level} size={compact ? 112 : 168} />
    <div class="state">
      <p class="word" style="--ring:var(--{level})">{word}</p>
      {#if !(live.error || v?.error)}<p class="hint">{hint}</p>{/if}
    </div>
  </section>

  {#if !errOpen}
  <ul class="parts" aria-label="Parts">
    {#each parts as p (p.label)}
      <li class:ok={p.ok}>
        <span class="mark" aria-hidden="true">{p.ok ? "✓" : ""}</span>
        <span class="plabel">{p.label}</span>
        <span class="ptext">{p.text}</span>
      </li>
    {/each}
  </ul>
  {/if}

  {#if live.error || v?.error}
    <ErrorNote bind:open={errOpen} detail={live.error || v?.error || ""} />
  {/if}

  {#if confirmShutdown}
    <section class="confirm">
      <p>Shut down WSL? This stops <strong>every</strong> WSL distro on this account, not only Ubuntu-24.04.</p>
      <div class="pair">
        <button class="btn danger-solid" onclick={shutdown}>Shut down</button>
        <button class="btn ghost" onclick={() => (confirmShutdown = false)}>Cancel</button>
      </div>
    </section>
  {:else}
    <section class="actions">
      <button class="btn primary" disabled={primary.off} onclick={primary.act}>{primary.label}</button>
      {#if !notInstalled}
        <button class="btn ghost" disabled={!!busy || !running} onclick={() => run(() => service.StopDistro())}
          >Stop distro</button
        >
      {/if}
    </section>
  {/if}
</main>

<style>
  main {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 12px;
    height: 100%;
    padding: 14px 16px 16px;
    background: var(--sea-0);
  }
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  h1 {
    margin: 0;
    font: 15px/1 var(--font-display);
    letter-spacing: 0.01em;
    color: var(--teal-strong);
  }
  .more {
    display: grid;
    place-items: center;
    width: 30px;
    height: 30px;
    border: 0;
    border-radius: 8px;
    background: transparent;
    color: var(--text-muted);
    cursor: pointer;
  }
  .more svg {
    fill: currentColor;
  }
  .more:hover,
  .more[aria-expanded="true"] {
    background: var(--sea-1);
    color: var(--text);
  }
  .menuwrap {
    position: relative;
  }
  .menu {
    position: absolute;
    right: 0;
    top: 34px;
    z-index: 5;
    min-width: 188px;
    padding: 6px;
    background: var(--card);
    border: 1px solid var(--line);
    border-radius: 12px;
    box-shadow: 0 8px 24px rgb(0 0 0 / 0.18);
  }
  .menu button,
  .menu .item {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 8px 10px;
    border: 0;
    border-radius: 8px;
    background: none;
    text-align: left;
    cursor: pointer;
  }
  .menu button:hover:not(:disabled),
  .menu .item:hover {
    background: var(--sea-1);
  }
  .menu button:disabled {
    opacity: 0.45;
    cursor: default;
  }
  .menu .danger {
    color: var(--failed);
  }
  .menu hr {
    margin: 4px 6px;
    border: 0;
    border-top: 1px solid var(--line);
  }

  .hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    text-align: center;
  }
  .word {
    margin: 0;
    font: 24px/1.1 var(--font-display);
    color: var(--text);
  }
  .word::after {
    content: "";
    display: inline-block;
    width: 9px;
    height: 9px;
    margin-left: 8px;
    border-radius: 50%;
    background: var(--ring);
    vertical-align: middle;
  }
  .hint {
    margin: 6px 0 0;
    max-width: 290px;
    color: var(--text-muted);
    font-size: 13px;
  }

  .parts {
    margin: 0;
    padding: 0;
    list-style: none;
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 4px;
  }
  .parts li {
    display: grid;
    grid-template-columns: 1fr;
    justify-items: center;
    padding: 6px 2px;
    text-align: center;
    background: var(--card);
    border: 1px solid var(--line);
    border-radius: 10px;
  }
  .mark {
    margin-bottom: 2px;
    align-self: center;
    width: 14px;
    height: 14px;
    border: 1.5px solid var(--text-muted);
    border-radius: 50%;
    font-size: 9px;
    line-height: 11px;
    text-align: center;
    color: var(--on-teal);
  }
  li.ok .mark {
    background: var(--running);
    border-color: var(--running);
  }
  .plabel {
    font-size: 12px;
    color: var(--text-muted);
  }
  .ptext {
    font-size: 11px;
    font-weight: 500;
    color: var(--text-muted);
  }


  .actions,
  .pair {
    display: flex;
    gap: 8px;
    margin-top: auto;
  }
  .btn {
    min-height: 44px;
    padding: 0 18px;
    border: 1px solid var(--line);
    border-radius: 12px;
    background: var(--card);
    font-weight: 600;
    cursor: pointer;
  }
  .btn:disabled {
    opacity: 0.5;
    cursor: default;
  }
  .btn.primary {
    flex: 1;
    background: var(--teal);
    border-color: var(--teal);
    color: var(--on-teal);
    font-size: 15px;
  }
  .btn.primary:hover:not(:disabled) {
    background: var(--teal-strong);
  }
  .btn.ghost:hover:not(:disabled) {
    background: var(--sea-1);
  }
  .btn.danger-solid {
    flex: 1;
    background: var(--failed);
    border-color: var(--failed);
    color: #fff;
  }
  .pair .ghost {
    flex: 1;
  }
  .confirm {
    margin-top: auto;
    padding: 12px;
    background: var(--card);
    border: 1px solid var(--line);
    border-radius: var(--radius);
  }
  .confirm p {
    margin: 0 0 10px;
  }
</style>
