<script lang="ts">
  import "./sx.css";
  import Maintenance from "./Maintenance.svelte";
  import Icon from "./Icon.svelte";
  import Hint from "./Hint.svelte";
  import ErrorNote from "./ErrorNote.svelte";
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
      icon: "anchor",
      pos: 50,
      title: "Keep alive",
      scope: "This app only",
      note: "A hidden session holds Ubuntu-24.04 up. Changes nothing outside this app.",
    },
    {
      mode: IdleMode.IdleWSLConfig,
      icon: "infinity",
      pos: 94,
      title: "Never idle",
      scope: "All WSL distros",
      note: "Sets instanceIdleTimeout=-1 in .wslconfig. Applies to every WSL2 distro of your account after WSL restarts.",
    },
    {
      mode: IdleMode.IdleDefault,
      icon: "moon",
      pos: 6,
      title: "WSL default",
      scope: "Stops when idle",
      note: "WSL stops the distro, and Portainer with it, when nothing uses it.",
    },
  ];

  const expires = $derived(v?.cert.notAfter ? new Date(v.cert.notAfter).toLocaleDateString() : "");
  const stopPos = $derived(v?.idleMode === IdleMode.IdleDefault ? 6 : v?.idleMode === IdleMode.IdleKeepalive ? 50 : 94);
  const idleNote = $derived(idleModes.find((m) => m.mode === v?.idleMode)?.note ?? "");
  const thumb = $derived(v?.cert.thumbprint ? v.cert.thumbprint.slice(0, 8) + "..." + v.cert.thumbprint.slice(-4) : "");

  load();
  window.addEventListener("focus", load);
</script>

<main class="sx">
  <h1>Settings</h1>

  {#if !v}
    <span class="sr">Loading...</span>
    <div class="grid">
      <div class="skel w6"></div>
      <div class="skel w3"></div>
      <div class="skel w3"></div>
      <div class="skel w6"></div>
    </div>
  {:else}
    <div class="grid">
      <section class="card w6">
        <div class="lab">
          <h2>Start with Windows</h2>
          <Hint
            label="Start with Windows"
            text="Docker Zero opens when you sign in. With the second switch on it also starts Ubuntu-24.04 and Portainer."
          />
        </div>
        <div class="start">
          <label class="disc" class:on={v.autoStart.enabled}>
            <input
              class="pick"
              type="checkbox"
              aria-label="Start Docker Zero when I sign in"
              checked={v.autoStart.enabled}
              disabled={!!busy}
              onchange={(e) => act("autostart", () => S.SetAutoStart(e.currentTarget.checked, v!.autoStart.startDistro))}
            />
            <Icon name="power" size={36} />
          </label>
          <div>
            <div class="big">{v.autoStart.enabled ? "On at sign-in" : "Off"}</div>
            <label class="subline" class:off={!v.autoStart.enabled}>
              Ubuntu + Portainer too
              <input
                class="pick"
                type="checkbox"
                aria-label="Also start Ubuntu-24.04 and Portainer"
                checked={v.autoStart.startDistro}
                disabled={!!busy || !v.autoStart.enabled}
                onchange={(e) => act("autostart", () => S.SetAutoStart(true, e.currentTarget.checked))}
              />
              <span class="sw" class:on={v.autoStart.startDistro}></span>
            </label>
          </div>
        </div>
      </section>

      <section class="card w3">
        <div class="lab">
          <h2>Shortcuts</h2>
          <Hint label="Shortcuts" text={"Shortcuts and auto start point at " + v.installPath + "; the app copies itself there first."} />
        </div>
        <div class="apps">
          <label class="app" class:on={v.shortcuts.desktop}>
            <input
              class="pick"
              type="checkbox"
              aria-label="Desktop"
              checked={v.shortcuts.desktop}
              disabled={!!busy}
              onchange={(e) => act("shortcuts", () => S.SetShortcuts(e.currentTarget.checked, v!.shortcuts.startMenu))}
            />
            <span class="face"><Icon name="monitor" size={26} /></span>Desktop
          </label>
          <label class="app" class:on={v.shortcuts.startMenu}>
            <input
              class="pick"
              type="checkbox"
              aria-label="Start Menu"
              checked={v.shortcuts.startMenu}
              disabled={!!busy}
              onchange={(e) => act("shortcuts", () => S.SetShortcuts(v!.shortcuts.desktop, e.currentTarget.checked))}
            />
            <span class="face"><Icon name="menu" size={26} /></span>Start
          </label>
        </div>
      </section>

      <section class="card w3">
        <div class="lab"><h2>Open Portainer in</h2></div>
        <div class="rounds">
          <label class="rb" class:on={v.browser === BrowserKind.BrowserDefault}>
            <input
              class="pick"
              type="radio"
              name="browser"
              aria-label="System default"
              checked={v.browser === BrowserKind.BrowserDefault}
              onchange={() => act("browser", () => S.SetBrowser(BrowserKind.BrowserDefault, ""))}
            />
            <span class="o"><Icon name="globe" size={20} /></span>Default
          </label>
          <label class="rb" class:on={v.browser === BrowserKind.BrowserChrome} class:off={!v.chromePath}>
            <input
              class="pick"
              type="radio"
              name="browser"
              aria-label="Google Chrome"
              checked={v.browser === BrowserKind.BrowserChrome}
              disabled={!v.chromePath}
              onchange={() => act("browser", () => S.SetBrowser(BrowserKind.BrowserChrome, ""))}
            />
            <span class="o"><Icon name="compass" size={20} /></span>{v.chromePath ? "Chrome" : "no Chrome"}
          </label>
          <label class="rb" class:on={v.browser === BrowserKind.BrowserCustom} class:off={!v.browserPath}>
            <input
              class="pick"
              type="radio"
              name="browser"
              aria-label="Other browser"
              checked={v.browser === BrowserKind.BrowserCustom}
              disabled={!v.browserPath}
              onchange={() => act("browser", () => S.SetBrowser(BrowserKind.BrowserCustom, v!.browserPath))}
            />
            <span class="o"><Icon name="window" size={20} /></span>Other
          </label>
        </div>
        <div class="chips">
          {#if v.browserPath}
            <span class="chip mono" title={v.browserPath}><Icon name="folder" size={14} /><span class="clip"><bdi>{v.browserPath}</bdi></span></span>
          {/if}
          <button class="btn small" onclick={pickBrowser}><Icon name="folder" size={14} />Browse...</button>
        </div>
      </section>

      <section class="card w6">
        <div class="lab">
          <h2>Keep Ubuntu-24.04 running</h2>
          <Hint label="Keep Ubuntu-24.04 running" text={idleNote} />
        </div>
        <div class="tide">
          <svg viewBox="0 0 400 20" preserveAspectRatio="none" aria-hidden="true">
            <path d="M0 10Q25 0 50 10T100 10T150 10T200 10T250 10T300 10T350 10T400 10" fill="none" stroke="var(--line)" stroke-width="4" vector-effect="non-scaling-stroke" />
            <path
              d="M0 10Q25 0 50 10T100 10T150 10T200 10T250 10T300 10T350 10T400 10"
              fill="none"
              stroke="var(--teal)"
              stroke-width="4"
              vector-effect="non-scaling-stroke"
              style="clip-path: inset(-4px {100 - stopPos}% -4px 0)"
            />
          </svg>
          {#each [...idleModes].sort((a, b) => a.pos - b.pos) as m (m.mode)}
            <label class="stop" class:on={v.idleMode === m.mode} style="left: {m.pos}%">
              <input
                class="pick"
                type="radio"
                name="idle"
                aria-label={m.title}
                checked={v.idleMode === m.mode}
                disabled={!!busy}
                onchange={() => act("idle", () => S.SetIdleMode(m.mode))}
              />
              <span class="dot"><Icon name={m.icon} size={m.pos === stopPos ? 22 : 18} /></span>{m.title}
            </label>
          {/each}
        </div>
        <div class="scope"><span>stops when idle</span><span>this app only</span><span>all distros</span></div>
        <div class="chips">
          <span class="chip" class:ok={v.wslConfigIdle}><Icon name="file" size={14} />.wslconfig {v.wslConfigIdle ? "idle override on" : "no idle override"}</span>
          {#if v.needsWslReboot}
            {#if confirmApply}
              <span class="chip bad"><Icon name="alert" size={14} />Shut down WSL now? Every distro stops</span>
              <button
                class="btn danger small"
                onclick={() => {
                  confirmApply = false;
                  act("apply", () => S.ApplyWSLConfig());
                }}><Icon name="power" size={14} />Shut down WSL</button
              >
              <button class="btn small" onclick={() => (confirmApply = false)}>Later</button>
            {:else}
              <span class="chip warn"><Icon name="wave" size={14} />Restart WSL to apply</span>
              <button class="btn small" disabled={!!busy} onclick={() => (confirmApply = true)}>Apply now...</button>
            {/if}
          {/if}
        </div>
      </section>

      <section class="card w4">
        <div class="lab">
          <h2>Certificate</h2>
          <Hint
            label="Certificate"
            text="Windows asks you to confirm each change to your trusted certificates. The private key stays inside Ubuntu-24.04. Regenerate restarts Portainer."
          />
        </div>
        {#if v.cert.present}
          <div class="stamp" class:warn={!v.cert.trusted}>
            <svg class="arc" viewBox="0 0 108 108" aria-hidden="true">
              <defs><path id="ca" d="M17 54a37 37 0 0 1 74 0" /><path id="cb" d="M15 58a39 39 0 0 0 78 0" /></defs>
              <text font-size="10" font-weight="700" letter-spacing="2.5" fill="currentColor"
                ><textPath href="#ca" startOffset="50%" text-anchor="middle">{v.cert.trusted ? "TRUSTED" : "NOT TRUSTED"}</textPath></text
              >
              <text font-size="9" fill="currentColor"><textPath href="#cb" startOffset="50%" text-anchor="middle">until {expires}</textPath></text>
            </svg>
            <Icon name={v.cert.trusted ? "shield" : "shieldAlert"} size={34} />
          </div>
          <div class="stampbar"><span class="chip mono" title={v.cert.thumbprint}><Icon name="hash" size={12} />{thumb}</span></div>
        {:else}
          <div class="stamp none"><Icon name={v.certUnread ? "moon" : "shieldOff"} size={34} /></div>
          <div class="stampbar">
            <span class="chip"><Icon name={v.certUnread ? "power" : "alert"} size={14} />{v.certUnread ? "Start Ubuntu-24.04 to read it" : "No certificate yet"}</span>
          </div>
        {/if}
        <div class="rounds" style="margin-top: 12px">
          <button class="rbtn" aria-label="Install (trust)" disabled={!!busy || !v.cert.present || v.cert.trusted} onclick={() => act("cert", () => S.TrustCert())}
            ><span class="o"><Icon name="shield" size={20} /></span>Trust</button
          >
          <button class="rbtn" aria-label="Remove" disabled={!!busy || !v.cert.trusted} onclick={() => act("cert", () => S.UntrustCert())}
            ><span class="o"><Icon name="shieldOff" size={20} /></span>Remove</button
          >
          <button class="rbtn" aria-label="Regenerate" disabled={!!busy || v.certUnread} onclick={() => act("cert", () => S.RegenerateCert())}
            ><span class="o"><Icon name="refresh" size={20} /></span>Renew</button
          >
        </div>
      </section>

      <Maintenance />
    </div>
  {/if}

  {#if busy}
    <div class="toast"><span class="spin"></span>Working...</div>
  {/if}
  {#if error}
    <div style="margin-top: 10px"><ErrorNote title="Could not apply the change" detail={error} /></div>
  {/if}
  {#each v?.problems ?? [] as p (p)}
    <div style="margin-top: 8px"><ErrorNote title="Could not read a setting" detail={p} /></div>
  {/each}
</main>
