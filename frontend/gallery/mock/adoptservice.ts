import type * as Real from "../../bindings/github.com/the-khiem7/docker-zero/internal/ui/adoptservice";
import { call, scenario, type Mocked } from "../fixtures";

let adopted = false;

const mock = {
  Adopt: (items: string[] | null) =>
    call("AdoptService.Adopt", () => {
      adopted = true;
      return items;
    }),
  Detect: () =>
    call("AdoptService.Detect", () =>
      adopted && scenario.reportAfterAdopt
        ? scenario.reportAfterAdopt
        : (scenario.report ?? { distro: "running", container: "running", certificate: true, findings: [] }),
    ),
} satisfies Mocked<typeof Real>;

export const { Adopt, Detect } = mock;
