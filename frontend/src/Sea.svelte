<script lang="ts">
  /* The logo's two-tone sea inside a round control: a back swell, a front swell with a foam crest,
     and two bubbles. The water rises and the swells shift once when the level changes; nothing loops,
     so a resident tray app stays idle on CPU. */
  // simple: one pale swell for small controls, where layers, foam and bubbles would bury the icon.
  let { level = 0.5, deep = false, simple = false }: { level?: number; deep?: boolean; simple?: boolean } = $props();
  const y = $derived(100 - 100 * Math.min(1, Math.max(0, level)));
</script>

<svg class="sea" class:deep viewBox="0 0 100 100" preserveAspectRatio="none" aria-hidden="true">
  <g class="tide" style="transform: translateY({y}px)">
    {#if simple}
      <path class="calm" d="M0 6Q25 -4 50 6T100 6V120H0Z" />
    {:else}
      <path class="back" d="M-60 4Q-47.5 -6 -35 4T-10 4T15 4T40 4T65 4T90 4T115 4T140 4T165 4V120H-60Z" />
      <path class="front" d="M-60 12Q-45 3 -30 12T0 12T30 12T60 12T90 12T120 12T150 12V120H-60Z" />
      <path class="foam" d="M-60 12Q-45 3 -30 12T0 12T30 12T60 12T90 12T120 12T150 12" />
      <circle class="bub" cx="30" cy="34" r="3.2" />
      <circle class="bub" cx="70" cy="46" r="2.2" />
    {/if}
  </g>
</svg>

<style>
  .sea {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
    border-radius: 50%;
    pointer-events: none;
  }
  .tide {
    transition: transform 900ms var(--ease);
  }
  .back {
    fill: var(--water-a);
    opacity: 0.45;
    transition: transform 1200ms var(--ease);
  }
  .front {
    fill: var(--water-b);
    opacity: 0.85;
    transition: transform 1200ms var(--ease);
  }
  .foam {
    fill: none;
    stroke: #fff;
    stroke-width: 2.4;
    stroke-linecap: round;
    stroke-dasharray: 14 10 4 10;
    opacity: 0.55;
    transition: transform 1200ms var(--ease);
  }
  .calm {
    fill: var(--water-a);
    opacity: 0.22;
  }
  .bub {
    fill: none;
    stroke: #fff;
    stroke-width: 1.4;
    opacity: 0.5;
  }
  /* On the filled accent disc the sea is drawn in lighter and darker tones of the accent. */
  .deep .back {
    fill: color-mix(in srgb, var(--water-a) 55%, #fff);
    opacity: 0.35;
  }
  .deep .front {
    fill: color-mix(in srgb, var(--water-b) 70%, #000);
    opacity: 0.75;
  }
  /* One gentle shift of the swells on each level change, so the water looks poured, not painted. */
  :global(.on) > .sea .back {
    transform: translateX(-25px);
  }
  :global(.on) > .sea .front,
  :global(.on) > .sea .foam {
    transform: translateX(30px);
  }
  @media (prefers-reduced-motion: reduce) {
    .tide,
    .back,
    .front,
    .foam {
      transition: none;
    }
  }
</style>
