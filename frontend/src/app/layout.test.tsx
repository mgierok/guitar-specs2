import { renderToString } from "react-dom/server";
import { describe, expect, it } from "vitest";

import RootLayout, { metadata } from "./layout";

describe("RootLayout", () => {
  it("exposes default metadata", () => {
    expect(metadata).toBeDefined();
    expect(metadata?.title).toBeDefined();
    expect(metadata?.description).toBeDefined();
  });

  it("renders children", () => {
    const html = renderToString(
      <RootLayout>
        <div>Child content</div>
      </RootLayout>
    );

    expect(html).toContain("Child content");
  });
});
