<script lang="ts">
  import Popup from "./Popup.svelte";
  import Corner from "./Corner.svelte";
  import Settings from "./Settings.svelte";
  import Setup from "./Setup.svelte";

  const q = new URLSearchParams(location.search);
  const view = q.get("view");
  const f = q.get("frame");
  const frame = f === null ? undefined : Number(f);
</script>

{#if view === "corner"}
  <Corner />
{:else if view === "settings"}
  <Settings />
{:else if view === "setup"}
  <Setup />
{:else}
  <Popup {frame} />
{/if}
