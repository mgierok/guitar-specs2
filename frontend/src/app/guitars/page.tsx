import Link from "next/link";

import { fetchGuitars } from "@/lib/api";

export const dynamic = "force-dynamic";

export const metadata = {
  title: "Guitars",
  alternates: {
    canonical: "/guitars",
  },
};

export default async function GuitarsPage() {
  const guitars = await fetchGuitars();

  return (
    <div className="min-h-screen bg-zinc-950 text-white">
      <main className="mx-auto flex w-full max-w-6xl flex-col gap-8 px-6 pb-16 pt-12 md:px-10">
        <header>
          <p className="text-sm uppercase tracking-[0.3em] text-purple-300">
            Guitars
          </p>
          <h1 className="mt-3 text-3xl font-semibold md:text-4xl">
            Browse the full catalog
          </h1>
          <p className="mt-2 max-w-2xl text-sm text-purple-100/70">
            Filter by brand, pickup configuration, scale length, and more.
          </p>
        </header>

        <section className="grid gap-4 md:grid-cols-3">
          {guitars.length === 0 ? (
            <div className="rounded-2xl border border-white/10 bg-white/5 p-6 text-sm text-purple-100/70">
              No guitars found yet. Seed the database to populate the catalog.
            </div>
          ) : (
            guitars.map((guitar) => (
              <Link
                key={guitar.slug}
                href={`/guitars/${guitar.slug}`}
                className="rounded-2xl border border-white/10 bg-white/5 p-5 transition hover:border-purple-400/60"
              >
                <h2 className="text-lg font-semibold">{guitar.name}</h2>
                <p className="mt-2 text-sm text-purple-100/70">
                  {guitar.brand} · {guitar.model}
                </p>
              </Link>
            ))
          )}
        </section>
      </main>
    </div>
  );
}
