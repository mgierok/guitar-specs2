import Link from "next/link";

export const metadata = {
  title: "Home",
  alternates: {
    canonical: "/",
  },
};

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-purple-950 via-purple-900 to-zinc-950 text-white">
      <main className="mx-auto flex w-full max-w-6xl flex-col gap-16 px-6 pb-20 pt-16 md:px-10">
        <section className="rounded-3xl border border-white/10 bg-white/5 p-10 backdrop-blur">
          <p className="text-sm uppercase tracking-[0.3em] text-purple-200">
            Guitar-Specs
          </p>
          <h1 className="mt-4 text-4xl font-semibold tracking-tight md:text-5xl">
            Search, compare, and explore guitars with confidence.
          </h1>
          <p className="mt-4 max-w-2xl text-lg text-purple-100/80">
            Find full specs across brands, filter by key dimensions, and build
            head-to-head comparisons before you buy or build.
          </p>
          <form action="/guitars" className="mt-8 flex flex-col gap-3 md:flex-row">
            <label className="sr-only" htmlFor="search">
              Search guitars
            </label>
            <input
              id="search"
              type="text"
              name="q"
              placeholder="Search by brand, model, or spec"
              className="h-12 flex-1 rounded-full border border-white/20 bg-white/10 px-5 text-white placeholder:text-purple-100/60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-purple-300"
            />
            <button
              type="submit"
              className="h-12 rounded-full bg-purple-400 px-6 font-medium text-purple-950 transition hover:bg-purple-300"
            >
              Search
            </button>
          </form>
        </section>

        <section className="grid gap-6 md:grid-cols-3">
          <div className="rounded-2xl border border-white/10 bg-white/5 p-6">
            <h2 className="text-xl font-semibold">Browse the catalog</h2>
            <p className="mt-2 text-sm text-purple-100/70">
              Filter by body shape, pickups, scale length, and finish.
            </p>
            <Link
              className="mt-4 inline-flex text-sm font-medium text-purple-200 hover:text-purple-100"
              href="/guitars"
            >
              Go to list
            </Link>
          </div>
          <div className="rounded-2xl border border-white/10 bg-white/5 p-6">
            <h2 className="text-xl font-semibold">Compare instruments</h2>
            <p className="mt-2 text-sm text-purple-100/70">
              Line up two or more guitars and compare specs side by side.
            </p>
            <Link
              className="mt-4 inline-flex text-sm font-medium text-purple-200 hover:text-purple-100"
              href="/compare"
            >
              Start comparing
            </Link>
          </div>
          <div className="rounded-2xl border border-white/10 bg-white/5 p-6">
            <h2 className="text-xl font-semibold">Dive into details</h2>
            <p className="mt-2 text-sm text-purple-100/70">
              Every guitar gets a full, searchable spec sheet.
            </p>
            <Link
              className="mt-4 inline-flex text-sm font-medium text-purple-200 hover:text-purple-100"
              href="/guitars/fender-player-stratocaster"
            >
              See an example
            </Link>
          </div>
        </section>
      </main>
    </div>
  );
}
