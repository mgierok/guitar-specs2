export const metadata = {
  title: "Compare",
  alternates: {
    canonical: "/compare",
  },
};

export default function ComparePage() {
  return (
    <div className="min-h-screen bg-zinc-950 text-white">
      <main className="mx-auto flex w-full max-w-6xl flex-col gap-6 px-6 pb-16 pt-12 md:px-10">
        <header>
          <p className="text-sm uppercase tracking-[0.3em] text-purple-300">
            Compare
          </p>
          <h1 className="mt-3 text-3xl font-semibold md:text-4xl">
            Compare guitars side by side
          </h1>
          <p className="mt-2 max-w-2xl text-sm text-purple-100/70">
            Select two or more guitars to compare specs across pickups, scale
            length, and materials.
          </p>
        </header>

        <section className="rounded-2xl border border-white/10 bg-white/5 p-6">
          <p className="text-sm text-purple-100/70">
            Comparison grid placeholder. Wire this to the API once endpoints
            are ready.
          </p>
        </section>
      </main>
    </div>
  );
}
