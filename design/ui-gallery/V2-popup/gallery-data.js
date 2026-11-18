window.GALLERY = {
 "title": "Docker Zero UI redesign board, V2 popup (popup and corner box)",
 "lede": "Proposed redesign (direction B: centred ring with the logo's orca, brand colours, long-step feedback), rendered from the app frontend with the same mocked services as the current-UI board. The Settings window moved to V2-setting and the Setup window to V2-setup. Every state keeps the id of its current-UI twin; the motion frames and the long-step setup states are extra. Mark each Keep, Revise or Discuss.",
 "regenerate": [
  "cd frontend && npm run gallery",
  "uvx --from playwright --with pillow python -I <frontdoor-ui-gallery>/scripts/shoot.py render design/ui-gallery/V2-popup/gallery.json"
 ],
 "generated": "2026-11-19",
 "commit": "3579917",
 "storageKey": "docker-zero-ui-gallery2",
 "views": {
  "popup": {
   "scroller": "main",
   "title": "Tray popup",
   "size": [
    360,
    420
   ],
   "source": "frontend/src/Popup.svelte",
   "opens": "Click the tray icon, or launch the exe while it already runs.",
   "url": "/?view=popup"
  },
  "corner": {
   "scroller": "main",
   "title": "Corner box",
   "size": [
    280,
    96
   ],
   "source": "frontend/src/Corner.svelte",
   "opens": "Corner box switch in the popup, or Show corner box in the tray menu; stays on top at the bottom right.",
   "url": "/?view=corner"
  }
 },
 "states": [
  {
   "id": "popup-loading",
   "view": "popup",
   "state": "loading",
   "title": "Loading",
   "how": "The first moment after the popup opens, before a status arrives.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 16
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 16
   }
  },
  {
   "id": "popup-running",
   "view": "popup",
   "state": "running",
   "title": "Running",
   "how": "Distro, Docker and Portainer up; keepalive hold on.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 22
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 21
   }
  },
  {
   "id": "popup-starting",
   "view": "popup",
   "state": "starting",
   "title": "Starting",
   "how": "Right after Start, until Portainer answers (up to 90 s).",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
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
   "id": "popup-partial",
   "view": "popup",
   "state": "partial",
   "title": "Distro up, Portainer down",
   "how": "The container exited while the distro runs.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 22
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 22
   }
  },
  {
   "id": "popup-stopping",
   "view": "popup",
   "state": "stopping",
   "title": "Stopping",
   "how": "Right after Stop distro.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 19
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 18
   }
  },
  {
   "id": "popup-stopped",
   "view": "popup",
   "state": "stopped",
   "title": "Stopped",
   "how": "Distro stopped by the user or by the WSL idle timeout.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
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
   "view": "popup",
   "state": "not-installed",
   "title": "Not installed",
   "how": "Before setup, or after Uninstall all.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 21
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 21
   }
  },
  {
   "id": "popup-no-wsl",
   "view": "popup",
   "state": "no-wsl",
   "title": "WSL missing",
   "how": "WSL is not enabled on the machine.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 19
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 18
   }
  },
  {
   "id": "popup-no-wsl-info",
   "view": "popup",
   "state": "no-wsl",
   "title": "WSL missing, details open",
   "how": "The i button opened: the raw error text sits behind it.",
   "note": "",
   "actions": [
    "click More information"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
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
   "id": "popup-confirm-shutdown",
   "view": "popup",
   "state": "running",
   "title": "Confirm Shut down WSL",
   "how": "Shut down WSL clicked.",
   "note": "",
   "actions": [
    "click More",
    "click Shut down WSL"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 28
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 27
   }
  },
  {
   "id": "popup-start-error",
   "view": "popup",
   "state": "start-error",
   "title": "Start failed",
   "how": "Start clicked and Portainer did not answer in time.",
   "note": "",
   "actions": [
    "click Start"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 20
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 19
   }
  },
  {
   "id": "popup-start-hover",
   "view": "popup",
   "state": "stopped",
   "title": "Pointer on Start",
   "how": "The pointer rests on the Start button.",
   "note": "",
   "actions": [
    "hover Start"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 19
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 18
   }
  },
  {
   "id": "popup-more",
   "view": "popup",
   "state": "running",
   "title": "More menu open",
   "how": "The three-dot menu: Settings, Corner box, Set up..., Shut down WSL.",
   "note": "",
   "actions": [
    "click More"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 23
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 23
   }
  },
  {
   "id": "popup-starting-1",
   "view": "popup",
   "state": "starting-1",
   "title": "Motion: Start, 1 of 4 parts",
   "how": "Right after Start: only WSL is up; the sea is low.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
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
   "id": "popup-starting-2",
   "view": "popup",
   "state": "starting-2",
   "title": "Motion: Start, 2 of 4 parts",
   "how": "Ubuntu is up, Docker is not yet.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 20
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 19
   }
  },
  {
   "id": "popup-leap-0",
   "view": "popup",
   "state": "running",
   "title": "Motion: leap, start",
   "how": "The moment Portainer answers: the orca low in the sea, tilted.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 22
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 21
   }
  },
  {
   "id": "popup-leap-50",
   "view": "popup",
   "state": "running",
   "title": "Motion: leap, peak",
   "how": "Half way through the single 1.4 s leap: the orca at its highest.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 22
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 21
   }
  },
  {
   "id": "popup-leap-100",
   "view": "popup",
   "state": "running",
   "title": "Motion: leap, end",
   "how": "The leap has finished; the orca rests as in the logo. Nothing loops after this.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 360,
    "h": 420,
    "kb": 22
   },
   "dark": {
    "w": 360,
    "h": 420,
    "kb": 21
   }
  },
  {
   "id": "corner-running",
   "view": "corner",
   "state": "running",
   "title": "Running",
   "how": "Portainer up.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
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
   "view": "corner",
   "state": "starting",
   "title": "Working",
   "how": "Any action in progress.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
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
   "id": "corner-stopped",
   "view": "corner",
   "state": "stopped",
   "title": "Stopped",
   "how": "Distro stopped.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
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
   "id": "corner-start-error",
   "view": "corner",
   "state": "start-error",
   "title": "Start failed",
   "how": "Start clicked and Portainer did not answer in time.",
   "note": "",
   "actions": [
    "click Start"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 280,
    "h": 96,
    "kb": 7
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 7
   }
  },
  {
   "id": "corner-start-error-info",
   "view": "corner",
   "state": "start-error",
   "title": "Start failed, details open",
   "how": "The i button opened; the raw error replaces the buttons until it is closed.",
   "note": "",
   "actions": [
    "click Start",
    "click More information"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 280,
    "h": 96,
    "kb": 8
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 8
   }
  },
  {
   "id": "corner-not-installed",
   "view": "corner",
   "state": "not-installed-start",
   "title": "Not installed",
   "how": "Corner box on while the distro is missing; Start clicked.",
   "note": "",
   "actions": [
    "click Start"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 280,
    "h": 96,
    "kb": 8
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 8
   }
  },
  {
   "id": "corner-not-installed-info",
   "view": "corner",
   "state": "not-installed-start",
   "title": "Not installed, details open",
   "how": "The i button opened.",
   "note": "",
   "actions": [
    "click Start",
    "click More information"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 280,
    "h": 96,
    "kb": 10
   },
   "dark": {
    "w": 280,
    "h": 96,
    "kb": 10
   }
  }
 ],
 "skipped": [],
 "real": []
};
