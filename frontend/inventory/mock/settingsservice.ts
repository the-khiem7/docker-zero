import type * as Real from "../../bindings/github.com/the-khiem7/docker-zero/internal/ui/settingsservice";
import { baseSettings, call, scenario, type Mocked } from "../fixtures";

const view = { ...baseSettings, ...scenario.settings };

const mock = {
  ApplyWSLConfig: () => call("SettingsService.ApplyWSLConfig", undefined),
  Get: () => call("SettingsService.Get", view),
  PickBrowser: () => call("SettingsService.PickBrowser", scenario.pick ?? ""),
  RegenerateCert: () => call("SettingsService.RegenerateCert", undefined),
  SetAutoStart: () => call("SettingsService.SetAutoStart", undefined),
  SetBrowser: () => call("SettingsService.SetBrowser", undefined),
  SetIdleMode: () => call("SettingsService.SetIdleMode", false),
  SetShortcuts: () => call("SettingsService.SetShortcuts", undefined),
  TrustCert: () => call("SettingsService.TrustCert", undefined),
  UntrustCert: () => call("SettingsService.UntrustCert", undefined),
} satisfies Mocked<typeof Real>;

export const {
  ApplyWSLConfig,
  Get,
  PickBrowser,
  RegenerateCert,
  SetAutoStart,
  SetBrowser,
  SetIdleMode,
  SetShortcuts,
  TrustCert,
  UntrustCert,
} = mock;
