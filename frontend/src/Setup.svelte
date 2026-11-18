<script lang="ts">
  import { Events } from "@wailsio/runtime";
  import * as S from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/setupservice";
  import * as A from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/adoptservice";
  import type { Report } from "../bindings/github.com/the-khiem7/docker-zero/internal/adopt/models";
  import type { SetupView } from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/models";
  import "./sx.css";
  import Ring from "./Ring.svelte";
  import Icon from "./Icon.svelte";
  import Hint from "./Hint.svelte";
  import ErrorNote from "./ErrorNote.svelte";
  import { run, service, live } from "./lib/status.svelte";

  const minLength = 12;

  let v = $state<SetupView | null>(null);
  let username = $state("admin");
  let secret = $state("");
  let again = $state("");
  let error = $state("");

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
      ((found.findings?.length ?? 0) > 0 ||
        (found.container !== "" && found.container !== "missing") ||
        // A stopped distro reports no container, so without this the card and
        // its result vanish once the last finding is replaced.
        adopted.length > 0),
  );


  // A step that has run this long gets a "still working" notice; the backend reports no timestamps,
  // so the clock starts when this window first sees the step running.
  const SLOW_AFTER_S = 120;
  const expect: Record<string, string> = {
    distro: "mostly downloading, several minutes",
    docker: "downloading packages, a few minutes",
  };

  let now = $state(Date.now());
  let since = $state({ id: "", at: Date.now() });

  const steps = $derived(v?.steps ?? []);
  const total = $derived(steps.length);
  const finished = $derived(steps.filter((s) => s.status === "done" || s.status === "skipped").length);
  const current = $derived(steps.find((s) => s.status === "running"));
  const failedStep = $derived(steps.find((s) => s.status === "failed"));
  const level = $derived(failedStep || v?.error ? "failed" : v?.done ? "running" : current ? "partial" : "stopped");

  $effect(() => {
    const id = current?.id ?? "";
    if (id !== since.id) since = { id, at: Date.now() };
    if (!id) return;
    const t = setInterval(() => (now = Date.now()), 1000);
    return () => clearInterval(t);
  });

  const elapsed = $derived(current ? Math.max(0, Math.floor((now - since.at) / 1000)) : 0);
  const clock = $derived(`${Math.floor(elapsed / 60)}:${String(elapsed % 60).padStart(2, "0")}`);

  const headline = $derived(
    v?.done ? "All set" : failedStep || v?.error ? "Setup stopped" : current ? "Setting up" : "Set up Docker Zero",
  );
  const busy = $derived(v?.running ?? false);
  const done = $derived(v?.done ?? false);
  const editable = $derived(!busy && !done);

  Events.On("setup", (e) => {
    v = e.data;
  });
  S.State().then((s) => (v = s));
  detect();
  window.addEventListener("focus", detect);

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

  const short = (t: string) => {
    const k = t.toLowerCase();
    if (k.includes("wsl")) return "WSL";
    if (k.includes("distro")) return "Ubuntu";
    if (k.includes("docker engine")) return "Docker";
    if (k.includes("linux user")) return "User";
    if (k.includes("certificate")) return "Certificate";
    if (k.includes("start portainer")) return "Portainer";
    if (k.includes("admin")) return "Admin";
    if (k.includes("docker command")) return "Command";
    return t;
  };
  const stepIcon = (t: string) =>
    ({ WSL: "window", Ubuntu: "distro", Docker: "app", User: "user", Certificate: "shield", Portainer: "globe", Admin: "key", Command: "terminal" })[short(t)] ?? "sliders";
  const itemIcon = (t: string) => {
    const k = t.toLowerCase();
    if (k.includes("shortcut") && k.includes("start")) return "power";
    if (k.includes("shortcut")) return "menu";
    if (k.includes("docker") || k.includes("command")) return "terminal";
    return "file";
  };
  const shortItem = (t: string) => {
    const k = t.toLowerCase();
    if (k.includes("portainer")) return "Portainer shortcuts";
    if (k.includes("sign-in") || k.includes("startup")) return "Auto start";
    if (k.includes("docker")) return "docker command";
    if (k.includes("manage") || k.includes("idle")) return ".wslconfig line";
    return t.length > 20 ? t.slice(0, 18) + "..." : t;
  };
  const ready = $derived(!!username.trim() && secret.length >= minLength && secret === again);
  const fill = $derived(Math.min(100, (secret.length / minLength) * 100));
</script>


<main class="sx">
  <section class="card w6 hero">
    <div class="dz">
      <Ring id="setup" orca parts={finished} max={Math.max(total, 1)} {level} size={104} />
      <div class="grow">
        <div class="big">{headline}</div>
        <div class="chips" style="margin-top: 8px">
          {#if total > 0}<span class="chip"><Icon name="check" size={14} />{finished} / {total} steps</span>{/if}
          {#if current}<span class="chip"><Icon name="wave" size={14} />{clock} on this step</span>{/if}
          {#if failedStep}<span class="chip bad"><Icon name="alert" size={14} />{short(failedStep.title)} stopped</span>{/if}
          {#if !current && !failedStep && !done}
            <span class="chip"><Icon name="shield" size={14} />no admin rights</span>
            <span class="chip"><Icon name="clock" size={14} />several minutes</span>
          {/if}
          {#if current && expect[current.id]}<span class="chip"><Icon name="clock" size={14} />{expect[current.id]}</span>{/if}
          {#if current && elapsed >= SLOW_AFTER_S}<span class="chip warn" role="status"><Icon name="alert" size={14} />Still working, nothing is stuck</span>{/if}
        </div>
      </div>
      {#if !current && !done}
        <Hint
          label="Setup"
          text="Installs Ubuntu-24.04 in WSL, Docker Engine and Portainer for your account. It needs no admin rights and takes several minutes, mostly for downloads."
        />
      {/if}
    </div>
  </section>

  {#if hasExisting && found}
    <section class="card w6">
      <div class="lab">
        <h2>Existing setup found</h2>
        <Hint
          label="Existing setup"
          text={`Ubuntu-24.04${found.container && found.container !== "missing" ? `, Portainer (${found.container})` : ""} already exists. Taking over keeps the distro, containers and Portainer data as they are and replaces the launchers below with Docker Zero's own. You do not need to run setup.`}
        />
      </div>
      {#if found.findings?.length}
        <div class="apps wrap">
          {#each found.findings ?? [] as f (f.item)}
            <label class="app" class:on={chosen[f.item]} title={f.detail}>
              <input class="pick" type="checkbox" aria-label={f.title} bind:checked={chosen[f.item]} disabled={adopting} />
              <span class="face"><Icon name={itemIcon(f.title)} size={24} /></span>{shortItem(f.title)}
            </label>
          {/each}
        </div>
        <div class="actions" style="justify-content: center">
          <button class="btn primary" disabled={adopting || !Object.values(chosen).some(Boolean)} onclick={takeOver}>
            {#if adopting}<span class="spin"></span>{/if}Take over
          </button>
        </div>
      {:else}
        <div class="chips"><span class="chip ok"><Icon name="check" size={14} />Nothing left to replace</span></div>
      {/if}
      {#if adopted.length}
        <div class="chips">
          <span class="chip ok"><Icon name="check" size={14} />Replaced {adopted.length} item{adopted.length === 1 ? "" : "s"}</span>
          <span class="chip"><Icon name="terminal" size={14} />open a new terminal for docker</span>
        </div>
      {/if}
      {#if adoptError}<div style="margin-top: 10px"><ErrorNote title="Take over stopped" detail={adoptError} /></div>{/if}
    </section>
  {/if}

  {#if editable}
    <section class="card w6">
      <div class="lab">
        <h2>Portainer admin account</h2>
        <Hint label="Admin account" text={`At least ${minLength} characters. Docker Zero never stores the password; Portainer keeps its own copy.`} />
      </div>
      <div class="resetrow">
        <div class="dcol">
          <button class="disc key" class:on={ready} aria-label="Start setup" style="--fill: {fill}%" onclick={start}>
            <Icon name="power" size={34} />
          </button>
          <span class="dcap">{v?.error ? "Try again" : "Start setup"}</span>
        </div>
        <div class="uf">
          <label class="ul"><Icon name="user" size={18} /><input type="text" aria-label="Username" placeholder="Username" bind:value={username} autocomplete="off" spellcheck="false" /></label>
          <label class="ul"><Icon name="lock" size={18} /><input type="password" aria-label="Password" placeholder="Password" bind:value={secret} autocomplete="new-password" /></label>
          <label class="ul"><Icon name="lock" size={18} /><input type="password" aria-label="Repeat password" placeholder="Repeat password" bind:value={again} autocomplete="new-password" /></label>
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

  {#if v && v.steps?.length}
    <section class="card w6">
      <div class="lab"><h2>Steps</h2></div>
      <ol class="stepgrid">
        {#each v.steps as step (step.id)}
          <li class="st {step.status || 'pending'}" title={step.title} aria-current={step.status === "running" ? "step" : undefined}>
            <span class="nd">
              {#if step.status === "done"}<Icon name="check" size={16} />
              {:else if step.status === "failed"}<Icon name="cross" size={16} />
              {:else if step.status === "skipped"}<Icon name="skip" size={16} />
              {:else}<Icon name={stepIcon(step.title)} size={15} />
              {/if}
            </span>
            <span class="cap">{short(step.title)}</span>
            <span class="sr">{step.title}: {step.status || "pending"}</span>
          </li>
        {/each}
      </ol>
      {#if current}<span class="indet" role="progressbar" aria-label={current.title}><i></i></span>{/if}
      {#if v.log?.length}
        <details>
          <summary><Icon name="terminal" size={14} />Details</summary>
          <pre class="log">{v.log.join("\n")}</pre>
        </details>
      {/if}
    </section>
  {/if}

  {#if error || v?.error}
    <div class="w6"><ErrorNote title={error && error.length <= 60 ? error : "Setup stopped"} detail={error || v?.error || ""} /></div>
  {/if}
  {#if live.error}
    <div class="w6"><ErrorNote title="Could not open Portainer" detail={live.error} /></div>
  {/if}

  {#if done}
    <section class="card w6">
      <div class="dz">
        <Ring id="done" orca parts={total} max={Math.max(total, 1)} level="running" size={72} />
        <div class="grow">
          <div class="big">Portainer is ready</div>
          <div class="chips" style="margin-top: 6px">
            <span class="chip"><Icon name="user" size={14} />sign in with your account</span>
            <span class="chip"><Icon name="terminal" size={14} />new terminals have docker</span>
          </div>
        </div>
        <button class="btn primary" onclick={open}><Icon name="globe" size={16} />Open Portainer</button>
      </div>
    </section>
  {/if}
</main>
