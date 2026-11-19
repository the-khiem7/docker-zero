window.GALLERY = {
 "title": "Docker Zero Settings redesign board, V2 control deck",
 "lede": "The Settings window as a control deck: a grid of big visual controls (power disc, app tiles, tide slider, certificate stamp) with drawn icons, chips instead of paragraphs and long explanations behind the i button. Chosen from three layouts of the exploration board. Same mocked services and state ids as before. Mark each Keep, Revise or Discuss.",
 "regenerate": [
  "cd frontend && npm run gallery",
  "uvx --from playwright --with pillow python -I <frontdoor-ui-gallery>/scripts/shoot.py render design/ui-gallery/V2-setting/gallery.json"
 ],
 "generated": "2026-11-19",
 "commit": "e966701",
 "storageKey": "docker-zero-ui-gallery3",
 "views": {
  "settings": {
   "scroller": "main",
   "title": "Settings",
   "size": [
    520,
    680
   ],
   "source": "frontend/src/Settings.svelte",
   "opens": "Gear in the popup, or Settings... in the tray menu.",
   "url": "/?view=settings"
  },
  "maintenance": {
   "scroller": "main",
   "title": "Settings: Portainer and Danger zone",
   "size": [
    520,
    680
   ],
   "source": "frontend/src/Maintenance.svelte",
   "opens": "Bottom of the Settings window, below the certificate.",
   "url": "/?view=settings",
   "short": "Portainer and Danger zone"
  }
 },
 "states": [
  {
   "id": "settings-loading",
   "view": "settings",
   "state": "settings-loading",
   "title": "Loading",
   "how": "The first moment after the window opens.",
   "note": "",
   "actions": [],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 680,
    "kb": 7
   },
   "dark": {
    "w": 520,
    "h": 680,
    "kb": 6
   }
  },
  {
   "id": "settings",
   "view": "settings",
   "state": "settings",
   "title": "All settings",
   "how": "Typical state after setup: auto start, shortcuts, keepalive, trusted certificate.",
   "note": "",
   "actions": [],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 1021,
    "kb": 75
   },
   "dark": {
    "w": 520,
    "h": 1021,
    "kb": 70
   }
  },
  {
   "id": "settings-off",
   "view": "settings",
   "state": "settings-off",
   "title": "Everything off",
   "how": "Fresh portable run: no auto start, no shortcuts, WSL default idle.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 520,
    "kb": 35
   },
   "dark": {
    "w": 520,
    "h": 520,
    "kb": 32
   }
  },
  {
   "id": "settings-wslconfig",
   "view": "settings",
   "state": "settings-wslconfig",
   "title": ".wslconfig mode pending",
   "how": "Never stop idle distros chosen; WSL has not restarted yet.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 212,
    "kb": 20
   },
   "dark": {
    "w": 520,
    "h": 212,
    "kb": 19
   }
  },
  {
   "id": "settings-wslconfig-confirm",
   "view": "settings",
   "state": "settings-wslconfig",
   "title": "Confirm Apply now",
   "how": "Apply now... clicked.",
   "note": "",
   "actions": [
    "click Apply now..."
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 240,
    "kb": 24
   },
   "dark": {
    "w": 520,
    "h": 240,
    "kb": 22
   }
  },
  {
   "id": "settings-browser",
   "view": "settings",
   "state": "settings-browser",
   "title": "Custom browser, no Chrome",
   "how": "Other browser chosen with Browse...; Chrome not installed.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 195,
    "kb": 15
   },
   "dark": {
    "w": 520,
    "h": 195,
    "kb": 14
   }
  },
  {
   "id": "settings-untrusted",
   "view": "settings",
   "state": "settings-untrusted",
   "title": "Certificate not trusted",
   "how": "After setup, before Install (trust).",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 296,
    "kb": 20
   },
   "dark": {
    "w": 520,
    "h": 296,
    "kb": 20
   }
  },
  {
   "id": "settings-no-cert",
   "view": "settings",
   "state": "settings-no-cert",
   "title": "No certificate",
   "how": "Distro running but setup did not reach the certificate step.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 296,
    "kb": 15
   },
   "dark": {
    "w": 520,
    "h": 296,
    "kb": 13
   }
  },
  {
   "id": "settings-cert-unread",
   "view": "settings",
   "state": "settings-cert-unread",
   "title": "Distro stopped",
   "how": "Settings opened while Ubuntu-24.04 is stopped; the certificate is not read.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 296,
    "kb": 15
   },
   "dark": {
    "w": 520,
    "h": 296,
    "kb": 14
   }
  },
  {
   "id": "settings-hint-open",
   "view": "settings",
   "state": "settings",
   "title": "Explanation behind the i button",
   "how": "The i of the Certificate section opened: the long text is read on request, not printed on the page.",
   "note": "",
   "actions": [
    "click About Certificate"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 296,
    "kb": 26
   },
   "dark": {
    "w": 520,
    "h": 296,
    "kb": 25
   }
  },
  {
   "id": "settings-hint-shortcuts",
   "view": "settings",
   "state": "settings",
   "title": "Explanation of a card on the left edge",
   "how": "The i of the Shortcuts card opened: the bubble stays inside the window and the install folder sits in a monospace pill.",
   "note": "",
   "actions": [
    "click About Shortcuts"
   ],
   "crop": "window",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 680,
    "kb": 57
   },
   "dark": {
    "w": 520,
    "h": 680,
    "kb": 53
   }
  },
  {
   "id": "settings-busy",
   "view": "settings",
   "state": "settings-busy",
   "title": "Working",
   "how": "A setting is being applied.",
   "note": "",
   "actions": [
    "click Desktop"
   ],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 1074,
    "kb": 76
   },
   "dark": {
    "w": 520,
    "h": 1074,
    "kb": 71
   }
  },
  {
   "id": "settings-problems",
   "view": "settings",
   "state": "settings-problems",
   "title": "Errors",
   "how": "A setting failed to apply, and two settings could not be read.",
   "note": "",
   "actions": [
    "click Desktop"
   ],
   "crop": "full",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 1245,
    "kb": 89
   },
   "dark": {
    "w": 520,
    "h": 1245,
    "kb": 83
   }
  },
  {
   "id": "maint",
   "view": "maintenance",
   "state": "maint",
   "title": "Idle",
   "how": "Settings window scrolled to the bottom.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 444,
    "kb": 34
   },
   "dark": {
    "w": 520,
    "h": 444,
    "kb": 32
   }
  },
  {
   "id": "maint-updating",
   "view": "maintenance",
   "state": "maint-updating",
   "title": "Updating",
   "how": "Update Portainer clicked; pull in progress.",
   "note": "",
   "actions": [
    "click Update Portainer"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 444,
    "kb": 31
   },
   "dark": {
    "w": 520,
    "h": 444,
    "kb": 29
   }
  },
  {
   "id": "maint-updated",
   "view": "maintenance",
   "state": "maint-updated",
   "title": "Updated",
   "how": "A newer lts image was pulled and Portainer restarted.",
   "note": "",
   "actions": [
    "click Update Portainer"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 444,
    "kb": 34
   },
   "dark": {
    "w": 520,
    "h": 444,
    "kb": 32
   }
  },
  {
   "id": "maint-current",
   "view": "maintenance",
   "state": "maint-current",
   "title": "Already latest",
   "how": "Update clicked with the image current.",
   "note": "",
   "actions": [
    "click Update Portainer"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 444,
    "kb": 34
   },
   "dark": {
    "w": 520,
    "h": 444,
    "kb": 32
   }
  },
  {
   "id": "maint-update-failed",
   "view": "maintenance",
   "state": "maint-update-failed",
   "title": "Update failed",
   "how": "Docker Hub unreachable.",
   "note": "",
   "actions": [
    "click Update Portainer"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 510,
    "kb": 37
   },
   "dark": {
    "w": 520,
    "h": 510,
    "kb": 36
   }
  },
  {
   "id": "maint-reset-open",
   "view": "maintenance",
   "state": "maint",
   "title": "Reset form",
   "how": "Reset admin password... clicked.",
   "note": "",
   "actions": [
    "click Reset admin password..."
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 654,
    "kb": 46
   },
   "dark": {
    "w": 520,
    "h": 654,
    "kb": 43
   }
  },
  {
   "id": "maint-reset-invalid",
   "view": "maintenance",
   "state": "maint",
   "title": "Reset form: short password",
   "how": "Reset password clicked with the fields empty.",
   "note": "",
   "actions": [
    "click Reset admin password...",
    "click Reset password"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 716,
    "kb": 52
   },
   "dark": {
    "w": 520,
    "h": 716,
    "kb": 48
   }
  },
  {
   "id": "maint-reset-wrong-user",
   "view": "maintenance",
   "state": "maint-reset-wrong-user",
   "title": "Reset: wrong username",
   "how": "Password replaced, but the username check failed.",
   "note": "",
   "actions": [
    "click Reset admin password...",
    "type Admin username",
    "type New password",
    "type Repeat new password",
    "click Reset password"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 720,
    "kb": 51
   },
   "dark": {
    "w": 520,
    "h": 720,
    "kb": 48
   }
  },
  {
   "id": "maint-reset-done",
   "view": "maintenance",
   "state": "maint-reset-done",
   "title": "Reset done",
   "how": "Password replaced and checked.",
   "note": "",
   "actions": [
    "click Reset admin password...",
    "type New password",
    "type Repeat new password",
    "click Reset password"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 444,
    "kb": 34
   },
   "dark": {
    "w": 520,
    "h": 444,
    "kb": 32
   }
  },
  {
   "id": "maint-plan",
   "view": "maintenance",
   "state": "maint",
   "title": "Uninstall plan",
   "how": "Uninstall all... clicked; backup on, no file yet.",
   "note": "",
   "actions": [
    "click Uninstall all..."
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 531,
    "kb": 41
   },
   "dark": {
    "w": 520,
    "h": 531,
    "kb": 37
   }
  },
  {
   "id": "maint-plan-file",
   "view": "maintenance",
   "state": "maint-plan-file",
   "title": "Uninstall plan with backup file",
   "how": "Choose file... used.",
   "note": "",
   "actions": [
    "click Uninstall all...",
    "click Choose file..."
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 561,
    "kb": 42
   },
   "dark": {
    "w": 520,
    "h": 561,
    "kb": 39
   }
  },
  {
   "id": "maint-confirm",
   "view": "maintenance",
   "state": "maint",
   "title": "Uninstall confirmation",
   "how": "Backup unticked, Continue clicked.",
   "note": "",
   "actions": [
    "click Uninstall all...",
    "click Save a backup of Ubuntu-24.04 first",
    "click Continue"
   ],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 444,
    "kb": 33
   },
   "dark": {
    "w": 520,
    "h": 444,
    "kb": 31
   }
  },
  {
   "id": "maint-uninstalling",
   "view": "maintenance",
   "state": "maint-uninstalling",
   "title": "Uninstalling",
   "how": "Uninstall now clicked; backup export running.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 521,
    "kb": 38
   },
   "dark": {
    "w": 520,
    "h": 521,
    "kb": 36
   }
  },
  {
   "id": "maint-uninstall-failed",
   "view": "maintenance",
   "state": "maint-uninstall-failed",
   "title": "Uninstall failed",
   "how": "The Windows certificate dialog was answered No.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 522,
    "kb": 36
   },
   "dark": {
    "w": 520,
    "h": 522,
    "kb": 33
   }
  },
  {
   "id": "maint-uninstalled",
   "view": "maintenance",
   "state": "maint-uninstalled",
   "title": "Uninstalled",
   "how": "Every step finished; only Quit is left.",
   "note": "",
   "actions": [],
   "crop": "cards",
   "themes": [
    "light",
    "dark"
   ],
   "light": {
    "w": 520,
    "h": 443,
    "kb": 35
   },
   "dark": {
    "w": 520,
    "h": 443,
    "kb": 33
   }
  }
 ],
 "skipped": [],
 "real": []
};
