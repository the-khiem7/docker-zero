<script lang="ts">
  interface Props {
    /** The raw error text; shown only behind the info button. */
    detail: string;
    /** A short human sentence; derived from the detail when left out. */
    title?: string;
    /** One line: badge, title and info button only. */
    compact?: boolean;
    open?: boolean;
    /** False when the parent shows the detail itself. */
    showDetail?: boolean;
  }

  let { detail, title, compact = false, open = $bindable(false), showDetail = true }: Props = $props();

  function friendly(msg: string, short: boolean): string {
    const m = msg.toLowerCase();
    if (short) {
      if (m.includes("wsl is not available") || m.includes("wsl.exe")) return "No WSL";
      if (m.includes("timed out")) return "Start timed out";
      if (m.includes("not installed")) return "Not set up yet";
      return "Failed";
    }
    if (m.includes("wsl is not available") || m.includes("wsl.exe")) return "WSL is not available";
    if (m.includes("timed out")) return "Portainer did not answer in time";
    if (m.includes("not installed")) return "Ubuntu is not set up yet";
    return "That did not work";
  }

  const heading = $derived(title ?? friendly(detail, compact));
</script>

<div class="note" class:compact role="alert">
  <span class="badge" aria-hidden="true">!</span>
  <strong class="title">{heading}</strong>
  {#if detail && detail !== heading}
    <button class="info" type="button" aria-label="More information" aria-expanded={open} onclick={() => (open = !open)}>
      i
    </button>
  {/if}
  {#if open && showDetail}
    <p class="detail">{detail}</p>
  {/if}
</div>

<style>
  .note {
    display: grid;
    grid-template-columns: auto 1fr auto;
    align-items: center;
    gap: 8px 10px;
    flex: none;
    padding: 9px 10px 9px 12px;
    border: 1px solid color-mix(in srgb, var(--failed) 40%, transparent);
    border-radius: 12px;
    background: color-mix(in srgb, var(--failed) 9%, var(--card));
    color: var(--text);
    text-align: left;
  }
  .note.compact {
    padding: 3px 4px 3px 7px;
    gap: 6px;
    border-radius: 9px;
  }
  .badge {
    display: grid;
    place-items: center;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: var(--failed);
    color: #fff;
    font: 12px/1 var(--font-display);
  }
  .compact .badge {
    width: 16px;
    height: 16px;
    font-size: 10px;
  }
  .title {
    min-width: 0;
    overflow: hidden;
    font-size: 13px;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .compact .title {
    font-size: 12px;
  }
  .info {
    display: grid;
    place-items: center;
    width: 24px;
    height: 24px;
    padding: 0;
    border: 1.5px solid color-mix(in srgb, var(--failed) 55%, transparent);
    border-radius: 50%;
    background: transparent;
    color: var(--failed);
    font: italic 700 13px/1 Georgia, serif;
    cursor: pointer;
    --wails-draggable: no-drag;
  }
  .compact .info {
    width: 18px;
    height: 18px;
    font-size: 11px;
  }
  .info:hover,
  .info[aria-expanded="true"] {
    background: var(--failed);
    color: #fff;
  }
  .detail {
    grid-column: 1 / -1;
    margin: 0;
    padding-top: 6px;
    border-top: 1px dashed color-mix(in srgb, var(--failed) 35%, transparent);
    color: var(--text-muted);
    font: 11.5px/1.4 Consolas, "Cascadia Mono", monospace;
    word-break: break-word;
    user-select: text;
  }
</style>
