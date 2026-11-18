<script lang="ts">
  /* Docker Zero glyphs: round 1.8 strokes on a 24 grid, one accent shape filled
     with the logo's sea. Every glyph is drawn for this app, none comes from a set. */
  const c = (x: number, y: number, r: number) => `M${x - r} ${y}a${r} ${r} 0 1 0 ${2 * r} 0a${r} ${r} 0 1 0 ${-2 * r} 0Z`;
  const r = (x: number, y: number, w: number, h: number, k: number) =>
    `M${x + k} ${y}h${w - 2 * k}a${k} ${k} 0 0 1 ${k} ${k}v${h - 2 * k}a${k} ${k} 0 0 1 ${-k} ${k}h${-(w - 2 * k)}a${k} ${k} 0 0 1 ${-k} ${-k}v${-(h - 2 * k)}a${k} ${k} 0 0 1 ${k} ${-k}Z`;

  type G = { s: string[]; w?: string[] };
  const shield = "M12 3 4.5 6v5.5c0 4.5 3 7.5 7.5 9.5 4.5-2 7.5-5 7.5-9.5V6Z";
  const set: Record<string, G> = {
    app: { s: [c(12, 12, 9)], w: ["M3.2 13.5Q7.6 10.5 12 13.5T20.8 13.5A9 9 0 0 1 3.2 13.5Z"] },
    power: { s: ["M12 3v8", "M6.6 6.6a8 8 0 1 0 10.8 0"], w: ["M6.6 17.4a8 8 0 0 0 10.8 0Z"] },
    distro: { s: ["M12 3 3.5 7.5v9L12 21l8.5-4.5v-9Z", "M3.5 7.5 12 12l8.5-4.5", "M12 12v9"], w: ["M12 12l8.5-4.5v9L12 21Z"] },
    monitor: { s: [r(3, 4, 18, 12, 2.5), "M9 20h6", "M12 16v4"], w: ["M3 12h18v1.5a2.5 2.5 0 0 1-2.5 2.5h-13A2.5 2.5 0 0 1 3 13.5Z"] },
    menu: { s: [r(4, 4, 7, 7, 2), r(13, 4, 7, 7, 2), r(4, 13, 7, 7, 2), r(13, 13, 7, 7, 2)], w: [r(13, 13, 7, 7, 2)] },
    anchor: { s: [c(12, 5, 2), "M12 7v13", "M7.5 11.5h9", "M4.5 14a7.5 7.5 0 0 0 15 0"], w: ["M4.5 14a7.5 7.5 0 0 0 15 0Z"] },
    infinity: {
      s: ["M12 12c-2-2.6-3.4-4-5.4-4a4 4 0 0 0 0 8c2 0 3.4-1.4 5.4-4s3.4-4 5.4-4a4 4 0 0 1 0 8c-2 0-3.4-1.4-5.4-4Z"],
      w: ["M6.6 8a4 4 0 0 0 0 8c2 0 3.4-1.4 5.4-4-2-2.6-3.4-4-5.4-4Z"],
    },
    moon: { s: ["M19.5 14.5A8 8 0 0 1 9.5 4.5a8 8 0 1 0 10 10Z"], w: ["M19.5 14.5A8 8 0 0 1 9.5 4.5 8 8 0 0 0 19.5 14.5Z"] },
    globe: { s: [c(12, 12, 9), "M3 12h18", "M12 3c3.2 3 3.2 15 0 18", "M12 3c-3.2 3-3.2 15 0 18"], w: ["M3.4 14.5h17.2A9 9 0 0 1 3.4 14.5Z"] },
    compass: { s: [c(12, 12, 9), "M15.8 8.2l-2 5.6-5.6 2 2-5.6Z"], w: ["M15.8 8.2l-2 5.6-5.6 2Z"] },
    window: { s: [r(3, 5, 18, 14, 3), "M3 9.5h18", "M6.5 7.3v.1", "M9.3 7.3v.1"], w: ["M3 9.5h18V16a3 3 0 0 1-3 3H6a3 3 0 0 1-3-3Z"] },
    shield: { s: [shield, "M8.8 12.2l2.4 2.4 4-4.6"], w: ["M4.6 12h14.8c-.4 4-3.2 6.5-7.4 8.4C7.8 18.5 5 16 4.6 12Z"] },
    shieldAlert: { s: [shield, "M12 8.5v4", "M12 15.6v.1"], w: ["M4.6 12h14.8c-.4 4-3.2 6.5-7.4 8.4C7.8 18.5 5 16 4.6 12Z"] },
    shieldOff: { s: ["M12 3 4.5 6v5.5c0 4.5 3 7.5 7.5 9.5 4.5-2 7.5-5 7.5-9.5V6Z", "M9 9l6 6", "M15 9l-6 6"] },
    calendar: { s: [r(3.5, 5, 17, 15, 3), "M3.5 10h17", "M8 3v4", "M16 3v4"], w: ["M3.5 10h17v7a3 3 0 0 1-3 3h-11a3 3 0 0 1-3-3Z"] },
    clock: { s: [c(12, 12, 9), "M12 7v5l3 2"], w: ["M3.4 14.5h17.2A9 9 0 0 1 3.4 14.5Z"] },
    hash: { s: ["M9.5 4 7.5 20", "M16.5 4l-2 16", "M4.5 9h16", "M3.5 15h16"] },
    update: { s: ["M19.5 12a7.5 7.5 0 1 1-2.4-5.5", "M19.5 4.5v4h-4", "M12 9v6", "M9.5 11.5 12 9l2.5 2.5"], w: [c(12, 12, 4)] },
    refresh: { s: ["M19.5 12a7.5 7.5 0 1 1-2.4-5.5", "M19.5 4.5v4h-4"] },
    key: { s: [c(8, 15, 4), "M11 12l9-9", "M16 7l3 3", "M14 9l2 2"], w: [c(8, 15, 4)] },
    user: { s: [c(12, 8, 4), "M4.5 20c.8-3.8 3.7-6 7.5-6s6.700 2.200 7.500 6"], w: ["M4.5 20c.8-3.8 3.7-6 7.5-6s6.7 2.2 7.5 6Z"] },
    lock: { s: [r(5, 11, 14, 9, 2.5), "M8 11V8a4 4 0 0 1 8 0v3"], w: [r(5, 14.5, 14, 5.5, 2)] },
    trash: { s: ["M4 7h16", "M9 7V4.5h6V7", "M6 7l1 13h10l1-13", "M10 11v5", "M14 11v5"], w: ["M6.4 12h11.2l-.6 8H7Z"] },
    backup: { s: ["M4 15v4h16v-4", "M12 4v11", "M7.500 10.500 12 15l4.500-4.500"], w: ["M4 15v4h16v-4Z"] },
    terminal: { s: [r(3, 5, 18, 14, 3), "M7 10l3 2.500L7 15", "M13 15h4"], w: [r(3, 14, 18, 5, 3)] },
    sliders: { s: ["M4 8h9", "M17 8h3", "M4 16h3", "M11 16h9", c(15, 8, 2), c(9, 16, 2)], w: [c(15, 8, 2), c(9, 16, 2)] },
    file: { s: ["M7 3h7l5 5v13H7Z", "M14 3v5h5"], w: ["M14 3v5h5Z"] },
    folder: { s: ["M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2Z"], w: ["M3 12h18v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2Z"] },
    check: { s: ["M5 12.500 9.500 17 19 7.500"] },
    cross: { s: ["M6 6l12 12", "M18 6 6 18"] },
    skip: { s: ["M5 12h9", "M12 8l4 4-4 4", "M19 7v10"] },
    wave: { s: ["M3 9.500q3-3 6 0t6 0 6 0", "M3 15.500q3-3 6 0t6 0 6 0"], w: ["M3 15.500q3-3 6 0t6 0 6 0V19H3Z"] },
    alert: { s: ["M12 4 2.800 19.500h18.400Z", "M12 10v4", "M12 16.800v.1"], w: ["M12 4l9.200 15.500H2.800Z"] },
    spark: { s: ["M12 3v4", "M12 17v4", "M3 12h4", "M17 12h4", "M6 6l2.500 2.500", "M15.500 15.500 18 18", "M18 6l-2.500 2.500", "M8.500 15.500 6 18"] },
  };

  let { name, size = 24 }: { name: string; size?: number } = $props();
  const g = $derived(set[name] ?? set.app);
</script>

<svg
  class="ic"
  viewBox="0 0 24 24"
  width={size}
  height={size}
  fill="none"
  stroke="currentColor"
  stroke-width="1.8"
  stroke-linecap="round"
  stroke-linejoin="round"
  aria-hidden="true"
>
  {#each g.w ?? [] as d, i (i)}<path {d} class="w" />{/each}
  {#each g.s as d, i (i)}<path {d} />{/each}
</svg>

<style>
  .ic {
    flex: none;
  }
  .w {
    fill: var(--water-a);
    fill-opacity: 0.4;
    stroke: none;
  }
</style>
