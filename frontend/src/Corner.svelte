<script lang="ts">
  import { live, run, service } from "./lib/status.svelte";
  import Ring from "./Ring.svelte";
  import ErrorNote from "./ErrorNote.svelte";

  let errOpen = $state(false);
  const v = $derived(live.view);
  const err = $derived(live.error || v?.error || "");
  const level = $derived(v?.level ?? "stopped");
  const busy = $derived(v?.busy ?? "");
  const up = $derived([v?.wslAvailable, v?.distro === "running", v?.docker === "active", v?.portainerUp].filter(Boolean).length);
  const word = $derived(
    busy ? "Working" : ({ running: "Running", partial: "Not ready", stopped: "Stopped", failed: "Needs setup" }[level] ?? "Stopped"),
  );
</script>

<main>
  <Ring id="corner" orca parts={up} {level} size={64} />
  <div class="col">
    <div class="head">
      <strong class="word">{word}</strong>
      <button class="close" title="Hide corner box" aria-label="Hide corner box" onclick={() => service.SetCornerBox(false)}
        >&#10005;</button
      >
    </div>
    {#if err}
      <ErrorNote compact showDetail={false} bind:open={errOpen} detail={err} />
    {/if}
    {#if errOpen}
      <p class="detail">{err}</p>
    {:else}
      <div class="actions">
        {#if v?.portainerUp}
          <button class="primary" onclick={() => run(() => service.OpenPortainer())}>Open</button>
          <button disabled={!!busy} onclick={() => run(() => service.StopDistro())}>Stop</button>
        {:else}
          <button class="primary" disabled={!!busy} onclick={() => run(() => service.Start())}>Start</button>
        {/if}
      </div>
    {/if}
  </div>
</main>

<style>
  main {
    display: flex;
    align-items: center;
    gap: 12px;
    height: 100%;
    padding: 8px 12px;
    background: var(--sea-0);
    border: 1px solid var(--line);
    border-radius: var(--radius);
    --wails-draggable: drag;
  }
  .col {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 5px;
  }
  .head {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .word {
    font: 15px/1 var(--font-display);
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
    gap: 6px;
  }
  .actions button {
    flex: 1;
    min-height: 28px;
    border: 1px solid var(--line);
    border-radius: 9px;
    background: var(--card);
    font-weight: 600;
    cursor: pointer;
    --wails-draggable: no-drag;
  }
  .actions button.primary {
    background: var(--teal);
    border-color: var(--teal);
    color: var(--on-teal);
  }
  .actions button:disabled {
    opacity: 0.5;
    cursor: default;
  }
  .detail {
    margin: 0;
    max-height: 38px;
    overflow: auto;
    color: var(--text-muted);
    font: 11px/1.35 Consolas, monospace;
    word-break: break-word;
    user-select: text;
  }
</style>
