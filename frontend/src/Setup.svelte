<script lang="ts">
  import { Events } from "@wailsio/runtime";
  import * as S from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/setupservice";
  import * as A from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/adoptservice";
  import type { Report } from "../bindings/github.com/the-khiem7/docker-zero/internal/adopt/models";
  import type { SetupView } from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/models";
  import { run, service, live } from "./lib/status.svelte";

  const minLength = 12;

  let v = $state<SetupView | null>(null);
  let username = $state("admin");
  let secret = $state("");
  let again = $state("");
  let error = $state("");
  let logBox: HTMLElement | undefined = $state();

  let found = $state<Report | null>(null);
  let chosen = $state<Record<string, boolean>>({});
  let adopting = $state(false);
  let adoptError = $state("");
  let adopted = $state<string[]>([]);

  async function detect() {
    try {
      const r = await A.Detect();
      found = r;
      for (const f of r.findings ?? []) chosen[f.item] ??= true;
    } catch (err) {
      adoptError = err instanceof Error ? err.message : String(err);
    }
  }

  async function takeOver() {
    adopting = true;
    adoptError = "";
    try {
      adopted = (await A.Adopt(Object.keys(chosen).filter((k) => chosen[k]))) ?? [];
    } catch (err) {
      adoptError = err instanceof Error ? err.message : String(err);
    } finally {
      adopting = false;
      await detect();
    }
  }

  const hasExisting = $derived(
    !!found &&
      found.distro !== "not-installed" &&
      ((found.findings?.length ?? 0) > 0 || (found.container !== "" && found.container !== "missing")),
  );

  const busy = $derived(v?.running ?? false);
  const done = $derived(v?.done ?? false);
  const editable = $derived(!busy && !done);

  const icons: Record<string, string> = {
    running: "◐",
    done: "✓",
    skipped: "↷",
    failed: "✕",
  };

  Events.On("setup", (e) => {
    v = e.data;
  });
  S.State().then((s) => (v = s));
  detect();
  window.addEventListener("focus", detect);

  $effect(() => {
    // Keep the newest log line in view.
    void v?.log?.length;
    if (logBox) logBox.scrollTop = logBox.scrollHeight;
  });

  async function start() {
    error = "";
    if (!username.trim()) {
      error = "Enter an admin username.";
    } else if (secret.length < minLength) {
      error = `The password needs at least ${minLength} characters.`;
    } else if (secret !== again) {
      error = "The two passwords differ.";
    }
    if (error) return;
    try {
      await S.Start(username.trim(), secret);
      secret = "";
      again = "";
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  async function open() {
    await run(() => service.OpenPortainer());
  }
</script>

<main>
  <h1>Set up Docker Zero</h1>
  <p class="muted">
    Installs Ubuntu-24.04 in WSL, Docker Engine and Portainer for your account. It needs no admin rights and takes
    several minutes, mostly for downloads.
  </p>

  {#if hasExisting && found}
    <section class="card existing">
      <h2>Existing setup found</h2>
      <p class="muted small">
        Ubuntu-24.04{found.container && found.container !== "missing" ? `, Portainer (${found.container})` : ""} already
        exists. Taking over keeps the distro, containers and Portainer data as they are and replaces the launchers below
        with Docker Zero's own. You do not need to run setup.
      </p>
      {#each found.findings ?? [] as f (f.item)}
        <label class="row">
          <input type="checkbox" bind:checked={chosen[f.item]} disabled={adopting} />
          <span>
            {f.title}
            <span class="muted small block">{f.detail}</span>
          </span>
        </label>
      {/each}
      {#if found.findings?.length}
        <button class="primary" disabled={adopting || !Object.values(chosen).some(Boolean)} onclick={takeOver}>
          {adopting ? "Working..." : "Take over"}
        </button>
      {:else}
        <p class="small">Nothing left to replace.</p>
      {/if}
      {#if adopted.length}
        <p class="ok small">
          Replaced {adopted.length} item{adopted.length === 1 ? "" : "s"}. Open a new terminal to use the docker command.
        </p>
      {/if}
      {#if adoptError}<p class="error">{adoptError}</p>{/if}
    </section>
  {/if}

  {#if editable}
    <section class="card">
      <h2>Portainer admin account</h2>
      <label class="field">
        Username
        <input type="text" bind:value={username} autocomplete="off" spellcheck="false" />
      </label>
      <label class="field">
        Password
        <input type="password" bind:value={secret} autocomplete="new-password" />
      </label>
      <label class="field">
        Repeat password
        <input type="password" bind:value={again} autocomplete="new-password" />
      </label>
      <p class="muted small">At least {minLength} characters. Docker Zero never stores it; Portainer keeps its own copy.</p>
      <button class="primary" onclick={start}>{v?.error ? "Try again" : "Start setup"}</button>
    </section>
  {/if}

  {#if v}
    <section class="card">
      <h2>Steps</h2>
      <ul class="steps">
        {#each v.steps as step (step.id)}
          <li class={step.status || "pending"}>
            <span class="icon" aria-hidden="true">{icons[step.status] ?? "○"}</span>
            <span class="title">{step.title}</span>
            {#if step.status === "skipped"}<span class="note">already done</span>{/if}
          </li>
        {/each}
      </ul>
    </section>

    {#if v.log?.length}
      <pre class="log" bind:this={logBox}>{v.log.join("\n")}</pre>
    {/if}
  {/if}

  {#if error || v?.error}
    <p class="error">{error || v?.error}</p>
  {/if}
  {#if live.error}
    <p class="error">{live.error}</p>
  {/if}

  {#if done}
    <section class="card finish">
      <h2>All set</h2>
      <p class="muted">Portainer is ready. Sign in with the account you just chose. Open a new terminal to use the docker command.</p>
      <button class="primary" onclick={open}>Open Portainer</button>
    </section>
  {/if}
</main>

<style>
  main {
    height: 100%;
    overflow-y: auto;
    padding: 16px 20px 24px;
    user-select: text;
  }
  h1 {
    margin: 0 0 6px;
    font-size: 20px;
    font-weight: 600;
  }
  h2 {
    margin: 0 0 8px;
    font-size: 14px;
    font-weight: 600;
  }
  .card {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 12px 14px;
    margin: 12px 0;
  }
  .muted {
    color: var(--text-muted);
  }
  .small {
    font-size: 12px;
  }
  .block {
    display: block;
  }
  .row {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    padding: 3px 0;
  }
  .ok {
    color: var(--ok);
  }
  .existing {
    border-color: var(--accent);
  }
  .field {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 8px;
  }
  .field input {
    padding: 6px 8px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--surface);
    color: inherit;
  }
  button.primary {
    padding: 8px 14px;
    border: 1px solid var(--accent);
    border-radius: var(--radius);
    background: var(--accent);
    color: var(--accent-text);
    cursor: pointer;
  }
  .steps {
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .steps li {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 0;
    color: var(--text-muted);
  }
  .steps li.running,
  .steps li.done {
    color: inherit;
  }
  .steps li.done .icon {
    color: var(--ok);
  }
  .steps li.failed .icon {
    color: var(--bad);
  }
  .steps li.running .icon {
    animation: pulse 1s ease-in-out infinite;
  }
  .icon {
    width: 16px;
    text-align: center;
  }
  .note {
    margin-left: auto;
    font-size: 12px;
  }
  .log {
    max-height: 180px;
    overflow: auto;
    margin: 12px 0;
    padding: 8px 10px;
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    font-size: 11px;
    white-space: pre-wrap;
    word-break: break-word;
  }
  .error {
    color: var(--bad);
    font-size: 13px;
    word-break: break-word;
  }
  @keyframes pulse {
    50% {
      opacity: 0.3;
    }
  }
</style>
