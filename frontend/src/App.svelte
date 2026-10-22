<script lang="ts">
  import Popup from "./Popup.svelte";
  import Corner from "./Corner.svelte";
  import Settings from "./Settings.svelte";
  import Setup from "./Setup.svelte";

  const view = new URLSearchParams(location.search).get("view");
</script>

{#if view === "corner"}
  <Corner />
{:else if view === "settings"}
  <Settings />
{:else if view === "setup"}
  <Setup />
{:else}
  <Popup />
{/if}
