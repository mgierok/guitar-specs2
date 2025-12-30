import type { Metadata } from "next";
import { JetBrains_Mono, Space_Grotesk } from "next/font/google";
import "./globals.css";

const spaceGrotesk = Space_Grotesk({
  variable: "--font-brand",
  subsets: ["latin"],
});

const jetBrainsMono = JetBrains_Mono({
  variable: "--font-code",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL("https://www.guitar-specs.com"),
  title: {
    default: "Guitar-Specs",
    template: "%s | Guitar-Specs",
  },
  description:
    "Search, compare, and explore detailed specifications for electric, acoustic, and classical guitars.",
  alternates: {
    canonical: "/",
  },
  openGraph: {
    title: "Guitar-Specs",
    description:
      "Search, compare, and explore detailed specifications for electric, acoustic, and classical guitars.",
    url: "/",
    siteName: "Guitar-Specs",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Guitar-Specs",
    description:
      "Search, compare, and explore detailed specifications for electric, acoustic, and classical guitars.",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${spaceGrotesk.variable} ${jetBrainsMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
