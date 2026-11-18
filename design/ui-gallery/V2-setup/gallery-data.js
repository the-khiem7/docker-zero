window.GALLERY = {
 "title": "Docker Zero Setup redesign board, V2 setup (bento)",
 "lede": "The Setup window in the same bento language as V2-setting: a hero card with the ring and chips, a track of the eight steps, a key disc as the start button, take-over tiles, details behind the i button. Same mocked services and state ids as V2-popup. Mark each Keep, Revise or Discuss.",
 "regenerate": [
  "cd frontend && npm run gallery",
  "uvx --from playwright --with pillow python -I <frontdoor-ui-gallery>/scripts/shoot.py render design/ui-gallery/V2-setup/gallery.json"
 ],
 "generated": "2026-11-19",
 "commit": "3579917",
 "storageKey": "docker-zero-ui-gallery-v2-setup",
 "views": {
  "setup": {
   "scroller": "main",
   "title": "Setup",
   "size": [
    520,
    700
   ],
   "source": "frontend/src/Setup.svelte",
   "opens": "Opens by itself at launch while Ubuntu-24.04 is missing; Set up tile or Set up... link in the popup; Set up Docker Zero... in the tray menu.",
   "url": "/?view=setup",
   "timers": "paused"
  }
 },
 "states": [
  {
   "id": "setup-fresh",
   "view": "setup",
   "state": "setup-fresh",
   "title": "First run",
   "how": "First launch on a machine without Ubuntu-24.04.",
   "note": "",
   "actions": [],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 700,
    "kb": 35
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 31
   }
  },
  {
   "id": "setup-invalid",
   "view": "setup",
   "state": "setup-fresh",
   "title": "Form error",
   "how": "Start setup clicked with the password empty.",
   "note": "",
   "actions": [
    "click Start setup"
   ],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 700,
    "kb": 40
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 36
   }
  },
  {
   "id": "setup-distro",
   "view": "setup",
   "state": "setup-distro",
   "title": "Installing the distro",
   "how": "Step 2 of 8, just started.",
   "note": "",
   "actions": [],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 700,
    "kb": 24
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 22
   }
  },
  {
   "id": "setup-distro-slow",
   "view": "setup",
   "state": "setup-distro",
   "title": "Same step, 2 min 30 s later",
   "how": "The clock runs from when the window first saw the step; past 2 minutes a still-working notice appears.",
   "note": "",
   "actions": [
    "2 min 30 s later"
   ],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 700,
    "kb": 27
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 25
   }
  },
  {
   "id": "setup-running",
   "view": "setup",
   "state": "setup-running",
   "title": "Running",
   "how": "Setup in progress, third step.",
   "note": "",
   "actions": [],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 700,
    "kb": 24
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 23
   }
  },
  {
   "id": "setup-failed",
   "view": "setup",
   "state": "setup-failed",
   "title": "Failed",
   "how": "A step failed; the form returns with Try again.",
   "note": "",
   "actions": [],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 700,
    "kb": 35
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 33
   }
  },
  {
   "id": "setup-done",
   "view": "setup",
   "state": "setup-done",
   "title": "All set",
   "how": "Every step finished.",
   "note": "",
   "actions": [],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 700,
    "kb": 30
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 29
   }
  },
  {
   "id": "setup-existing",
   "view": "setup",
   "state": "setup-existing",
   "title": "Existing setup found",
   "how": "A setup made by hand from the old guide is detected.",
   "note": "",
   "actions": [],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 770,
    "kb": 47
   },
   "dark": {
    "w": 520,
    "h": 770,
    "kb": 42
   }
  },
  {
   "id": "setup-adopted",
   "view": "setup",
   "state": "setup-existing",
   "title": "Taken over",
   "how": "Take over clicked; every item replaced.",
   "note": "",
   "actions": [
    "click Take over"
   ],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 704,
    "kb": 44
   },
   "dark": {
    "w": 520,
    "h": 704,
    "kb": 39
   }
  },
  {
   "id": "setup-adopt-error",
   "view": "setup",
   "state": "setup-adopt-error",
   "title": "Take over failed",
   "how": "A legacy file was locked by another process.",
   "note": "",
   "actions": [
    "click Take over"
   ],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 826,
    "kb": 51
   },
   "dark": {
    "w": 520,
    "h": 826,
    "kb": 46
   }
  },
  {
   "id": "setup-installed",
   "view": "setup",
   "state": "setup-installed",
   "title": "Opened after setup",
   "how": "Set up... opened again on a machine Docker Zero already set up.",
   "note": "",
   "actions": [],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 700,
    "kb": 39
   },
   "dark": {
    "w": 520,
    "h": 700,
    "kb": 35
   }
  }
 ],
 "skipped": [],
 "real": []
};
