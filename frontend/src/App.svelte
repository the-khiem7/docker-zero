<script lang="ts">
  import Popup from "./Popup.svelte";
  import Corner from "./Corner.svelte";

  const view = new URLSearchParams(location.search).get("view");
</script>

{#if view === "corner"}
  <Corner />
{:else}
  <Popup />
{/if}
