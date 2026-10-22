<script lang="ts">
  import Maintenance from "./Maintenance.svelte";
  import * as S from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/settingsservice";
  import type { SettingsView } from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/models";
  import { BrowserKind, IdleMode } from "../bindings/github.com/the-khiem7/docker-zero/internal/settings/models";

  let v = $state<SettingsView | null>(null);
  let busy = $state("");
  let error = $state("");
  let confirmApply = $state(false);

  async function load() {
    v = await S.Get();
  }

  async function act(name: string, fn: () => Promise<unknown>) {
    busy = name;
    error = "";
    try {
      await fn();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      busy = "";
      await load();
    }
  }

  async function pickBrowser() {
    const path = await S.PickBrowser();
    if (path) await act("browser", () => S.SetBrowser(BrowserKind.BrowserCustom, path));
  }

  const idleModes = [
    {
      mode: IdleMode.IdleKeepalive,
      title: "Keep running while Docker Zero runs",
      note: "A hidden session holds Ubuntu-24.04 up. Changes nothing outside this app.",
    },
    {
      mode: IdleMode.IdleWSLConfig,
      title: "Never stop idle distros (.wslconfig)",
      note: "Sets instanceIdleTimeout=-1. Applies to every WSL2 distro of your account after WSL restarts.",
    },
    {
      mode: IdleMode.IdleDefault,
      title: "WSL default",
      note: "WSL stops the distro, and Portainer with it, when nothing uses it.",
    },
  ];

  const expires = $derived(v?.cert.notAfter ? new Date(v.cert.notAfter).toLocaleDateString() : "");

  load();
  window.addEventListener("focus", load);
</script>

<main>
  <h1>Settings</h1>

  {#if !v}
    <p class="muted">Loading...</p>
  {:else}
    <section class="card">
      <h2>Start with Windows</h2>
      <label class="row">
        <input
          type="checkbox"
          checked={v.autoStart.enabled}
          disabled={!!busy}
          onchange={(e) => act("autostart", () => S.SetAutoStart(e.currentTarget.checked, v!.autoStart.startDistro))}
        />
        Start Docker Zero when I sign in
      </label>
      <label class="row indent">
        <input
          type="checkbox"
          checked={v.autoStart.startDistro}
          disabled={!!busy || !v.autoStart.enabled}
          onchange={(e) => act("autostart", () => S.SetAutoStart(true, e.currentTarget.checked))}
        />
        Also start Ubuntu-24.04 and Portainer
      </label>
    </section>

    <section class="card">
      <h2>Shortcuts</h2>
      <label class="row">
        <input
          type="checkbox"
          checked={v.shortcuts.desktop}
          disabled={!!busy}
          onchange={(e) => act("shortcuts", () => S.SetShortcuts(e.currentTarget.checked, v!.shortcuts.startMenu))}
        />
        Desktop
      </label>
      <label class="row">
        <input
          type="checkbox"
          checked={v.shortcuts.startMenu}
          disabled={!!busy}
          onchange={(e) => act("shortcuts", () => S.SetShortcuts(v!.shortcuts.desktop, e.currentTarget.checked))}
        />
        Start Menu
      </label>
      <p class="muted small">Shortcuts and auto start point at {v.installPath}; the app copies itself there first.</p>
    </section>

    <section class="card">
      <h2>Keep Ubuntu-24.04 running</h2>
      {#each idleModes as m (m.mode)}
        <label class="row option">
          <input
            type="radio"
            name="idle"
            checked={v.idleMode === m.mode}
            disabled={!!busy}
            onchange={() => act("idle", () => S.SetIdleMode(m.mode))}
          />
          <span>
            <strong>{m.title}</strong>
            <span class="muted small block">{m.note}</span>
          </span>
        </label>
      {/each}
      <p class="muted small">.wslconfig: {v.wslConfigIdle ? "instanceIdleTimeout=-1 is set" : "no idle override"}</p>
      {#if v.needsWslReboot}
        {#if confirmApply}
          <div class="warn">
            <p>Shut down WSL now? Every WSL distro of your account stops, then starts again when used.</p>
            <button
              class="danger"
              onclick={() => {
                confirmApply = false;
                act("apply", () => S.ApplyWSLConfig());
              }}>Shut down WSL</button
            >
            <button onclick={() => (confirmApply = false)}>Later</button>
          </div>
        {:else}
          <div class="warn">
            <p>The .wslconfig change takes effect after WSL shuts down.</p>
            <button disabled={!!busy} onclick={() => (confirmApply = true)}>Apply now...</button>
          </div>
        {/if}
      {/if}
    </section>

    <section class="card">
      <h2>Browser for Open Portainer</h2>
      <label class="row">
        <input
          type="radio"
          name="browser"
          checked={v.browser === BrowserKind.BrowserDefault}
          onchange={() => act("browser", () => S.SetBrowser(BrowserKind.BrowserDefault, ""))}
        />
        System default
      </label>
      <label class="row">
        <input
          type="radio"
          name="browser"
          checked={v.browser === BrowserKind.BrowserChrome}
          disabled={!v.chromePath}
          onchange={() => act("browser", () => S.SetBrowser(BrowserKind.BrowserChrome, ""))}
        />
        Google Chrome {v.chromePath ? "" : "(not found)"}
      </label>
      <label class="row">
        <input
          type="radio"
          name="browser"
          checked={v.browser === BrowserKind.BrowserCustom}
          disabled={!v.browserPath}
          onchange={() => act("browser", () => S.SetBrowser(BrowserKind.BrowserCustom, v!.browserPath))}
        />
        Other: <span class="path">{v.browserPath || "none chosen"}</span>
        <button class="small-btn" onclick={pickBrowser}>Browse...</button>
      </label>
    </section>

    <section class="card">
      <h2>Certificate for https://localhost:9443</h2>
      {#if v.cert.present}
        <p class="small">
          Expires {expires}.
          {v.cert.trusted ? "Trusted by this Windows account." : "Not trusted: the browser shows a warning."}
        </p>
        <p class="muted small mono">{v.cert.thumbprint}</p>
      {:else}
        <p class="small">No certificate in Ubuntu-24.04 yet.</p>
      {/if}
      <div class="buttons">
        <button disabled={!!busy || !v.cert.present || v.cert.trusted} onclick={() => act("cert", () => S.TrustCert())}
          >Install (trust)</button
        >
        <button disabled={!!busy || !v.cert.trusted} onclick={() => act("cert", () => S.UntrustCert())}>Remove</button>
        <button disabled={!!busy} onclick={() => act("cert", () => S.RegenerateCert())}>Regenerate</button>
      </div>
      <p class="muted small">
        Windows asks you to confirm each change to your trusted certificates. The private key stays inside
        Ubuntu-24.04. Regenerate restarts Portainer.
      </p>
    </section>

    <Maintenance />
  {/if}

  {#if busy}
    <p class="muted small">Working...</p>
  {/if}
  {#if error}
    <p class="error">{error}</p>
  {/if}
  {#each v?.problems ?? [] as p (p)}
    <p class="error small">{p}</p>
  {/each}
</main>

<style>
  main {
    height: 100%;
    overflow-y: auto;
    padding: 16px 20px 24px;
    user-select: text;
  }
  h1 {
    margin: 0 0 12px;
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
    margin-bottom: 12px;
  }
  .row {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 3px 0;
  }
  .row.option {
    align-items: flex-start;
  }
  .indent {
    padding-left: 24px;
  }
  .block {
    display: block;
  }
  .muted {
    color: var(--text-muted);
  }
  .small {
    font-size: 12px;
  }
  .mono {
    font-family: Consolas, monospace;
    word-break: break-all;
  }
  .path {
    color: var(--text-muted);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 220px;
  }
  .buttons {
    display: flex;
    gap: 8px;
    margin: 8px 0;
  }
  button {
    padding: 5px 12px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--surface);
    cursor: pointer;
  }
  button:disabled {
    opacity: 0.45;
    cursor: default;
  }
  .small-btn {
    margin-left: auto;
    padding: 2px 10px;
  }
  .danger {
    background: var(--bad);
    border-color: var(--bad);
    color: #fff;
  }
  .warn {
    margin-top: 8px;
    padding: 8px 10px;
    border-radius: 6px;
    border: 1px solid var(--warn);
  }
  .warn p {
    margin: 0 0 8px;
  }
  .error {
    color: var(--bad);
  }
</style>
