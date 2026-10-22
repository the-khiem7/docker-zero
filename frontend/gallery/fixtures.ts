// Scenarios for the UI gallery: each one is the data the Go services would
// return in one state of the app. The page picks one with ?state=<name>.
import type { SettingsView, SetupView, StepView, View } from "../bindings/github.com/the-khiem7/docker-zero/internal/ui/models";
import type { Report } from "../bindings/github.com/the-khiem7/docker-zero/internal/adopt/models";
import { Item } from "../bindings/github.com/the-khiem7/docker-zero/internal/adopt/models";
import { BrowserKind, IdleMode } from "../bindings/github.com/the-khiem7/docker-zero/internal/settings/models";
import { Status, StepID } from "../bindings/github.com/the-khiem7/docker-zero/internal/setup/models";
import { State } from "../bindings/github.com/the-khiem7/docker-zero/internal/wsl/models";
import { caller, stateName, type Behaviour, type Mocked } from "./mock-kit";

export interface Scenario extends Behaviour {
  view?: Partial<View>;
  settings?: Partial<SettingsView>;
  setup?: SetupView;
  uninstall?: SetupView;
  report?: Report;
  /** Report returned by Detect after a successful Adopt. */
  reportAfterAdopt?: Report;
  updateReplaced?: boolean;
  /** Path the file pickers return. */
  pick?: string;
}

const checkedAt = "2026-10-22T09:00:00Z";

export const runningView: View = {
  wslAvailable: true,
  distro: State.StateRunning,
  docker: "active",
  portainer: "running",
  portainerUp: true,
  portainerVersion: "2.45.3",
  checkedAt,
  busy: "",
  holding: true,
  level: "running",
  corner: false,
};

const stoppedView: Partial<View> = {
  distro: State.StateStopped,
  docker: "unknown",
  portainer: "unknown",
  portainerUp: false,
  portainerVersion: "",
  holding: false,
  level: "stopped",
};

const notInstalledView: Partial<View> = { ...stoppedView, distro: State.StateNotInstalled, level: "failed" };

export const baseSettings: SettingsView = {
  autoStart: { enabled: true, startDistro: true },
  shortcuts: { desktop: true, startMenu: true },
  idleMode: IdleMode.IdleKeepalive,
  wslConfigIdle: false,
  browser: BrowserKind.BrowserDefault,
  browserPath: "",
  chromePath: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
  cert: {
    present: true,
    thumbprint: "3F2A9C41D07B6E58A1C4F09D2E7B835C61A0F4E2",
    notAfter: "2031-10-21T00:00:00Z",
    trusted: true,
  },
  certUnread: false,
  installPath: "C:\\Users\\you\\AppData\\Local\\Programs\\docker-zero\\docker-zero.exe",
  problems: null,
  wslConfigPath: "C:\\Users\\you\\.wslconfig",
  needsWslReboot: false,
};

const setupTitles: [StepID, string][] = [
  [StepID.StepWSL, "Check WSL"],
  [StepID.StepDistro, "Install the Ubuntu distro"],
  [StepID.StepDocker, "Install Docker Engine"],
  [StepID.StepUser, "Create the Linux user"],
  [StepID.StepCert, "Create the local certificate"],
  [StepID.StepPortainer, "Start Portainer"],
  [StepID.StepAdmin, "Create the Portainer admin"],
  [StepID.StepWrapper, "Add the docker command"],
];

const uninstallTitles: [StepID, string][] = [
  [StepID.StepBackup, "Back up the Ubuntu distro"],
  [StepID.StepTrust, "Remove the trusted certificate"],
  [StepID.StepRemove, "Remove the Ubuntu distro"],
  [StepID.StepWindows, "Remove auto start, shortcuts and the docker command"],
  [StepID.StepData, "Remove Docker Zero settings"],
];

/** Steps with the given statuses, in order; missing ones are pending. */
function steps(titles: [StepID, string][], statuses: Status[] = []): StepView[] {
  return titles.map(([id, title], i) => ({ id, title, status: statuses[i] ?? Status.$zero }));
}

const { StatusDone: D, StatusRunning: R, StatusSkipped: S, StatusFailed: F } = Status;

const setupLog = [
  "Checking WSL",
  "Downloading and registering Ubuntu-24.04 as Ubuntu-24.04",
  "Installing, this may take a few minutes...",
  "Ubuntu-24.04 registered",
  "Hit:1 http://archive.ubuntu.com/ubuntu noble InRelease",
  "Get:2 https://download.docker.com/linux/ubuntu noble InRelease [48.5 kB]",
  "Setting up docker-ce (5:29.9.2-1~ubuntu.24.04~noble) ...",
];

const freshSetup: SetupView = { steps: steps(setupTitles), log: null, running: false, done: false };

const legacyFindings = [
  {
    item: Item.ItemShortcuts,
    title: "Replace the Portainer shortcuts with Docker Zero shortcuts",
    detail: "Portainer.lnk on the Desktop and in the Start Menu",
  },
  {
    item: Item.ItemStartup,
    title: "Replace the sign-in start shortcut with Docker Zero's auto start",
    detail: "Startup\\keep-ubuntu-running.lnk",
  },
  {
    item: Item.ItemCommand,
    title: "Replace the manual docker command with Docker Zero's",
    detail: "C:\\Users\\you\\docker-wsl",
  },
  {
    item: Item.ItemIdle,
    title: "Let Docker Zero manage the .wslconfig idle line",
    detail: "C:\\Users\\you\\.wslconfig",
  },
];

const installed: Report = { distro: "running", container: "running", certificate: true, findings: [] };
const noDistro: Report = { distro: "not-installed", container: "", certificate: false, findings: [] };

export const scenarios: Record<string, Scenario> = {
  // Popup and corner box.
  loading: { view: { checkedAt: "" } },
  running: {},
  starting: { view: { busy: "start", level: "partial", portainerUp: false, portainerVersion: "" } },
  partial: { view: { portainer: "exited", portainerUp: false, portainerVersion: "", level: "partial" } },
  stopping: { view: { busy: "stop" } },
  stopped: { view: stoppedView },
  "not-installed": { view: notInstalledView },
  "no-wsl": {
    view: {
      ...notInstalledView,
      wslAvailable: false,
      error: 'WSL is not available: exec: "wsl.exe": executable file not found in %PATH%',
    },
  },
  "start-error": { view: stoppedView, fail: { "Service.Start": "timed out waiting for Portainer" } },
  "not-installed-start": {
    view: notInstalledView,
    fail: { "Service.Start": "Ubuntu-24.04 is not installed; run setup first" },
  },

  // Settings.
  "settings-loading": { hang: ["SettingsService.Get"] },
  settings: {},
  "settings-untrusted": { settings: { cert: { ...baseSettings.cert, trusted: false } } },
  "settings-no-cert": { settings: { cert: { present: false, trusted: false } } },
  "settings-cert-unread": { view: stoppedView, settings: { cert: { present: false, trusted: false }, certUnread: true } },
  "settings-wslconfig": {
    settings: { idleMode: IdleMode.IdleWSLConfig, wslConfigIdle: true, needsWslReboot: true },
  },
  "settings-browser": {
    settings: {
      browser: BrowserKind.BrowserCustom,
      browserPath: "C:\\Program Files\\Mozilla Firefox\\firefox.exe",
      chromePath: "",
    },
  },
  "settings-off": {
    settings: {
      autoStart: { enabled: false, startDistro: false },
      shortcuts: { desktop: false, startMenu: false },
      idleMode: IdleMode.IdleDefault,
    },
  },
  "settings-busy": { hang: ["SettingsService.SetShortcuts"] },
  "settings-problems": {
    settings: {
      problems: [
        "auto start: open HKCU Run key: Access is denied.",
        ".wslconfig: open C:\\Users\\you\\.wslconfig: The process cannot access the file because it is being used by another process.",
      ],
    },
    fail: { "SettingsService.SetShortcuts": "create shortcut: Access is denied." },
  },

  // Maintenance, inside the settings window.
  maint: {},
  "maint-updating": { hang: ["MaintainService.UpdatePortainer"] },
  "maint-updated": { updateReplaced: true },
  "maint-current": { updateReplaced: false },
  "maint-update-failed": {
    fail: {
      "MaintainService.UpdatePortainer":
        "pull portainer/portainer-ce:lts: Error response from daemon: Get \"https://registry-1.docker.io/v2/\": dial tcp 127.0.0.9:443: connect: connection refused",
    },
  },
  "maint-reset-wrong-user": {
    fail: {
      "MaintainService.ResetAdmin":
        'the password was replaced, but signing in as "nobody" failed; check the username: sign-in answered HTTP 422',
    },
  },
  "maint-reset-done": {},
  "maint-plan-file": { pick: "D:\\Backups\\ubuntu-24.04.tar" },
  "maint-uninstalling": {
    uninstall: {
      steps: steps(uninstallTitles, [D, R]),
      log: ["Exporting Ubuntu-24.04 to D:\\Backups\\ubuntu-24.04.tar", "Export in progress, this may take a few minutes."],
      running: true,
      done: false,
    },
  },
  "maint-uninstall-failed": {
    uninstall: {
      steps: steps(uninstallTitles, [S, F]),
      log: ["Removing the certificate from the current user's Root store"],
      running: false,
      done: false,
      error: "Remove the trusted certificate: the operation was canceled by the user",
    },
  },
  "maint-uninstalled": {
    uninstall: { steps: steps(uninstallTitles, [S, D, D, D, D]), log: null, running: false, done: true },
  },

  // Setup window.
  "setup-fresh": { view: notInstalledView, report: noDistro, setup: freshSetup },
  "setup-running": {
    view: notInstalledView,
    report: noDistro,
    setup: { steps: steps(setupTitles, [D, D, R]), log: setupLog, running: true, done: false },
  },
  "setup-failed": {
    view: { ...stoppedView, distro: State.StateRunning, level: "partial" },
    report: noDistro,
    setup: {
      steps: steps(setupTitles, [D, D, F]),
      log: [...setupLog.slice(0, 5), "E: Unable to fetch some archives, maybe run apt-get update or try with --fix-missing?"],
      running: false,
      done: false,
      error: "Install Docker Engine: exit status 100",
    },
  },
  "setup-done": {
    report: noDistro,
    setup: {
      steps: steps(setupTitles, [D, D, D, D, D, D, D, D]),
      log: [...setupLog, "Portainer 2.45.3 answers on https://localhost:9443", "docker command added to the user PATH"],
      running: false,
      done: true,
    },
  },
  "setup-existing": {
    report: { ...installed, findings: legacyFindings },
    reportAfterAdopt: installed,
    setup: freshSetup,
  },
  "setup-installed": { report: installed, setup: freshSetup },
  "setup-adopt-error": {
    report: { ...installed, findings: legacyFindings },
    setup: freshSetup,
    fail: {
      "AdoptService.Adopt":
        "Replace the manual docker command with Docker Zero's: remove C:\\Users\\you\\docker-wsl\\bin\\docker.cmd: The process cannot access the file because it is being used by another process.",
    },
  },
};

export const scenario: Scenario = scenarios[stateName || "running"] ?? {};
export const call = caller(scenario);

export type { Mocked };
