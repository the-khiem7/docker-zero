import type * as Real from "../../bindings/github.com/the-khiem7/docker-zero/internal/ui/maintainservice";
import { call, scenario, type Mocked } from "../fixtures";

const mock = {
  ResetAdmin: () => call("MaintainService.ResetAdmin", undefined),
  UpdatePortainer: () => call("MaintainService.UpdatePortainer", scenario.updateReplaced ?? false),
} satisfies Mocked<typeof Real>;

export const { ResetAdmin, UpdatePortainer } = mock;
