// Theme switch for rendered illustrations: ?theme=dark selects the dark tokens in the stylesheet.
document.documentElement.dataset.theme = new URLSearchParams(location.search).get("theme") === "dark" ? "dark" : "light";

// Line glyphs for the badges and the section icons, drawn in a 100 box and inked in the logo's dark teal.
const GLYPHS = {
  setup: '<rect x="30" y="26" width="40" height="32" rx="5"/><path d="M30,35 H70"/><path d="M41,47 L47,53 L59,42"/>',
  state: '<circle cx="50" cy="42" r="14"/><circle cx="50" cy="42" r="4.5" fill="#12343b" stroke="none"/>',
  open: '<path d="M46,30 H38 a6,6 0 0 0 -6,6 V52 a6,6 0 0 0 6,6 H56 a6,6 0 0 0 6,-6 V46"/><path d="M52,48 L68,32"/><path d="M58,30 H70 V42"/>',
  up: '<path d="M30,43 H41 L46,30 L55,56 L60,43 H70"/>',
  upkeep: '<path d="M65,33 A16,16 0 1 0 66,48"/><path d="M67,25 V34 H58"/>',
  down: '<path d="M50,25 V47"/><path d="M39,37 L50,48 L61,37"/><path d="M34,58 H66"/>',
  take: '<path d="M50,26 V47"/><path d="M40,38 L50,48 L60,38"/><path d="M32,53 V58 H68 V53"/>',
  story: '<path d="M50,32 C43,28 36,28 31,30 V56 C36,54 43,54 50,58 C57,54 64,54 69,56 V30 C64,28 57,28 50,32 Z"/><path d="M50,32 V58"/>',
  get: '<rect x="31" y="28" width="16" height="13" rx="3"/><rect x="53" y="28" width="16" height="13" rx="3"/><rect x="31" y="46" width="16" height="13" rx="3"/><rect x="53" y="46" width="16" height="13" rx="3"/>',
  start: '<path d="M42,28 L68,43 L42,58 Z"/>',
  day: '<circle cx="50" cy="43" r="8"/><path d="M50,28 V31 M50,55 V58 M35,43 H38 M62,43 H65 M39.4,32.4 L41.5,34.5 M60.5,51.5 L62.6,53.6 M60.5,34.5 L62.6,32.4 M39.4,53.6 L41.5,51.5"/>',
  safe: '<path d="M50,26 L67,32 V43 C67,52 60,57 50,61 C40,57 33,52 33,43 V32 Z"/><path d="M43,43 L48,48 L58,37"/>',
  shield: '<path d="M50,26 L67,32 V43 C67,52 60,57 50,61 C40,57 33,52 33,43 V32 Z"/>',
  faq: '<path d="M39,35 C39,24 61,24 61,35 C61,43 50,43 50,50"/><circle cx="50" cy="58" r="2.8" fill="#12343b" stroke="none"/>',
  build: '<path d="M42,30 L30,43 L42,56"/><path d="M58,30 L70,43 L58,56"/>',
  gear: '<circle cx="50" cy="43" r="11"/><path d="M50,24 V29 M50,57 V62 M31,43 H36 M64,43 H69 M36.6,29.6 L40.1,33.1 M59.9,52.9 L63.4,56.4 M63.4,29.6 L59.9,33.1 M40.1,52.9 L36.6,56.4" stroke-width="5.5"/>',
  folder: '<path d="M30,32 H43 L48,38 H70 V58 H30 Z"/>',
  lock: '<rect x="35" y="39" width="30" height="21" rx="4"/><path d="M41,39 V34 a9,9 0 0 1 18,0 V39"/>',
  term: '<path d="M33,32 L45,43 L33,54"/><path d="M52,56 H67"/>',
  inf: '<path d="M50,43 C44,32 31,33 31,43 C31,53 44,54 50,43 C56,32 69,33 69,43 C69,53 56,54 50,43 Z"/>',
  moon: '<path d="M58,28 A16,16 0 1 0 69,52 A13,13 0 0 1 58,28 Z"/>',
};
const glyphs = Object.entries(GLYPHS).map(([id, d]) =>
  `<symbol id="g-${id}" viewBox="0 0 100 100"><g fill="none" stroke="#12343b" stroke-width="4.5" stroke-linecap="round" stroke-linejoin="round">${d}</g></symbol>`).join("");

// The sea along the foot of a section banner: the same two teals as the logo.
const wave = (y) => `M0,${y} Q60,${y - 8} 120,${y}` + Array.from({ length: 9 }, (_, i) => ` T${240 + i * 120},${y}`).join("") + " L1200,26 L0,26 Z";

// The logo as one sprite: the ring takes currentColor, so a figure can recolour it per state.
// The badge is the logo's ring and sea without the orca: the frame of every icon.
document.body.insertAdjacentHTML("afterbegin", `<svg width="0" height="0" style="position:absolute" aria-hidden="true"><defs>
<clipPath id="mark-sea"><circle cx="50" cy="50" r="24.5"/></clipPath>
<clipPath id="badge-sea"><circle cx="50" cy="50" r="34"/></clipPath>
<symbol id="mark" viewBox="0 0 100 100">
<g clip-path="url(#mark-sea)"><rect width="100" height="100" fill="#fff4e6"/>
<path d="M0,60 Q10,56 20,60 T40,60 T60,60 L100,60 L100,100 L0,100 Z" fill="#7fd1c7"/>
<path d="M0,68 Q10,64 20,68 T40,68 T60,68 L100,68 L100,100 L0,100 Z" fill="#2a9d8f"/></g>
<circle cx="50" cy="50" r="30" fill="none" stroke="currentColor" stroke-width="11"/>
<g transform="translate(50,50) scale(0.8) translate(-54,-52)">
<path d="M10,47 C11,33 26,22 44,21 C49,21 53,21 56,22 C60,13 66,7 75,4 C71,12 69,18 69,25 C80,32 86,44 86,57 C86,64 84,70 80,75 C86,77 92,80 97,85 C90,88 83,88 77,85 C77,90 75,95 71,99 C68,92 68,85 71,78 C72,66 68,55 60,48 C50,41 36,41 26,46 C20,49 13,50 10,47 Z" fill="#111111"/>
<path d="M44,43 C42,54 37,62 30,68 C40,66 49,58 54,45 Z" fill="#111111"/>
<g fill="#ffffff"><path d="M11,47 C20,44 30,38 44,37 C58,36 70,44 76,56 C79,63 79,70 76,76 L71,78 C72,66 68,55 60,48 C50,41 36,41 26,46 C20,49 13,50 11,47 Z"/><ellipse cx="30" cy="32" rx="6" ry="2.6" transform="rotate(-22 30 32)"/></g>
<path d="M44,43 C42,54 37,62 30,68 C40,66 49,58 54,45 Z" fill="#111111"/></g>
<path d="M78.98,57.76 A30,30 0 0 1 55.21,79.54" fill="none" stroke="currentColor" stroke-width="11" stroke-linecap="butt"/>
</symbol>
<symbol id="badge" viewBox="0 0 100 100">
<g clip-path="url(#badge-sea)"><rect width="100" height="100" fill="#fff4e6"/>
<path d="M0,66 Q10,62 20,66 T40,66 T60,66 T80,66 T100,66 L100,100 L0,100 Z" fill="#7fd1c7"/>
<path d="M0,74 Q10,70 20,74 T40,74 T60,74 T80,74 T100,74 L100,100 L0,100 Z" fill="#2a9d8f"/></g>
<circle cx="50" cy="50" r="38" fill="none" stroke="currentColor" stroke-width="9"/>
</symbol>
<symbol id="edge" viewBox="0 0 1200 26" preserveAspectRatio="none">
<path d="${wave(10)}" fill="#7fd1c7"/><path d="${wave(17)}" fill="#2a9d8f"/>
</symbol>
${glyphs}</defs></svg>`);
