import type * as Real from "../../bindings/github.com/the-khiem7/docker-zero/internal/ui/setupservice";
import { call, scenario, type Mocked } from "../fixtures";

const mock = {
  Start: () => call("SetupService.Start", undefined),
  State: () => call("SetupService.State", scenario.setup ?? { steps: null, log: null, running: false, done: false }),
} satisfies Mocked<typeof Real>;

export const { Start, State } = mock;
