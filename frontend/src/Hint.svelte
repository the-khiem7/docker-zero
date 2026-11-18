<script lang="ts">
  /** The i button: the long explanation stays behind it until asked for. */
  let { label, text }: { label: string; text: string } = $props();
  let open = $state(false);
</script>

<span class="hint">
  <button class="i" type="button" aria-label={"About " + label} aria-expanded={open} onclick={() => (open = !open)}>i</button>
  {#if open}<span class="bubble" role="note">{text}</span>{/if}
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
  .bubble {
    position: absolute;
    z-index: 5;
    top: 30px;
    right: 0;
    width: 250px;
    padding: 10px 12px;
    border: 1px solid var(--line);
    border-radius: 14px 4px 14px 14px;
    background: var(--card);
    box-shadow: 0 10px 24px rgb(0 0 0 / 0.16);
    color: var(--text);
    font-size: 12px;
    line-height: 1.45;
    text-align: left;
  }
</style>
