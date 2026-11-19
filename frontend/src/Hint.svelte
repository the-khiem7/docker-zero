<script lang="ts">
  /** The i button: the long explanation stays behind it until asked for. */
  let { label, text, code = "" }: { label: string; text: string; code?: string } = $props();
  let open = $state(false);
  let btn: HTMLButtonElement | undefined = $state();
  let pos = $state("");

  const WIDTH = 260;
  const EDGE = 10;

  // The bubble is fixed to the window and kept inside it: a card on the left edge would
  // otherwise push a right-aligned bubble past the window border, where it is cut off.
  function place() {
    if (!btn) return;
    const b = btn.getBoundingClientRect();
    const w = Math.min(WIDTH, window.innerWidth - 2 * EDGE);
    const left = Math.max(EDGE, Math.min(b.right - w, window.innerWidth - EDGE - w));
    pos = `top:${b.bottom + 6}px;left:${left}px;width:${w}px`;
  }

  function toggle() {
    open = !open;
    if (open) place();
  }

  function outside(e: PointerEvent) {
    if (open && btn && !btn.parentElement?.contains(e.target as Node)) open = false;
  }
</script>

<svelte:window
  onkeydown={(e) => e.key === "Escape" && (open = false)}
  onpointerdown={outside}
  onresize={() => open && place()}
  onscrollcapture={() => (open = false)}
/>

<span class="hint">
  <button bind:this={btn} class="i" type="button" aria-label={"About " + label} aria-expanded={open} onclick={toggle}>i</button>
  {#if open}<span class="bubble" role="note" style={pos}>{text}{#if code}<code>{code}</code>{/if}</span>{/if}
</span>

<style>
  .hint {
    position: relative;
    flex: none;
  }
  .i {
    display: grid;
    place-items: center;
    width: 24px;
    height: 24px;
    padding: 0;
    border: 1.5px solid color-mix(in srgb, var(--teal) 55%, transparent);
    border-radius: 50%;
    background: transparent;
    color: var(--teal);
    font: italic 700 13px/1 Georgia, serif;
    cursor: pointer;
  }
  .i:hover,
  .i[aria-expanded="true"] {
    background: var(--teal);
    color: var(--on-teal);
  }
  .i:focus-visible {
    outline: 2px solid var(--teal);
    outline-offset: 2px;
  }
  .bubble {
    position: fixed;
    z-index: 20;
    display: block;
    padding: 10px 12px;
    border: 1px solid var(--line);
    border-radius: 16px;
    background: var(--card);
    box-shadow: 0 10px 24px rgb(0 0 0 / 0.16);
    color: var(--text);
    font-size: 12px;
    font-style: normal;
    font-weight: 400;
    letter-spacing: normal;
    line-height: 1.45;
    text-align: left;
    text-transform: none;
    white-space: normal;
  }
  code {
    display: block;
    margin-top: 8px;
    padding: 5px 9px;
    border-radius: 10px;
    background: var(--sea-1);
    font: 11px/1.4 Consolas, "Cascadia Mono", monospace;
    overflow-wrap: anywhere;
  }
</style>
