export type GuitarListItem = {
  id: string;
  slug: string;
  name: string;
  brand: string;
  model: string;
  type: string;
  year?: number;
  thumbnail?: string;
};

export type GuitarDetail = {
  id: string;
  slug: string;
  name: string;
  brand: string;
  model: string;
  type: string;
  year?: number;
  description?: string;
  specs: Record<string, string>;
  media: { kind: string; url: string }[];
};

const API_BASE =
  process.env.API_BASE ||
  process.env.NEXT_PUBLIC_API_BASE ||
  "http://localhost:8080/api/v1";

export async function fetchGuitars(): Promise<GuitarListItem[]> {
  const response = await fetch(`${API_BASE}/guitars`, {
    next: { revalidate: 60 },
  });

  if (!response.ok) {
    throw new Error("Failed to load guitars");
  }

  const data = await response.json();
  return data.items ?? [];
}

export async function fetchGuitar(slug: string): Promise<GuitarDetail> {
  const response = await fetch(`${API_BASE}/guitars/${slug}`, {
    next: { revalidate: 60 },
  });

  if (!response.ok) {
    throw new Error("Failed to load guitar");
  }

  return response.json();
}
