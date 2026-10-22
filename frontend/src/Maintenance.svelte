<script lang="ts">
  import { Events } from "@wailsio/runtime";
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

  const icons: Record<string, string> = { running: "◐", done: "✓", skipped: "↷", failed: "✕" };

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
</script>

<section class="card">
  <h2>Portainer</h2>
  <div class="buttons">
    <button disabled={working} onclick={update}>Update Portainer</button>
    <button disabled={working} onclick={() => (resetOpen = !resetOpen)}>Reset admin password...</button>
  </div>
  <p class="muted small">
    Update pulls the latest long-term-support image and restarts Portainer; your data stays. Reset replaces the
    password of the existing admin account.
  </p>
  {#if resetOpen}
    <div class="panel">
      <label class="field">
        Admin username
        <input type="text" bind:value={username} autocomplete="off" spellcheck="false" />
      </label>
      <label class="field">
        New password
        <input type="password" bind:value={secret} autocomplete="new-password" />
      </label>
      <label class="field">
        Repeat new password
        <input type="password" bind:value={again} autocomplete="new-password" />
      </label>
      <p class="muted small">
        At least {minLength} characters. The username is only used to check the new password; the account keeps its
        name. Portainer stops for a few seconds.
      </p>
      <button class="primary" disabled={working} onclick={reset}>Reset password</button>
    </div>
  {/if}
  {#if busy === "update" || busy === "reset"}<p class="muted small">Working...</p>{/if}
  {#if notice}<p class="ok small">{notice}</p>{/if}
  {#if error && busy === ""}<p class="error small">{error}</p>{/if}
</section>

<section class="card danger-zone">
  <h2>Danger zone</h2>
  {#if u && (removing || removed || u.error)}
    <ul class="steps">
      {#each u.steps as step (step.id)}
        <li class={step.status || "pending"}>
          <span class="icon" aria-hidden="true">{icons[step.status] ?? "○"}</span>
          <span>{step.title}</span>
          {#if step.status === "skipped"}<span class="note">nothing to do</span>{/if}
        </li>
      {/each}
    </ul>
    {#if u.log?.length}<pre class="log">{u.log.join("\n")}</pre>{/if}
    {#if u.error}<p class="error small">{u.error}</p>{/if}
    {#if removed}
      <p class="small">
        Docker Zero is removed. Quit to finish: its own folder is deleted once the app has closed. A new terminal no
        longer finds the docker command.
      </p>
      <button class="danger" disabled={!!busy} onclick={quit}>Quit Docker Zero</button>
    {:else if !removing}
      <p class="muted small">Fix the problem above and run the uninstall again; finished steps are skipped.</p>
      <button disabled={working} onclick={() => (stage = "plan")}>Uninstall all...</button>
    {/if}
  {:else if stage === "idle"}
    <p class="muted small">Removes Ubuntu-24.04 with Docker, Portainer and every container and volume in it.</p>
    <button class="danger-outline" disabled={working} onclick={() => (stage = "plan")}>Uninstall all...</button>
  {/if}

  {#if stage === "plan" && !removing && !removed}
    <div class="panel">
      <p class="small"><strong>This removes:</strong></p>
      <ul class="small list">
        <li>the Ubuntu-24.04 distro with all its images, containers and volumes</li>
        <li>the trusted Portainer certificate (Windows asks you to confirm)</li>
        <li>auto start, shortcuts, the docker command and its PATH entry</li>
        <li>the .wslconfig idle line this app added, and the Docker Zero settings</li>
      </ul>
      <label class="row">
        <input type="checkbox" bind:checked={backup} />
        Save a backup of Ubuntu-24.04 first
      </label>
      {#if backup}
        <p class="warn small">
          Making the backup stops Ubuntu-24.04 and every container in it. The file can be as large as the whole
          distro.
        </p>
        <div class="row">
          <span class="path small">{backupPath || "no file chosen"}</span>
          <button class="small-btn" onclick={pick}>Choose file...</button>
        </div>
      {/if}
      <div class="buttons">
        <button disabled={backup && !backupPath} onclick={() => (stage = "confirm")}>Continue</button>
        <button onclick={() => (stage = "idle")}>Cancel</button>
      </div>
    </div>
  {/if}

  {#if stage === "confirm" && !removing && !removed}
    <div class="panel warn-panel">
      <p class="small">
        <strong>This cannot be undone.</strong>
        {backup ? "Your backup is saved first." : "No backup will be made."} Remove Docker Zero and Ubuntu-24.04 now?
      </p>
      <div class="buttons">
        <button class="danger" disabled={working} onclick={begin}>Uninstall now</button>
        <button onclick={() => (stage = "idle")}>Cancel</button>
      </div>
    </div>
  {/if}

  {#if busy === "uninstall" || removing}<p class="muted small">Working...</p>{/if}
  {#if error && busy === "" && stage !== "idle"}<p class="error small">{error}</p>{/if}
</section>

<style>
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
  .danger-zone {
    border-color: var(--bad);
  }
  .panel {
    margin-top: 8px;
    padding: 8px 10px;
    border: 1px solid var(--border);
    border-radius: 6px;
  }
  .warn-panel {
    border-color: var(--bad);
  }
  .muted {
    color: var(--text-muted);
  }
  .small {
    font-size: 12px;
  }
  .ok {
    color: var(--ok);
  }
  .error {
    color: var(--bad);
    word-break: break-word;
  }
  .warn {
    color: var(--warn);
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
  .row {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 3px 0;
  }
  .list {
    margin: 0 0 8px;
    padding-left: 18px;
  }
  .path {
    color: var(--text-muted);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
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
    padding: 2px 10px;
  }
  .primary {
    border-color: var(--accent);
    background: var(--accent);
    color: var(--accent-text);
  }
  .danger {
    background: var(--bad);
    border-color: var(--bad);
    color: #fff;
  }
  .danger-outline {
    border-color: var(--bad);
    color: var(--bad);
  }
  .steps {
    list-style: none;
    margin: 0 0 8px;
    padding: 0;
  }
  .steps li {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 3px 0;
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
  .icon {
    width: 16px;
    text-align: center;
  }
  .note {
    margin-left: auto;
    font-size: 12px;
  }
  .log {
    max-height: 140px;
    overflow: auto;
    margin: 8px 0;
    padding: 8px 10px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    font-size: 11px;
    white-space: pre-wrap;
    word-break: break-word;
  }
</style>
