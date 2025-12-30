import "@testing-library/jest-dom";
import { vi } from "vitest";

vi.mock("next/font/google", () => ({
  Space_Grotesk: () => ({ variable: "--font-brand" }),
  JetBrains_Mono: () => ({ variable: "--font-code" }),
}));
