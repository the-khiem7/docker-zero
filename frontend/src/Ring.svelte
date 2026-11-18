<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    parts: number;
    max?: number;
    level: string;
    size?: number;
    orca?: boolean;
    /** Board only: freeze the leap at this point, 0 to 1. */
    frame?: number;
    id: string;
    children?: Snippet;
  }

  let { parts, max = 4, level, size = 120, id, orca = false, frame, children }: Props = $props();

  // The disc spans y 13..87 in the 100 unit box; the water top moves from the bottom to the top edge.
  const top = $derived(87 - 74 * Math.min(1, Math.max(0, parts / max)));
  // Logo geometry: sea inside r 24.5, ring r 30. The waterline runs from the disc bottom up to the orca's belly.
  const sea = $derived(80 - 26 * Math.min(1, Math.max(0, parts / max)));
  const LEAP_MS = 1400;
  let leaping = $state(false);
  let prev = "";
  // One leap when the state reaches running; never looped.
  $effect(() => {
    const now = level;
    if (now === "running" && prev && prev !== "running") {
      leaping = true;
      const t = setTimeout(() => (leaping = false), LEAP_MS);
      prev = now;
      return () => clearTimeout(t);
    }
    prev = now;
  });
  const orcaStyle = $derived(
    frame === undefined ? "" : `animation: leap ${LEAP_MS}ms var(--ease) both paused; animation-delay: -${Math.round(frame * LEAP_MS)}ms`,
  );
  const BODY =
    "M10,47 C11,33 26,22 44,21 C49,21 53,21 56,22 C60,13 66,7 75,4 C71,12 69,18 69,25 C80,32 86,44 86,57 C86,64 84,70 80,75 C86,77 92,80 97,85 C90,88 83,88 77,85 C77,90 75,95 71,99 C68,92 68,85 71,78 C72,66 68,55 60,48 C50,41 36,41 26,46 C20,49 13,50 10,47 Z";
</script>

<div class="ring" style="width:{size}px;height:{size}px;--ring:var(--{level})">
  {#if orca}
    <svg viewBox="0 0 100 100" aria-hidden="true">
      <defs>
        <clipPath id="sea-{id}"><circle cx="50" cy="50" r="24.5" /></clipPath>
      </defs>
      <g clip-path="url(#sea-{id})">
        <rect width="100" height="100" fill="#fff4e6" />
        <g class="water" style="transform: translateY({sea - 60}px)">
          <path d="M-10,60 Q0,56 10,60 T30,60 T50,60 T70,60 T90,60 T110,60 L110,120 L-10,120 Z" fill="#7fd1c7" />
          <path d="M-10,68 Q0,64 10,68 T30,68 T50,68 T70,68 T90,68 T110,68 L110,120 L-10,120 Z" fill="#2a9d8f" />
        </g>
      </g>
      <circle cx="50" cy="50" r="30" class="o-band" />
      <g class="orca" class:leaping style={orcaStyle}>
      <g transform="translate(50,50) scale(0.8) translate(-54,-52)">
        <path d={BODY} fill="#111111" />
        <path d="M44,43 C42,54 37,62 30,68 C40,66 49,58 54,45 Z" fill="#111111" />
        <g fill="#ffffff">
          <path d="M11,47 C20,44 30,38 44,37 C58,36 70,44 76,56 C79,63 79,70 76,76 L71,78 C72,66 68,55 60,48 C50,41 36,41 26,46 C20,49 13,50 11,47 Z" />
          <ellipse cx="30" cy="32" rx="6" ry="2.6" transform="rotate(-22 30 32)" />
        </g>
        <path d="M44,43 C42,54 37,62 30,68 C40,66 49,58 54,45 Z" fill="#111111" />
      </g>
      </g>
      <path d="M78.98,57.76 A30,30 0 0 1 55.21,79.54" class="o-band" stroke-linecap="butt" />
    </svg>
  {:else}
    <svg viewBox="0 0 100 100" aria-hidden="true">
      <defs>
        <clipPath id="disc-{id}"><circle cx="50" cy="50" r="37" /></clipPath>
      </defs>
      <circle cx="50" cy="50" r="37" class="disc" />
      <g clip-path="url(#disc-{id})">
        <g class="water" style="transform: translateY({top}px)">
          <path class="back" d="M-25 6 Q-12.5 -1 0 6 T25 6 T50 6 T75 6 T100 6 T125 6 V110 H-25 Z" />
          <path class="front" d="M-25 9 Q-12.5 2 0 9 T25 9 T50 9 T75 9 T100 9 T125 9 V110 H-25 Z" />
        </g>
      </g>
      <circle cx="50" cy="50" r="44" class="band" />
    </svg>
  {/if}
  {#if children}<div class="center">{@render children()}</div>{/if}
</div>

<style>
  .ring {
    position: relative;
    flex: none;
  }
  svg {
    width: 100%;
    height: 100%;
    display: block;
  }
  .disc {
    fill: var(--sea-1);
  }
  .water {
    transition: transform 1100ms var(--ease);
  }
  .back {
    fill: var(--water-a);
    opacity: 0.55;
  }
  .front {
    fill: var(--water-b);
  }
  .band {
    fill: none;
    stroke: var(--ring);
    stroke-width: 6.5;
    transition: stroke 400ms;
  }
  .orca {
    transform-origin: 50px 62px;
  }
  .orca.leaping {
    animation: leap 1400ms var(--ease) both;
  }
  @keyframes -global-leap {
    0% {
      transform: translateY(9px) rotate(-9deg);
    }
    45% {
      transform: translateY(-7px) rotate(5deg);
    }
    100% {
      transform: translateY(0) rotate(0);
    }
  }
  .o-band {
    fill: none;
    stroke: var(--ring);
    stroke-width: 11;
    transition: stroke 400ms;
  }
  .center {
    position: absolute;
    inset: 0;
    display: grid;
    place-items: center;
    text-align: center;
  }
</style>
