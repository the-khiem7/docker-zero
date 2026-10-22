<script lang="ts">
  import { live, run, service, levelLabel } from "./lib/status.svelte";

  let confirmShutdown = $state(false);

  const v = $derived(live.view);
  const busy = $derived(v?.busy ?? "");
  const level = $derived(v?.level ?? "stopped");
  const running = $derived(v?.distro === "running");

  const rows = $derived([
    { label: "WSL", value: v ? (v.wslAvailable ? "Available" : "Not available") : "...", ok: v?.wslAvailable },
    { label: "Ubuntu-24.04", value: distroText(v?.distro), ok: running },
    { label: "Docker Engine", value: v?.docker ?? "...", ok: v?.docker === "active" },
    {
      label: "Portainer",
      value: v?.portainerUp ? `Up, ${v.portainerVersion}` : (v?.portainer ?? "..."),
      ok: v?.portainerUp,
    },
  ]);

  function distroText(d?: string): string {
    switch (d) {
      case "running":
        return "Running";
      case "stopped":
        return "Stopped";
      case "not-installed":
        return "Not installed";
      default:
        return "...";
    }
  }

  async function start() {
    await run(() => service.Start());
    if (!live.error) await service.HidePopup();
  }

  async function open() {
    await run(() => service.OpenPortainer());
    if (!live.error) await service.HidePopup();
  }

  async function shutdown() {
    confirmShutdown = false;
    await run(() => service.ShutdownWSL());
  }
</script>

<main>
  <header>
    <div class="title">
      <span class="dot {level}"></span>
      <h1>Docker Zero</h1>
    </div>
    <div class="title">
      <span class="pill {level}">{busy ? busyLabel(busy) : levelLabel[level]}</span>
      <button class="gear" title="Settings" aria-label="Settings" onclick={() => service.OpenSettings()}>&#9881;</button>
    </div>
  </header>

  <section class="card status">
    {#each rows as row (row.label)}
      <div class="row">
        <span class="label">{row.label}</span>
        <span class="value" class:ok={row.ok}>{row.value}</span>
      </div>
    {/each}
  </section>

  {#if confirmShutdown}
    <section class="card confirm">
      <p>Shut down WSL? This stops <strong>every</strong> WSL distro on this account, not only Ubuntu-24.04.</p>
      <div class="confirm-actions">
        <button class="tile danger" onclick={shutdown}>Shut down</button>
        <button class="tile" onclick={() => (confirmShutdown = false)}>Cancel</button>
      </div>
    </section>
  {:else}
    <section class="tiles">
      <button class="tile primary" disabled={!!busy} onclick={start}>
        <span class="icon" aria-hidden="true">&#9654;</span>Start
      </button>
      <button class="tile" disabled={!v?.portainerUp} onclick={open}>
        <span class="icon" aria-hidden="true">&#8599;</span>Open Portainer
      </button>
      <button class="tile" disabled={!!busy || !running} onclick={() => run(() => service.StopDistro())}>
        <span class="icon" aria-hidden="true">&#9632;</span>Stop distro
      </button>
      <button class="tile" disabled={!!busy || !v?.wslAvailable} onclick={() => (confirmShutdown = true)}>
        <span class="icon" aria-hidden="true">&#9211;</span>Shut down WSL
      </button>
    </section>
  {/if}

  {#if live.error || v?.error}
    <p class="error">{live.error || v?.error}</p>
  {/if}

  <footer>
    <label class="switch">
      <input
        type="checkbox"
        checked={v?.corner ?? false}
        onchange={(e) => service.SetCornerBox(e.currentTarget.checked)}
      />
      Corner box
    </label>
    <span class="muted">{v?.holding ? "Kept running" : ""}</span>
  </footer>
</main>

<script module lang="ts">
  function busyLabel(b: string): string {
    return { start: "Starting...", stop: "Stopping...", shutdown: "Shutting down..." }[b] ?? "Working...";
  }
</script>

<style>
  main {
    display: flex;
    flex-direction: column;
    gap: 12px;
    height: 100%;
    padding: 14px;
    border: 1px solid var(--border);
    border-radius: var(--radius);
  }
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .title {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  h1 {
    margin: 0;
    font-size: 15px;
    font-weight: 600;
  }
  .card {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 8px 12px;
  }
  .row {
    display: flex;
    justify-content: space-between;
    padding: 4px 0;
  }
  .label {
    color: var(--text-muted);
  }
  .value.ok {
    color: var(--ok);
  }
  .tiles {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 8px;
  }
  .tile {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 6px;
    min-height: 64px;
    padding: 10px 12px;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    cursor: pointer;
    text-align: left;
  }
  .tile:hover:not(:disabled) {
    background: var(--card-hover);
  }
  .tile:disabled {
    opacity: 0.45;
    cursor: default;
  }
  .tile.primary {
    background: var(--accent);
    border-color: var(--accent);
    color: var(--accent-text);
  }
  .tile.danger {
    background: var(--bad);
    border-color: var(--bad);
    color: #fff;
  }
  .icon {
    font-size: 16px;
  }
  .confirm p {
    margin: 4px 0 10px;
  }
  .confirm-actions {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 8px;
  }
  .confirm-actions .tile {
    min-height: 0;
    align-items: center;
  }
  .error {
    margin: 0;
    color: var(--bad);
    font-size: 12px;
    word-break: break-word;
  }
  footer {
    margin-top: auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 12px;
  }
  .switch {
    display: flex;
    gap: 6px;
    align-items: center;
  }
  .muted {
    color: var(--text-muted);
  }
  .dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--text-muted);
  }
  .dot.running {
    background: var(--ok);
  }
  .dot.partial {
    background: var(--warn);
  }
  .dot.failed {
    background: var(--bad);
  }
  .gear {
    border: 0;
    background: none;
    cursor: pointer;
    font-size: 16px;
    color: var(--text-muted);
    padding: 0 2px;
  }
  .pill {
    font-size: 12px;
    padding: 2px 10px;
    border-radius: 999px;
    border: 1px solid var(--border);
    background: var(--card);
  }
</style>
