"""Render the animated README hero: steps motion.html through one loop and saves an animated WebP per theme.

Run from the repository root:
  uvx --from playwright --with pillow python -I design/readme/source/frames.py --skill-dir <frontdoor-illustrate folder>

Needs Edge or Chrome on the machine and network access for the web fonts. Output: hero-motion.webp and
hero-motion-dark.webp in design/readme/asset, rounded and bordered like the still figures.
"""
import argparse
import functools
import http.server
import importlib.util
import io
import tempfile
import threading
from pathlib import Path

from PIL import Image
from playwright.sync_api import sync_playwright

ap = argparse.ArgumentParser()
ap.add_argument("--skill-dir", required=True, help="the frontdoor-illustrate skill folder (for render.py helpers)")
ap.add_argument("--frames", type=int, default=30)
ap.add_argument("--ms", type=int, default=80, help="duration of one frame")
ap.add_argument("--quality", type=int, default=68)
ap.add_argument("--out", default="design/readme/asset")
a = ap.parse_args()

spec = importlib.util.spec_from_file_location("render", Path(a.skill_dir) / "scripts" / "render.py")
render = importlib.util.module_from_spec(spec)
spec.loader.exec_module(render)

root = Path(".").resolve()
out = Path(a.out).resolve()
out.mkdir(parents=True, exist_ok=True)


class Quiet(http.server.SimpleHTTPRequestHandler):
    def log_message(self, *args):
        pass


server = render.QuietServer(("127.0.0.1", 0), functools.partial(Quiet, directory=str(root)))
threading.Thread(target=server.serve_forever, daemon=True).start()
port = server.server_address[1]

with sync_playwright() as p, tempfile.TemporaryDirectory() as tmp:
    browser = render.launch(p, None)
    page = browser.new_page(viewport={"width": 2600, "height": 1400})
    for theme in ("light", "dark"):
        frames = []
        page.goto(f"http://127.0.0.1:{port}/design/readme/source/motion.html?theme={theme}", wait_until="networkidle")
        page.evaluate("document.fonts.ready")
        for i in range(a.frames):
            page.evaluate(f"setT({i / a.frames})")
            el = page.locator("#hero-motion")
            line = render.parse_color(
                el.evaluate("e => getComputedStyle(e).getPropertyValue('--cardline')"),
                render.FALLBACK_BORDER[theme])
            raw = Path(tmp) / f"{theme}-{i}.png"
            el.screenshot(path=str(raw), scale="css")
            framed = Path(tmp) / f"{theme}-{i}.webp"
            render.round_and_border(raw, framed, line)
            frames.append(Image.open(framed).convert("RGBA"))
        dst = out / ("hero-motion.webp" if theme == "light" else "hero-motion-dark.webp")
        frames[0].save(dst, "WEBP", save_all=True, append_images=frames[1:], duration=a.ms, loop=0,
                       quality=a.quality, method=6, minimize_size=False)
        print(f"{dst.name}  {frames[0].size[0]}x{frames[0].size[1]}  {len(frames)} frames  {dst.stat().st_size / 1024:.0f} KB")
    browser.close()
server.shutdown()
