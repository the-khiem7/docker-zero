<script lang="ts">
  import { Events } from "@wailsio/runtime";
  import Icon from "./Icon.svelte";
  import Hint from "./Hint.svelte";
  import Ring from "./Ring.svelte";
  import ErrorNote from "./ErrorNote.svelte";
  import * as M from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/maintainservice";
  import * as U from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/uninstallservice";
  import type { SetupView } from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/models";

  const minLength = 12;

  type Stage = "idle" | "plan" | "confirm";

  let busy = $state("");
  let error = $state("");
  let notice = $state("");

  let resetOpen = $state(false);
  let username = $state("admin");
  let secret = $state("");
  let again = $state("");

  let stage = $state<Stage>("idle");
  let backup = $state(true);
  let backupPath = $state("");
  let u = $state<SetupView | null>(null);

  const removing = $derived(u?.running ?? false);
  const removed = $derived(u?.done ?? false);
  const working = $derived(!!busy || removing);

  const subject = (title: string) => {
    const t = title.toLowerCase();
    if (t.includes("back up")) return "backup";
    if (t.includes("certificate")) return "shield";
    if (t.includes("distro")) return "distro";
    if (t.includes("auto start")) return "power";
    return "sliders";
  };

  Events.On("uninstall", (e) => {
    u = e.data;
  });
  U.State().then((s) => {
    if (s.running || s.done || s.error) u = s;
  });

  function message(err: unknown): string {
    return err instanceof Error ? err.message : String(err);
  }

  async function act(name: string, fn: () => Promise<void>) {
    busy = name;
    error = "";
    notice = "";
    try {
      await fn();
    } catch (err) {
      error = message(err);
    } finally {
      busy = "";
    }
  }

  const update = () =>
    act("update", async () => {
      const replaced = await M.UpdatePortainer();
      notice = replaced ? "Portainer updated and restarted." : "Portainer already runs the latest image.";
    });

  const reset = () => {
    if (!username.trim()) {
      error = "Enter the admin username.";
    } else if (secret.length < minLength) {
      error = `The password needs at least ${minLength} characters.`;
    } else if (secret !== again) {
      error = "The two passwords differ.";
    } else {
      return act("reset", async () => {
        await M.ResetAdmin(username.trim(), secret);
        secret = "";
        again = "";
        resetOpen = false;
        notice = "Password replaced. Sign in to Portainer with the new password.";
      });
    }
  };

  async function pick() {
    const path = await U.PickBackup();
    if (path) backupPath = path;
  }

  const begin = () =>
    act("uninstall", async () => {
      await U.Start(backup ? backupPath : "");
      stage = "idle";
    });

  const quit = () => act("quit", () => U.Quit());

  const removes = [
    { icon: "distro", t: "Ubuntu", s: "Ubuntu-24.04 with images, containers, volumes" },
    { icon: "shield", t: "Certificate", s: "the trusted Portainer certificate" },
    { icon: "power", t: "Auto start", s: "auto start, shortcuts and the PATH entry" },
    { icon: "terminal", t: "docker", s: "the docker command" },
    { icon: "file", t: ".wslconfig", s: "the idle line this app added" },
    { icon: "sliders", t: "Settings", s: "Docker Zero settings and data folder" },
  ];

  const short = (title: string) => {
    const k = subject(title);
    return ({ backup: "Backup", shield: "Certificate", distro: "Ubuntu", power: "Auto start", sliders: "Settings" } as Record<string, string>)[k] ?? title;
  };
  const finished = $derived((u?.steps ?? []).filter((s) => s.status === "done" || s.status === "skipped").length);
  const trackPct = $derived(Math.min(100, (Math.max(0, finished - (removing ? 0 : 1)) / Math.max(1, (u?.steps?.length ?? 2) - 1)) * 100));
  const ready = $derived(!!username.trim() && secret.length >= minLength && secret === again);
  const fill = $derived(Math.min(100, (secret.length / minLength) * 100));
</script>

<section class="card w2">
  <div class="lab">
    <h2>Portainer</h2>
    <Hint
      label="Portainer"
      text="Update pulls the latest long-term-support image and restarts Portainer; your data stays. Reset replaces the password of the existing admin account."
    />
  </div>
  <div class="col">
    <button class="sq" aria-label="Update Portainer" disabled={working} onclick={update}>
      {#if busy === "update"}<span class="spin big"></span>{:else}<Icon name="update" size={24} />{/if}Update
    </button>
    <button class="sq" aria-label="Reset admin password..." disabled={working} onclick={() => (resetOpen = !resetOpen)}
      ><Icon name="key" size={24} />Reset password</button
    >
  </div>
</section>

{#if resetOpen}
  <section class="card w6">
    <div class="lab">
      <h2>Reset admin password</h2>
      <Hint
        label="Reset password"
        text={`At least ${minLength} characters. The username is only used to check the new password; the account keeps its name. Portainer stops for a few seconds.`}
      />
    </div>
    <div class="resetrow">
      <button class="disc key" class:on={ready} aria-label="Reset password" disabled={working} style="--fill: {fill}%" onclick={reset}>
        {#if busy === "reset"}<span class="spin big"></span>{:else}<Icon name="key" size={34} />{/if}
      </button>
      <div class="uf">
        <label class="ul"><Icon name="user" size={18} /><input type="text" aria-label="Admin username" placeholder="Admin username" bind:value={username} autocomplete="off" spellcheck="false" /></label>
        <label class="ul"><Icon name="lock" size={18} /><input type="password" aria-label="New password" placeholder="New password" bind:value={secret} autocomplete="new-password" /></label>
        <label class="ul"><Icon name="lock" size={18} /><input type="password" aria-label="Repeat new password" placeholder="Repeat new password" bind:value={again} autocomplete="new-password" /></label>
        <div class="meter">
          <span class="ticks">{#each Array.from({ length: minLength }, (_, i) => i) as i (i)}<i class:on={i < secret.length}></i>{/each}</span>
          <span class="chip" class:ok={secret.length >= minLength}>{Math.min(secret.length, minLength)}/{minLength}</span>
          {#if again}<span class="chip" class:ok={again === secret} class:bad={again !== secret}
              ><Icon name={again === secret ? "check" : "cross"} size={14} />{again === secret ? "match" : "differ"}</span
            >{/if}
        </div>
      </div>
    </div>
  </section>
{/if}
{#if error && busy === ""}
  <div class="w6"><ErrorNote title={error.length > 60 ? "Portainer task failed" : error} detail={error} /></div>
{/if}

<section class="card w6 red">
  <div class="lab">
    <h2>Danger zone</h2>
    <Hint
      label="Danger zone"
      text="Uninstall removes the Ubuntu-24.04 distro with all its images, containers and volumes; the trusted Portainer certificate (Windows asks you to confirm); auto start, shortcuts, the docker command and its PATH entry; the .wslconfig idle line this app added; and the Docker Zero settings."
    />
  </div>

  {#if u && (removing || removed || u.error)}
    {#if removed}
      <div class="dz">
        <Ring id="gone" orca parts={0} level="stopped" size={72} />
        <div class="grow">
          <div class="big">Docker Zero is removed</div>
          <span class="chip"><Icon name="terminal" size={14} />new terminals lose the docker command</span>
        </div>
        <button class="btn danger" disabled={!!busy} onclick={quit}><Icon name="power" size={16} />Quit Docker Zero</button>
      </div>
    {:else}
      <div class="track" style="--p: {trackPct}%">
        <span class="rail"><i></i></span>
        {#each u.steps as step (step.id)}
          <div class="st {step.status || 'pending'}" title={step.title}>
            <span class="nd">
              {#if step.status === "done"}<Icon name="check" size={16} />
              {:else if step.status === "failed"}<Icon name="cross" size={16} />
              {:else if step.status === "skipped"}<Icon name="skip" size={16} />
              {:else}<Icon name={subject(step.title)} size={15} />
              {/if}
            </span>
            <span class="cap">{short(step.title)}</span>
            <span class="sr">{step.title}: {step.status || "pending"}</span>
          </div>
        {/each}
      </div>
      <div class="chips" style="justify-content: space-between">
        {#if u.log?.length}
          <details open={removing}>
            <summary><Icon name="terminal" size={14} />Details</summary>
            <pre class="log">{u.log.join("\n")}</pre>
          </details>
        {:else}<span></span>{/if}
        {#if !removing}
          <button class="btn ghost-danger small" disabled={working} onclick={() => (stage = "plan")}><Icon name="refresh" size={14} />Uninstall all...</button>
        {/if}
      </div>
    {/if}
    {#if u.error}
      <div style="margin-top: 10px"><ErrorNote title="Uninstall stopped" detail={u.error} /></div>
    {/if}
  {:else if stage === "idle"}
    <div class="dz">
      <button class="disc red" aria-label="Uninstall all..." disabled={working} onclick={() => (stage = "plan")}><Icon name="trash" size={32} /></button>
      <div class="grow">
        <div class="big">Uninstall everything</div>
        <div class="chips" style="margin-top: 6px">
          <span class="chip"><Icon name="distro" size={14} />Ubuntu-24.04</span>
          <span class="chip"><Icon name="app" size={14} />Docker + Portainer</span>
          <span class="chip bad"><Icon name="trash" size={14} />volumes</span>
        </div>
      </div>
    </div>
  {/if}

  {#if stage === "plan" && !removing && !removed}
    <div class="cargo" aria-label="What is removed">
      {#each removes as r (r.t)}
        <span class="cg" title={r.s}><span class="bub"><Icon name={r.icon} size={22} /></span>{r.t}</span>
      {/each}
    </div>
    <div class="planrow">
      <div class="bkcell" class:on={backup}>
        <label class="bkhead">
          <input class="pick" type="checkbox" aria-label="Save a backup of Ubuntu-24.04 first" bind:checked={backup} />
          <span class="bub on"><Icon name="backup" size={22} /></span>
          <span class="grow"><b>Backup first</b></span>
          <span class="sw" class:on={backup}></span>
        </label>
        {#if backup}
          <div class="chips" style="margin-top: 8px">
            <span class="chip mono" title={backupPath}><Icon name="file" size={14} /><span class="clip"><bdi>{backupPath || "no file chosen"}</bdi></span></span>
            <button class="btn small" onclick={pick}><Icon name="folder" size={14} />Choose file...</button>
            <Hint label="Backup" text="The backup stops Ubuntu-24.04 and every container in it. The file can be as large as the whole distro." />
          </div>
        {/if}
      </div>
      <div class="col2">
        <button class="btn primary" disabled={backup && !backupPath} onclick={() => (stage = "confirm")}>Continue</button>
        <button class="btn" onclick={() => (stage = "idle")}>Cancel</button>
      </div>
    </div>
  {/if}

  {#if stage === "confirm" && !removing && !removed}
    <div class="dz">
      <button class="disc red on" aria-label="Uninstall now" disabled={working} onclick={begin}><Icon name="trash" size={32} /></button>
      <div class="grow">
        <div class="big">This cannot be undone</div>
        <div class="chips" style="margin-top: 6px">
          <span class="chip" class:ok={backup}><Icon name={backup ? "backup" : "cross"} size={14} />{backup ? "backup first" : "no backup"}</span>
          <span class="chip bad">tap the disc to remove</span>
        </div>
      </div>
      <button class="btn" onclick={() => (stage = "idle")}>Cancel</button>
    </div>
  {/if}

  {#if error && busy === "" && stage !== "idle"}
    <div style="margin-top: 10px"><ErrorNote title="Could not start" detail={error} /></div>
  {/if}
</section>

{#if busy === "update" || busy === "reset" || busy === "uninstall" || removing || notice}
  <div class="toast" class:good={!!notice && !busy && !removing}>
    {#if notice && !busy && !removing}<Icon name="check" size={16} />{notice}{:else}<span class="spin"></span>Working...{/if}
  </div>
{/if}
