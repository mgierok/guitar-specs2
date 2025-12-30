import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import GuitarDetailPage from "./page";

vi.mock("@/lib/api", () => ({
  fetchGuitar: vi.fn(async () => ({
    id: "g1",
    slug: "fender-player-stratocaster",
    name: "Player Stratocaster",
    brand: "Fender",
    model: "Stratocaster",
    type: "electric",
    specs: { scale_length: "25.5" },
    media: [],
  })),
}));

describe("GuitarDetailPage", () => {
  it("renders guitar detail content", async () => {
    render(
      await GuitarDetailPage({
        params: { slug: "fender-player-stratocaster" },
      })
    );

    expect(screen.getByText(/Player Stratocaster/i)).toBeInTheDocument();
    expect(screen.getByText(/Fender/i)).toBeInTheDocument();
    expect(screen.getByText(/scale length/i)).toBeInTheDocument();
  });
});
