window.INVENTORY = {
 "generated": "2026-10-22",
 "commit": "1559594",
 "scale": 2,
 "windows": {
  "popup": {
   "title": "Tray popup",
   "size": [
    360,
    420
   ],
   "component": "frontend/src/Popup.svelte",
   "opens": "Click the tray icon, or launch the exe while it already runs."
  },
  "corner": {
   "title": "Corner box",
   "size": [
    280,
    96
   ],
   "component": "frontend/src/Corner.svelte",
   "opens": "Corner box switch in the popup, or Show corner box in the tray menu; stays on top at the bottom right."
  },
  "settings": {
   "title": "Settings",
   "size": [
    520,
    680
   ],
   "component": "frontend/src/Settings.svelte",
   "opens": "Gear in the popup, or Settings... in the tray menu."
  },
  "maintenance": {
   "title": "Settings: Portainer and Danger zone",
   "size": [
    520,
    680
   ],
   "component": "frontend/src/Maintenance.svelte",
   "opens": "Bottom of the Settings window, below the certificate.",
   "view": "settings"
  },
  "setup": {
   "title": "Setup",
   "size": [
    520,
    700
   ],
   "component": "frontend/src/Setup.svelte",
   "opens": "Opens by itself at launch while Ubuntu-24.04 is missing; Set up tile or Set up... link in the popup; Set up Docker Zero... in the tray menu."
  }
 },
 "shots": [
  {
   "id": "popup-loading",
   "window": "popup",
   "state": "loading",
   "title": "Loading",
   "how": "The first moment after the popup opens, before a status arrives.",
   "note": "Start is enabled and the pill reads Stopped before any status is known.",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 14
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 13
   }
  },
  {
   "id": "popup-running",
   "window": "popup",
   "state": "running",
   "title": "Running",
   "how": "Distro, Docker and Portainer up; keepalive hold on.",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 20
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 20
   }
  },
  {
   "id": "popup-starting",
   "window": "popup",
   "state": "starting",
   "title": "Starting",
   "how": "Right after Start, until Portainer answers (up to 90 s).",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 17
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 17
   }
  },
  {
   "id": "popup-partial",
   "window": "popup",
   "state": "partial",
   "title": "Distro up, Portainer down",
   "how": "The container exited while the distro runs.",
   "note": "The pill reads Starting although nothing is starting.",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 19
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 19
   }
  },
  {
   "id": "popup-stopping",
   "window": "popup",
   "state": "stopping",
   "title": "Stopping",
   "how": "Right after Stop distro.",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 18
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 18
   }
  },
  {
   "id": "popup-stopped",
   "window": "popup",
   "state": "stopped",
   "title": "Stopped",
   "how": "Distro stopped by the user or by the WSL idle timeout.",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 18
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 18
   }
  },
  {
   "id": "popup-not-installed",
   "window": "popup",
   "state": "not-installed",
   "title": "Not installed",
   "how": "Before setup, or after Uninstall all.",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 19
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 19
   }
  },
  {
   "id": "popup-no-wsl",
   "window": "popup",
   "state": "no-wsl",
   "title": "WSL missing",
   "how": "WSL is not enabled on the machine.",
   "note": "The raw Go error is the only guidance; nothing says how to get WSL.",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 25
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 24
   }
  },
  {
   "id": "popup-confirm-shutdown",
   "window": "popup",
   "state": "running",
   "title": "Confirm Shut down WSL",
   "how": "Shut down WSL clicked.",
   "note": "",
   "actions": [
    "Shut down WSL"
   ],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 24
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 23
   }
  },
  {
   "id": "popup-start-error",
   "window": "popup",
   "state": "start-error",
   "title": "Start failed",
   "how": "Start clicked and Portainer did not answer in time.",
   "note": "",
   "actions": [
    "Start"
   ],
   "crop": "window",
   "light": {
    "w": 360,
    "h": 420,
    "kb": 20
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 20
   }
  },
  {
   "id": "corner-running",
   "window": "corner",
   "state": "running",
   "title": "Running",
   "how": "Portainer up.",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 280,
    "h": 96,
    "kb": 5
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 5
   }
  },
  {
   "id": "corner-starting",
   "window": "corner",
   "state": "starting",
   "title": "Working",
   "how": "Any action in progress.",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 280,
    "h": 96,
    "kb": 4
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 4
   }
  },
  {
   "id": "corner-stopped",
   "window": "corner",
   "state": "stopped",
   "title": "Stopped",
   "how": "Distro stopped.",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 280,
    "h": 96,
    "kb": 4
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 4
   }
  },
  {
   "id": "corner-start-error",
   "window": "corner",
   "state": "start-error",
   "title": "Start failed",
   "how": "Start clicked and Portainer did not answer in time.",
   "note": "The error line is squeezed to zero height in the 96 px box, so a failed Start shows nothing.",
   "actions": [
    "Start"
   ],
   "crop": "window",
   "light": {
    "w": 280,
    "h": 96,
    "kb": 4
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 4
   }
  },
  {
   "id": "corner-not-installed",
   "window": "corner",
   "state": "not-installed-start",
   "title": "Not installed",
   "how": "Corner box on while the distro is missing; Start clicked.",
   "note": "Offers Start, which can only fail, and the failure is not visible; no way to set up from here.",
   "actions": [
    "Start"
   ],
   "crop": "window",
   "light": {
    "w": 280,
    "h": 96,
    "kb": 4
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 4
   }
  },
  {
   "id": "settings-loading",
   "window": "settings",
   "state": "settings-loading",
   "title": "Loading",
   "how": "The first moment after the window opens.",
   "note": "",
   "actions": [],
   "crop": "window",
   "light": {
    "w": 520,
    "h": 680,
    "kb": 5
   },
   "dark": {
    "w": 520,
    "h": 680,
    "kb": 5
   }
  },
  {
   "id": "settings",
   "window": "settings",
   "state": "settings",
   "title": "All settings",
   "how": "Typical state after setup: auto start, shortcuts, keepalive, trusted certificate.",
   "note": "",
   "actions": [],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 1271,
    "kb": 119
   },
   "dark": {
    "w": 520,
    "h": 1271,
    "kb": 117
   }
  },
  {
   "id": "settings-off",
   "window": "settings",
   "state": "settings-off",
   "title": "Everything off",
   "how": "Fresh portable run: no auto start, no shortcuts, WSL default idle.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 553,
    "kb": 53
   },
   "dark": {
    "w": 520,
    "h": 553,
    "kb": 52
   }
  },
  {
   "id": "settings-wslconfig",
   "window": "settings",
   "state": "settings-wslconfig",
   "title": ".wslconfig mode pending",
   "how": "Never stop idle distros chosen; WSL has not restarted yet.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 339,
    "kb": 39
   },
   "dark": {
    "w": 520,
    "h": 339,
    "kb": 38
   }
  },
  {
   "id": "settings-wslconfig-confirm",
   "window": "settings",
   "state": "settings-wslconfig",
   "title": "Confirm Apply now",
   "how": "Apply now... clicked.",
   "note": "",
   "actions": [
    "Apply now..."
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 359,
    "kb": 44
   },
   "dark": {
    "w": 520,
    "h": 359,
    "kb": 42
   }
  },
  {
   "id": "settings-browser",
   "window": "settings",
   "state": "settings-browser",
   "title": "Custom browser, no Chrome",
   "how": "Other browser chosen with Browse...; Chrome not installed.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 160,
    "kb": 14
   },
   "dark": {
    "w": 520,
    "h": 160,
    "kb": 14
   }
  },
  {
   "id": "settings-untrusted",
   "window": "settings",
   "state": "settings-untrusted",
   "title": "Certificate not trusted",
   "how": "After setup, before Install (trust).",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 228,
    "kb": 28
   },
   "dark": {
    "w": 520,
    "h": 228,
    "kb": 28
   }
  },
  {
   "id": "settings-no-cert",
   "window": "settings",
   "state": "settings-no-cert",
   "title": "No certificate",
   "how": "Distro running but setup did not reach the certificate step.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 199,
    "kb": 20
   },
   "dark": {
    "w": 520,
    "h": 199,
    "kb": 20
   }
  },
  {
   "id": "settings-cert-unread",
   "window": "settings",
   "state": "settings-cert-unread",
   "title": "Distro stopped",
   "how": "Settings opened while Ubuntu-24.04 is stopped; the certificate is not read.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 199,
    "kb": 20
   },
   "dark": {
    "w": 520,
    "h": 199,
    "kb": 20
   }
  },
  {
   "id": "settings-busy",
   "window": "settings",
   "state": "settings-busy",
   "title": "Working",
   "how": "A setting is being applied.",
   "note": "Working... appears at the very bottom, out of view of the 680 px window.",
   "actions": [
    "Desktop"
   ],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 1300,
    "kb": 117
   },
   "dark": {
    "w": 520,
    "h": 1300,
    "kb": 115
   }
  },
  {
   "id": "settings-problems",
   "window": "settings",
   "state": "settings-problems",
   "title": "Errors",
   "how": "A setting failed to apply, and two settings could not be read.",
   "note": "Errors land below the Danger zone, out of view of the item that failed.",
   "actions": [
    "Desktop"
   ],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 1381,
    "kb": 136
   },
   "dark": {
    "w": 520,
    "h": 1381,
    "kb": 132
   }
  },
  {
   "id": "maint",
   "window": "maintenance",
   "state": "maint",
   "title": "Idle",
   "how": "Settings window scrolled to the bottom.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 296,
    "kb": 26
   },
   "dark": {
    "w": 520,
    "h": 296,
    "kb": 26
   }
  },
  {
   "id": "maint-updating",
   "window": "maintenance",
   "state": "maint-updating",
   "title": "Updating",
   "how": "Update Portainer clicked; pull in progress.",
   "note": "",
   "actions": [
    "Update Portainer"
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 325,
    "kb": 25
   },
   "dark": {
    "w": 520,
    "h": 325,
    "kb": 24
   }
  },
  {
   "id": "maint-updated",
   "window": "maintenance",
   "state": "maint-updated",
   "title": "Updated",
   "how": "A newer lts image was pulled and Portainer restarted.",
   "note": "",
   "actions": [
    "Update Portainer"
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 325,
    "kb": 29
   },
   "dark": {
    "w": 520,
    "h": 325,
    "kb": 29
   }
  },
  {
   "id": "maint-current",
   "window": "maintenance",
   "state": "maint-current",
   "title": "Already latest",
   "how": "Update clicked with the image current.",
   "note": "",
   "actions": [
    "Update Portainer"
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 325,
    "kb": 30
   },
   "dark": {
    "w": 520,
    "h": 325,
    "kb": 29
   }
  },
  {
   "id": "maint-update-failed",
   "window": "maintenance",
   "state": "maint-update-failed",
   "title": "Update failed",
   "how": "Docker Hub unreachable.",
   "note": "",
   "actions": [
    "Update Portainer"
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 342,
    "kb": 40
   },
   "dark": {
    "w": 520,
    "h": 342,
    "kb": 38
   }
  },
  {
   "id": "maint-reset-open",
   "window": "maintenance",
   "state": "maint",
   "title": "Reset form",
   "how": "Reset admin password... clicked.",
   "note": "",
   "actions": [
    "Reset admin password..."
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 577,
    "kb": 46
   },
   "dark": {
    "w": 520,
    "h": 577,
    "kb": 46
   }
  },
  {
   "id": "maint-reset-invalid",
   "window": "maintenance",
   "state": "maint",
   "title": "Reset form: short password",
   "how": "Reset password clicked with the fields empty.",
   "note": "",
   "actions": [
    "Reset admin password...",
    "Reset password"
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 618,
    "kb": 50
   },
   "dark": {
    "w": 520,
    "h": 618,
    "kb": 49
   }
  },
  {
   "id": "maint-reset-wrong-user",
   "window": "maintenance",
   "state": "maint-reset-wrong-user",
   "title": "Reset: wrong username",
   "how": "Password replaced, but the username check failed.",
   "note": "",
   "actions": [
    "Reset admin password...",
    "type Admin username",
    "type New password",
    "type Repeat new password",
    "Reset password"
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 635,
    "kb": 58
   },
   "dark": {
    "w": 520,
    "h": 635,
    "kb": 57
   }
  },
  {
   "id": "maint-reset-done",
   "window": "maintenance",
   "state": "maint-reset-done",
   "title": "Reset done",
   "how": "Password replaced and checked.",
   "note": "",
   "actions": [
    "Reset admin password...",
    "type New password",
    "type Repeat new password",
    "Reset password"
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 325,
    "kb": 32
   },
   "dark": {
    "w": 520,
    "h": 325,
    "kb": 31
   }
  },
  {
   "id": "maint-plan",
   "window": "maintenance",
   "state": "maint",
   "title": "Uninstall plan",
   "how": "Uninstall all... clicked; backup on, no file yet.",
   "note": "",
   "actions": [
    "Uninstall all..."
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 528,
    "kb": 60
   },
   "dark": {
    "w": 520,
    "h": 528,
    "kb": 61
   }
  },
  {
   "id": "maint-plan-file",
   "window": "maintenance",
   "state": "maint-plan-file",
   "title": "Uninstall plan with backup file",
   "how": "Choose file... used.",
   "note": "",
   "actions": [
    "Uninstall all...",
    "Choose file..."
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 528,
    "kb": 62
   },
   "dark": {
    "w": 520,
    "h": 528,
    "kb": 63
   }
  },
  {
   "id": "maint-confirm",
   "window": "maintenance",
   "state": "maint",
   "title": "Uninstall confirmation",
   "how": "Backup unticked, Continue clicked.",
   "note": "",
   "actions": [
    "Uninstall all...",
    "Save a backup of Ubuntu-24.04 first",
    "Continue"
   ],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 347,
    "kb": 31
   },
   "dark": {
    "w": 520,
    "h": 347,
    "kb": 30
   }
  },
  {
   "id": "maint-uninstalling",
   "window": "maintenance",
   "state": "maint-uninstalling",
   "title": "Uninstalling",
   "how": "Uninstall now clicked; backup export running.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 457,
    "kb": 43
   },
   "dark": {
    "w": 520,
    "h": 457,
    "kb": 43
   }
  },
  {
   "id": "maint-uninstall-failed",
   "window": "maintenance",
   "state": "maint-uninstall-failed",
   "title": "Uninstall failed",
   "how": "The Windows certificate dialog was answered No.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 502,
    "kb": 53
   },
   "dark": {
    "w": 520,
    "h": 502,
    "kb": 52
   }
  },
  {
   "id": "maint-uninstalled",
   "window": "maintenance",
   "state": "maint-uninstalled",
   "title": "Uninstalled",
   "how": "Every step finished; only Quit is left.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "light": {
    "w": 520,
    "h": 449,
    "kb": 51
   },
   "dark": {
    "w": 520,
    "h": 449,
    "kb": 50
   }
  },
  {
   "id": "setup-fresh",
   "window": "setup",
   "state": "setup-fresh",
   "title": "First run",
   "how": "First launch on a machine without Ubuntu-24.04.",
   "note": "",
   "actions": [],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 739,
    "kb": 47
   },
   "dark": {
    "w": 520,
    "h": 739,
    "kb": 48
   }
  },
  {
   "id": "setup-invalid",
   "window": "setup",
   "state": "setup-fresh",
   "title": "Form error",
   "how": "Start setup clicked with the password empty.",
   "note": "",
   "actions": [
    "Start setup"
   ],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 771,
    "kb": 51
   },
   "dark": {
    "w": 520,
    "h": 771,
    "kb": 52
   }
  },
  {
   "id": "setup-running",
   "window": "setup",
   "state": "setup-running",
   "title": "Running",
   "how": "Setup in progress, third step.",
   "note": "",
   "actions": [],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 700,
    "kb": 59
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 59
   }
  },
  {
   "id": "setup-failed",
   "window": "setup",
   "state": "setup-failed",
   "title": "Failed",
   "how": "A step failed; the form returns with Try again.",
   "note": "",
   "actions": [],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 909,
    "kb": 74
   },
   "dark": {
    "w": 520,
    "h": 909,
    "kb": 74
   }
  },
  {
   "id": "setup-done",
   "window": "setup",
   "state": "setup-done",
   "title": "All set",
   "how": "Every step finished.",
   "note": "The steps, the log and the form fields stay above the result; the All set card is last.",
   "actions": [],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 752,
    "kb": 80
   },
   "dark": {
    "w": 520,
    "h": 752,
    "kb": 80
   }
  },
  {
   "id": "setup-existing",
   "window": "setup",
   "state": "setup-existing",
   "title": "Existing setup found",
   "how": "A setup made by hand from the old guide is detected.",
   "note": "The take-over card and the full setup form show together, though the card says setup is not needed.",
   "actions": [],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 1078,
    "kb": 98
   },
   "dark": {
    "w": 520,
    "h": 1078,
    "kb": 97
   }
  },
  {
   "id": "setup-adopted",
   "window": "setup",
   "state": "setup-existing",
   "title": "Taken over",
   "how": "Take over clicked; every item replaced.",
   "note": "",
   "actions": [
    "Take over"
   ],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 929,
    "kb": 71
   },
   "dark": {
    "w": 520,
    "h": 929,
    "kb": 71
   }
  },
  {
   "id": "setup-adopt-error",
   "window": "setup",
   "state": "setup-adopt-error",
   "title": "Take over failed",
   "how": "A legacy file was locked by another process.",
   "note": "",
   "actions": [
    "Take over"
   ],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 1159,
    "kb": 116
   },
   "dark": {
    "w": 520,
    "h": 1159,
    "kb": 114
   }
  },
  {
   "id": "setup-installed",
   "window": "setup",
   "state": "setup-installed",
   "title": "Opened after setup",
   "how": "Set up... opened again on a machine Docker Zero already set up.",
   "note": "Reads Existing setup found / Nothing left to replace, above a form that would run setup again.",
   "actions": [],
   "crop": "full",
   "light": {
    "w": 520,
    "h": 900,
    "kb": 66
   },
   "dark": {
    "w": 520,
    "h": 900,
    "kb": 67
   }
  }
 ],
 "real": [
  {
   "id": "real-popup",
   "window": "popup",
   "pair": "popup-running",
   "title": "Popup, built exe",
   "w": 360,
   "h": 420
  },
  {
   "id": "real-corner",
   "window": "corner",
   "pair": "corner-running",
   "title": "Corner box, built exe",
   "w": 280,
   "h": 96
  },
  {
   "id": "real-settings",
   "window": "settings",
   "pair": "settings",
   "title": "Settings, built exe, first screen",
   "w": 520,
   "h": 680
  },
  {
   "id": "real-setup",
   "window": "setup",
   "pair": "setup-installed",
   "title": "Setup, built exe, opened after setup",
   "w": 520,
   "h": 700
  }
 ]
};
