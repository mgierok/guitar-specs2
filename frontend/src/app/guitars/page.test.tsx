import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import GuitarsPage from "./page";

vi.mock("@/lib/api", () => ({
  fetchGuitars: vi.fn(async () => []),
}));

describe("GuitarsPage", () => {
  it("renders empty state when no guitars", async () => {
    render(await GuitarsPage());

    expect(
      screen.getByText(/No guitars found yet/i)
    ).toBeInTheDocument();
  });
});
