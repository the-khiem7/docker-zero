import type * as Real from "../../bindings/github.com/the-khiem7/docker-zero/internal/ui/service";
import { call, runningView, scenario, type Mocked } from "../fixtures";

const view = { ...runningView, ...scenario.view };

const mock = {
  HidePopup: () => call("Service.HidePopup", undefined),
  OpenPortainer: () => call("Service.OpenPortainer", undefined),
  OpenSettings: () => call("Service.OpenSettings", undefined),
  OpenSetup: () => call("Service.OpenSetup", undefined),
  Refresh: () => call("Service.Refresh", undefined),
  SetCornerBox: () => call("Service.SetCornerBox", undefined),
  ShutdownWSL: () => call("Service.ShutdownWSL", undefined),
  Start: () => call("Service.Start", undefined),
  Status: () => call("Service.Status", view),
  StopDistro: () => call("Service.StopDistro", undefined),
} satisfies Mocked<typeof Real>;

export const { HidePopup, OpenPortainer, OpenSettings, OpenSetup, Refresh, SetCornerBox, ShutdownWSL, Start, Status, StopDistro } =
  mock;
