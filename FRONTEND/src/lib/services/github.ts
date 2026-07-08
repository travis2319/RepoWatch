import { API } from '$lib/constants/index';
import type { Repo, AccessRepo, ValidateData } from '$lib/types/github';

export async function fetchRepos(owner: string): Promise<Repo[]> {
  const r = await fetch(`${API}/repos?owner=${owner}`);
  if (!r.ok) throw new Error(`${r.status}`);
  const d = await r.json();
  return d.data ?? d ?? [];
}

export async function fetchUserAccess(owner: string, user: string): Promise<AccessRepo[]> {
  const r = await fetch(`${API}/check`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ owner, user }),
  });
  if (!r.ok) throw new Error(`${r.status}`);
  const d = await r.json();
  return (d.data ?? []).filter((x: AccessRepo) => x.hasAccess);
}

export async function fetchValidate(): Promise<ValidateData> {
  const r = await fetch(`${API}/github/validate`);
  if (!r.ok) throw new Error(`${r.status}`);
  const d = await r.json();
  return d.data ?? d;
}