# Renders every window state of the app from the mocked frontend and writes
# the inventory data the page reads.
#
#   cd frontend && npm run inventory          (serves http://127.0.0.1:9246)
#   uvx --from playwright --with pillow python -I design/ui-inventory/shoot.py
#
# Options: --only <id,id> to re-shoot some states, --base <url> for another
# server. Uses the Edge already on the machine; nothing is downloaded.
import argparse
import io
import json
import subprocess
import sys
from datetime import date
from pathlib import Path

from PIL import Image
from playwright.sync_api import sync_playwright

HERE = Path(__file__).resolve().parent
SCALE = 2

WINDOWS = {
    "popup": {
        "title": "Tray popup",
        "size": [360, 420],
        "component": "frontend/src/Popup.svelte",
        "opens": "Click the tray icon, or launch the exe while it already runs.",
    },
    "corner": {
        "title": "Corner box",
        "size": [280, 96],
        "component": "frontend/src/Corner.svelte",
        "opens": "Corner box switch in the popup, or Show corner box in the tray menu; stays on top at the bottom right.",
    },
    "settings": {
        "title": "Settings",
        "size": [520, 680],
        "component": "frontend/src/Settings.svelte",
        "opens": "Gear in the popup, or Settings... in the tray menu.",
    },
    "maintenance": {
        "title": "Settings: Portainer and Danger zone",
        "size": [520, 680],
        "component": "frontend/src/Maintenance.svelte",
        "opens": "Bottom of the Settings window, below the certificate.",
        "view": "settings",
    },
    "setup": {
        "title": "Setup",
        "size": [520, 700],
        "component": "frontend/src/Setup.svelte",
        "opens": "Opens by itself at launch while Ubuntu-24.04 is missing; Set up tile or Set up... link in the popup; Set up Docker Zero... in the tray menu.",
    },
}

CERT = ["Certificate for https://localhost:9443"]
MAINT = ["Portainer", "Danger zone"]


def click(name, role="button"):
    return ("click", role, name)


def fill(label, value):
    return ("fill", label, value)


PASSWORDS = [fill("New password", "inventory-example-pass"), fill("Repeat new password", "inventory-example-pass")]

# id, window, scenario, title, how the user gets there, options
SHOTS = [
    ("popup-loading", "popup", "loading", "Loading", "The first moment after the popup opens, before a status arrives.",
     {"note": "Start is enabled and the pill reads Stopped before any status is known."}),
    ("popup-running", "popup", "running", "Running", "Distro, Docker and Portainer up; keepalive hold on.", {}),
    ("popup-starting", "popup", "starting", "Starting", "Right after Start, until Portainer answers (up to 90 s).", {}),
    ("popup-partial", "popup", "partial", "Distro up, Portainer down", "The container exited while the distro runs.",
     {"note": "The pill reads Starting although nothing is starting."}),
    ("popup-stopping", "popup", "stopping", "Stopping", "Right after Stop distro.", {}),
    ("popup-stopped", "popup", "stopped", "Stopped", "Distro stopped by the user or by the WSL idle timeout.", {}),
    ("popup-not-installed", "popup", "not-installed", "Not installed", "Before setup, or after Uninstall all.", {}),
    ("popup-no-wsl", "popup", "no-wsl", "WSL missing", "WSL is not enabled on the machine.",
     {"note": "The raw Go error is the only guidance; nothing says how to get WSL."}),
    ("popup-confirm-shutdown", "popup", "running", "Confirm Shut down WSL", "Shut down WSL clicked.", {"actions": [click("Shut down WSL")]}),
    ("popup-start-error", "popup", "start-error", "Start failed", "Start clicked and Portainer did not answer in time.", {"actions": [click("Start")]}),
    ("corner-running", "corner", "running", "Running", "Portainer up.", {}),
    ("corner-starting", "corner", "starting", "Working", "Any action in progress.", {}),
    ("corner-stopped", "corner", "stopped", "Stopped", "Distro stopped.", {}),
    ("corner-start-error", "corner", "start-error", "Start failed", "Start clicked and Portainer did not answer in time.",
     {"actions": [click("Start")], "note": "The error line is squeezed to zero height in the 96 px box, so a failed Start shows nothing."}),
    ("corner-not-installed", "corner", "not-installed-start", "Not installed", "Corner box on while the distro is missing; Start clicked.",
     {"actions": [click("Start")], "note": "Offers Start, which can only fail, and the failure is not visible; no way to set up from here."}),
    ("settings-loading", "settings", "settings-loading", "Loading", "The first moment after the window opens.", {"window": True}),
    ("settings", "settings", "settings", "All settings", "Typical state after setup: auto start, shortcuts, keepalive, trusted certificate.", {"full": True}),
    ("settings-off", "settings", "settings-off", "Everything off", "Fresh portable run: no auto start, no shortcuts, WSL default idle.",
     {"cards": ["Start with Windows", "Shortcuts", "Keep Ubuntu-24.04 running"]}),
    ("settings-wslconfig", "settings", "settings-wslconfig", ".wslconfig mode pending", "Never stop idle distros chosen; WSL has not restarted yet.",
     {"cards": ["Keep Ubuntu-24.04 running"]}),
    ("settings-wslconfig-confirm", "settings", "settings-wslconfig", "Confirm Apply now", "Apply now... clicked.",
     {"cards": ["Keep Ubuntu-24.04 running"], "actions": [click("Apply now...")]}),
    ("settings-browser", "settings", "settings-browser", "Custom browser, no Chrome", "Other browser chosen with Browse...; Chrome not installed.",
     {"cards": ["Browser for Open Portainer"]}),
    ("settings-untrusted", "settings", "settings-untrusted", "Certificate not trusted", "After setup, before Install (trust).", {"cards": CERT}),
    ("settings-no-cert", "settings", "settings-no-cert", "No certificate", "Distro running but setup did not reach the certificate step.", {"cards": CERT}),
    ("settings-cert-unread", "settings", "settings-cert-unread", "Distro stopped", "Settings opened while Ubuntu-24.04 is stopped; the certificate is not read.", {"cards": CERT}),
    ("settings-busy", "settings", "settings-busy", "Working", "A setting is being applied.",
     {"full": True, "actions": [click("Desktop", "checkbox")], "note": "Working... appears at the very bottom, out of view of the 680 px window."}),
    ("settings-problems", "settings", "settings-problems", "Errors", "A setting failed to apply, and two settings could not be read.",
     {"full": True, "actions": [click("Desktop", "checkbox")], "note": "Errors land below the Danger zone, out of view of the item that failed."}),
    ("maint", "maintenance", "maint", "Idle", "Settings window scrolled to the bottom.", {"cards": MAINT}),
    ("maint-updating", "maintenance", "maint-updating", "Updating", "Update Portainer clicked; pull in progress.", {"cards": MAINT, "actions": [click("Update Portainer")]}),
    ("maint-updated", "maintenance", "maint-updated", "Updated", "A newer lts image was pulled and Portainer restarted.", {"cards": MAINT, "actions": [click("Update Portainer")]}),
    ("maint-current", "maintenance", "maint-current", "Already latest", "Update clicked with the image current.", {"cards": MAINT, "actions": [click("Update Portainer")]}),
    ("maint-update-failed", "maintenance", "maint-update-failed", "Update failed", "Docker Hub unreachable.", {"cards": MAINT, "actions": [click("Update Portainer")]}),
    ("maint-reset-open", "maintenance", "maint", "Reset form", "Reset admin password... clicked.", {"cards": MAINT, "actions": [click("Reset admin password...")]}),
    ("maint-reset-invalid", "maintenance", "maint", "Reset form: short password", "Reset password clicked with the fields empty.",
     {"cards": MAINT, "actions": [click("Reset admin password..."), click("Reset password")]}),
    ("maint-reset-wrong-user", "maintenance", "maint-reset-wrong-user", "Reset: wrong username", "Password replaced, but the username check failed.",
     {"cards": MAINT, "actions": [click("Reset admin password..."), fill("Admin username", "nobody"), *PASSWORDS, click("Reset password")]}),
    ("maint-reset-done", "maintenance", "maint-reset-done", "Reset done", "Password replaced and checked.",
     {"cards": MAINT, "actions": [click("Reset admin password..."), *PASSWORDS, click("Reset password")]}),
    ("maint-plan", "maintenance", "maint", "Uninstall plan", "Uninstall all... clicked; backup on, no file yet.", {"cards": MAINT, "actions": [click("Uninstall all...")]}),
    ("maint-plan-file", "maintenance", "maint-plan-file", "Uninstall plan with backup file", "Choose file... used.",
     {"cards": MAINT, "actions": [click("Uninstall all..."), click("Choose file...")]}),
    ("maint-confirm", "maintenance", "maint", "Uninstall confirmation", "Backup unticked, Continue clicked.",
     {"cards": MAINT, "actions": [click("Uninstall all..."), click("Save a backup of Ubuntu-24.04 first", "checkbox"), click("Continue")]}),
    ("maint-uninstalling", "maintenance", "maint-uninstalling", "Uninstalling", "Uninstall now clicked; backup export running.", {"cards": MAINT}),
    ("maint-uninstall-failed", "maintenance", "maint-uninstall-failed", "Uninstall failed", "The Windows certificate dialog was answered No.", {"cards": MAINT}),
    ("maint-uninstalled", "maintenance", "maint-uninstalled", "Uninstalled", "Every step finished; only Quit is left.", {"cards": MAINT}),
    ("setup-fresh", "setup", "setup-fresh", "First run", "First launch on a machine without Ubuntu-24.04.", {"full": True}),
    ("setup-invalid", "setup", "setup-fresh", "Form error", "Start setup clicked with the password empty.", {"full": True, "actions": [click("Start setup")]}),
    ("setup-running", "setup", "setup-running", "Running", "Setup in progress, third step.", {"full": True}),
    ("setup-failed", "setup", "setup-failed", "Failed", "A step failed; the form returns with Try again.", {"full": True}),
    ("setup-done", "setup", "setup-done", "All set", "Every step finished.", {"full": True,
     "note": "The steps, the log and the form fields stay above the result; the All set card is last."}),
    ("setup-existing", "setup", "setup-existing", "Existing setup found", "A setup made by hand from the old guide is detected.", {"full": True,
     "note": "The take-over card and the full setup form show together, though the card says setup is not needed."}),
    ("setup-adopted", "setup", "setup-existing", "Taken over", "Take over clicked; every item replaced.", {"full": True, "actions": [click("Take over")]}),
    ("setup-adopt-error", "setup", "setup-adopt-error", "Take over failed", "A legacy file was locked by another process.", {"full": True, "actions": [click("Take over")]}),
    ("setup-installed", "setup", "setup-installed", "Opened after setup", "Set up... opened again on a machine Docker Zero already set up.", {"full": True,
     "note": "Reads Existing setup found / Nothing left to replace, above a form that would run setup again."}),
]

# Screenshots of the built exe on a real account, made with real.ps1.
REAL = [
    ("real-popup", "popup", "popup-running", "Popup, built exe"),
    ("real-corner", "corner", "corner-running", "Corner box, built exe"),
    ("real-settings", "settings", "settings", "Settings, built exe, first screen"),
    ("real-setup", "setup", "setup-installed", "Setup, built exe, opened after setup"),
]

MEASURE = """(cards) => {
  const main = document.querySelector('main');
  const pick = cards.length ? [...document.querySelectorAll('section.card')].filter(s => {
    const h = s.querySelector('h2');
    return h && cards.some(c => h.textContent.trim().startsWith(c));
  }) : [];
  const r = pick.map(s => s.getBoundingClientRect());
  return {
    height: main ? main.scrollHeight : document.body.scrollHeight,
    found: pick.length,
    clip: r.length ? { top: Math.min(...r.map(x => x.top)), bottom: Math.max(...r.map(x => x.bottom)) } : null,
  };
}"""


def act(page, step):
    kind = step[0]
    if kind == "click":
        _, role, name = step
        page.get_by_role(role, name=name, exact=True).click()
    elif kind == "fill":
        _, label, value = step
        page.get_by_label(label, exact=True).fill(value)
    page.wait_for_timeout(450)


def shoot(page, base, shot, scheme):
    sid, win, state, _title, _how, opt = shot
    meta = WINDOWS[win]
    w, h = meta["size"]
    page.set_viewport_size({"width": w, "height": h})
    page.emulate_media(color_scheme=scheme)
    page.goto(f"{base}/?view={meta.get('view', win)}&state={state}")
    page.wait_for_timeout(700)
    for step in opt.get("actions", []):
        act(page, step)
    tall = opt.get("full") or opt.get("cards")
    if tall:
        m = page.evaluate(MEASURE, opt.get("cards", []))
        page.set_viewport_size({"width": w, "height": max(h, m["height"])})
        page.wait_for_timeout(200)
    m = page.evaluate(MEASURE, opt.get("cards", []))
    clip = None
    if opt.get("cards"):
        if m["found"] != len(opt["cards"]) or not m["clip"]:
            raise SystemExit(f"{sid}: found {m['found']} of the cards {opt['cards']}")
        top = max(0, m["clip"]["top"] - 12)
        clip = {"x": 0, "y": top, "width": w, "height": m["clip"]["bottom"] + 12 - top}
    png = page.screenshot(clip=clip, full_page=False)
    img = Image.open(io.BytesIO(png)).convert("RGB")
    out = HERE / "shots" / f"{sid}-{scheme}.webp"
    img.save(out, "WEBP", quality=80, method=6)
    return {"w": img.width // SCALE, "h": img.height // SCALE, "kb": round(out.stat().st_size / 1024)}


def main():
    sys.stdout.reconfigure(encoding="utf-8")
    ap = argparse.ArgumentParser()
    ap.add_argument("--base", default="http://127.0.0.1:9246")
    ap.add_argument("--only", default="")
    args = ap.parse_args()
    only = {s for s in args.only.split(",") if s}
    (HERE / "shots").mkdir(exist_ok=True)
    data_file = HERE / "inventory.js"
    old = {}
    if data_file.exists():
        text = data_file.read_text(encoding="utf-8")
        old = {s["id"]: s for s in json.loads(text[text.index("{"): text.rindex("}") + 1])["shots"]}

    shots = []
    with sync_playwright() as p:
        browser = p.chromium.launch(channel="msedge")
        ctx = browser.new_context(device_scale_factor=SCALE)
        page = ctx.new_page()
        errors = []
        page.on("pageerror", lambda e: errors.append(str(e)))
        for shot in SHOTS:
            sid, win, state, title, how, opt = shot
            entry = {"id": sid, "window": win, "state": state, "title": title, "how": how,
                     "note": opt.get("note", ""), "actions": [s[-1] if s[0] == "click" else f"type {s[1]}" for s in opt.get("actions", [])],
                     "crop": "cards" if opt.get("cards") else ("full" if opt.get("full") else "window")}
            if only and sid not in only and sid in old:
                shots.append(old[sid])
                continue
            for scheme in ("light", "dark"):
                entry[scheme] = shoot(page, args.base, shot, scheme)
            if errors:
                raise SystemExit(f"{sid}: page error {errors}")
            print(f"{sid:28} {entry['light']['w']}x{entry['light']['h']}  {entry['light']['kb']} + {entry['dark']['kb']} KB")
            shots.append(entry)
        browser.close()

    real = []
    for rid, win, pair, title in REAL:
        f = HERE / "real" / f"{rid}.png"
        if f.exists():
            with Image.open(f) as im:
                real.append({"id": rid, "window": win, "pair": pair, "title": title, "w": im.width, "h": im.height})
    commit = subprocess.run(["git", "rev-parse", "--short", "HEAD"], capture_output=True, text=True, cwd=HERE).stdout.strip()
    data = {"generated": date.today().isoformat(), "commit": commit, "scale": SCALE, "windows": WINDOWS, "shots": shots, "real": real}
    data_file.write_bytes(("window.INVENTORY = " + json.dumps(data, indent=1, ensure_ascii=False) + ";\n").encode("utf-8"))
    print(f"{len(shots)} states, {len(real)} real screenshots -> {data_file.name}")


if __name__ == "__main__":
    main()
