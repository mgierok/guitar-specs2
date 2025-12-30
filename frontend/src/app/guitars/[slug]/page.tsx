import type { Metadata } from "next";

import { fetchGuitar } from "@/lib/api";

type PageProps = {
  params: {
    slug: string;
  };
};

function titleFromSlug(slug: string) {
  return slug
    .split("-")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

export async function generateMetadata({
  params,
}: PageProps): Promise<Metadata> {
  const title = titleFromSlug(params.slug);

  return {
    title,
    alternates: {
      canonical: `/guitars/${params.slug}`,
    },
  };
}

export default async function GuitarDetailPage({ params }: PageProps) {
  const guitar = await fetchGuitar(params.slug);

  return (
    <div className="min-h-screen bg-zinc-950 text-white">
      <main className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-6 pb-16 pt-12 md:px-10">
        <header>
          <p className="text-sm uppercase tracking-[0.3em] text-purple-300">
            Guitar Details
          </p>
          <h1 className="mt-3 text-3xl font-semibold md:text-4xl">
            {guitar.name}
          </h1>
          <p className="mt-2 max-w-2xl text-sm text-purple-100/70">
            {guitar.brand} · {guitar.model} · {guitar.type}
          </p>
        </header>

        <section className="grid gap-4 rounded-2xl border border-white/10 bg-white/5 p-6">
          <div className="flex flex-col gap-2 text-sm text-purple-100/70">
            {guitar.description && <span>{guitar.description}</span>}
            {Object.entries(guitar.specs).map(([key, value]) => (
              <span key={key}>
                {key.replace(/_/g, \" \")}: {value}
              </span>
            ))}
          </div>
        </section>
      </main>
    </div>
  );
}
