import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import RootLayout, { metadata } from "./layout";

describe("RootLayout", () => {
  it("exposes default metadata", () => {
    expect(metadata).toBeDefined();
    expect(metadata?.title).toBeDefined();
    expect(metadata?.description).toBeDefined();
  });

  it("renders children", () => {
    render(
      <RootLayout>
        <div>Child content</div>
      </RootLayout>
    );

    expect(screen.getByText(/Child content/i)).toBeInTheDocument();
  });
});
