import { Events } from "@wailsio/runtime";
import * as Svc from "../../bindings/github.com/the-khiem7/docker-zero/internal/ui/service";
import type { View } from "../../bindings/github.com/the-khiem7/docker-zero/internal/ui/models";

export const service = Svc;

/** Live status shared by every component of this window. */
export const live = $state<{ view: View | null; error: string }>({ view: null, error: "" });

Events.On("status", (e) => {
  live.view = e.data;
});

Svc.Status().then((v) => {
  if (live.view === null && v.checkedAt) live.view = v;
});

/** Runs a bound action and keeps its error message for the UI. */
export async function run(action: () => Promise<void>): Promise<void> {
  live.error = "";
  try {
    await action();
  } catch (err) {
    live.error = err instanceof Error ? err.message : String(err);
  }
}

export const levelLabel: Record<string, string> = {
  running: "Running",
  partial: "Starting",
  stopped: "Stopped",
  failed: "Needs setup",
};
