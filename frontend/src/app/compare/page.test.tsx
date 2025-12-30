import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import ComparePage from "./page";

describe("ComparePage", () => {
  it("renders compare heading", () => {
    render(<ComparePage />);

    expect(
      screen.getByText(/Compare guitars side by side/i)
    ).toBeInTheDocument();
  });
});
