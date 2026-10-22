<script lang="ts">
  import Popup from "./Popup.svelte";
  import Corner from "./Corner.svelte";
  import Settings from "./Settings.svelte";

  const view = new URLSearchParams(location.search).get("view");
</script>

{#if view === "corner"}
  <Corner />
{:else if view === "settings"}
  <Settings />
{:else}
  <Popup />
{/if}
