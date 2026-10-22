<script lang="ts">
  import { live, run, service, levelLabel } from "./lib/status.svelte";

  const v = $derived(live.view);
  const level = $derived(v?.level ?? "stopped");
  const busy = $derived(v?.busy ?? "");
</script>

<main>
  <div class="head">
    <span class="dot {level}"></span>
    <strong>Docker Zero</strong>
    <span class="state">{busy ? "Working..." : levelLabel[level]}</span>
    <button class="close" title="Hide corner box" onclick={() => service.SetCornerBox(false)}>&#10005;</button>
  </div>
  <div class="actions">
    {#if v?.portainerUp}
      <button onclick={() => run(() => service.OpenPortainer())}>Open Portainer</button>
      <button disabled={!!busy} onclick={() => run(() => service.StopDistro())}>Stop</button>
    {:else}
      <button class="primary" disabled={!!busy} onclick={() => run(() => service.Start())}>Start</button>
    {/if}
  </div>
  {#if live.error}
    <p class="error" title={live.error}>{live.error}</p>
  {/if}
</main>

<style>
  main {
    height: 100%;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    border: 1px solid var(--border);
    border-radius: var(--radius);
    --wails-draggable: drag;
  }
  .head {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .state {
    margin-left: auto;
    color: var(--text-muted);
    font-size: 12px;
  }
  .close {
    border: 0;
    background: none;
    cursor: pointer;
    color: var(--text-muted);
    --wails-draggable: no-drag;
  }
  .actions {
    display: flex;
    gap: 8px;
  }
  .actions button {
    flex: 1;
    padding: 6px 8px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--card);
    cursor: pointer;
    --wails-draggable: no-drag;
  }
  .actions button:disabled {
    opacity: 0.45;
    cursor: default;
  }
  .actions button.primary {
    background: var(--accent);
    border-color: var(--accent);
    color: var(--accent-text);
  }
  .error {
    margin: 0;
    font-size: 11px;
    color: var(--bad);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
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
</style>
