// Shared helpers for the gallery mocks. Copy into the repository's gallery
// mock folder; drop the comments where the repository forbids them.

/** The scenario this page shows, from ?state=<name>. */
export const stateName = new URLSearchParams(location.search).get("state") ?? "";

/** How the mocked calls of one scenario behave, by "<module>.<function>". */
export interface Behaviour {
  /** Calls that reject, with the error text the real backend would give. */
  fail?: Record<string, string>;
  /** Calls that never settle, to hold a busy or loading state on screen. */
  hang?: string[];
  /** Delay before a call settles, in ms; long enough to see a spinner, short enough to settle before the shot. */
  delayMs?: number;
}

/** Returns a call() that resolves, rejects or hangs as the behaviour says. */
export function caller(b: Behaviour) {
  const delay = b.delayMs ?? 150;
  return function call<T>(name: string, value: T | (() => T)): Promise<T> {
    if (b.hang?.includes(name)) return new Promise<T>(() => {});
    const failure = b.fail?.[name];
    if (failure) return new Promise<T>((_, reject) => setTimeout(() => reject(new Error(failure)), delay));
    return new Promise<T>((resolve) =>
      setTimeout(() => resolve(typeof value === "function" ? (value as () => T)() : value), delay),
    );
  };
}

/**
 * The type of a mock for every exported function of a real module. Write a
 * mock as `const mock = { ... } satisfies Mocked<typeof Real>` with
 * `import type * as Real from "<real module>"`: a function added to or changed
 * in the real module then fails the type check instead of the gallery.
 */
export type Mocked<T> = {
  [K in keyof T]: T[K] extends (...args: infer A) => PromiseLike<infer R>
    ? (...args: A) => Promise<R>
    : T[K] extends (...args: infer A) => infer R
      ? (...args: A) => R
      : T[K];
};
