import type * as Real from "../../bindings/github.com/the-khiem7/docker-zero/internal/ui/uninstallservice";
import { call, scenario, type Mocked } from "../fixtures";

const mock = {
  PickBackup: () => call("UninstallService.PickBackup", scenario.pick ?? ""),
  Quit: () => call("UninstallService.Quit", undefined),
  Start: () => call("UninstallService.Start", undefined),
  State: () =>
    call("UninstallService.State", scenario.uninstall ?? { steps: null, log: null, running: false, done: false }),
} satisfies Mocked<typeof Real>;

export const { PickBackup, Quit, Start, State } = mock;
