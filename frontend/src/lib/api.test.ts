import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { fetchGuitar, fetchGuitars } from "./api";

type MockResponse = {
  ok: boolean;
  json: () => Promise<unknown>;
};

describe("api client", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("fetchGuitars returns items", async () => {
    const mock = {
      ok: true,
      json: async () => ({ items: [{ id: "g1", name: "Test", slug: "test", brand: "Brand", model: "Model", type: "electric" }] }),
    } satisfies MockResponse;

    (fetch as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mock);

    const items = await fetchGuitars();
    expect(items).toHaveLength(1);
    expect(items[0].name).toBe("Test");
  });

  it("fetchGuitars throws on error", async () => {
    const mock = { ok: false, json: async () => ({}) } satisfies MockResponse;
    (fetch as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mock);

    await expect(fetchGuitars()).rejects.toThrow(/Failed to load guitars/i);
  });

  it("fetchGuitar returns detail", async () => {
    const mock = {
      ok: true,
      json: async () => ({ id: "g1", slug: "test", name: "Test", brand: "Brand", model: "Model", type: "electric", specs: {}, media: [] }),
    } satisfies MockResponse;

    (fetch as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mock);

    const detail = await fetchGuitar("test");
    expect(detail.slug).toBe("test");
  });

  it("fetchGuitar throws on error", async () => {
    const mock = { ok: false, json: async () => ({}) } satisfies MockResponse;
    (fetch as unknown as ReturnType<typeof vi.fn>).mockResolvedValue(mock);

    await expect(fetchGuitar("test")).rejects.toThrow(/Failed to load guitar/i);
  });
});
