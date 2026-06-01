import { describe, it, expect } from "vitest";
import { cn, formatDate, maskKey } from "./utils";

describe("utils", () => {
  it("cn merges classes", () => {
    expect(cn("a", "b")).toBe("a b");
    expect(cn("p-2", "p-4")).toBe("p-4");
  });

  it("formatDate handles dates", () => {
    const out = formatDate("2025-01-01T00:00:00Z");
    expect(typeof out).toBe("string");
    expect(out.length).toBeGreaterThan(0);
  });

  it("maskKey masks long keys", () => {
    expect(maskKey("sk-abcdefghijklmnop")).toContain("...");
    expect(maskKey("")).toBe("");
    expect(maskKey("short")).toBe("****");
  });
});
