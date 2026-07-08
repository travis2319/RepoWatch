import { writable } from 'svelte/store';
import { fetchUserAccess } from '$lib/services/github';
import type { AccessRepo } from '$lib/types/github';

export const userRepos = writable<AccessRepo[]>([]);
export const userLoading = writable(false);
export const userSearched = writable(false);

export async function findUser(owner: string, username: string): Promise<void> {
  if (!username.trim()) return;
  userRepos.set([]);
  userSearched.set(false);
  userLoading.set(true);
  try {
    userRepos.set(await fetchUserAccess(owner, username.trim()));
  } catch {
    userRepos.set([]);
  } finally {
    userLoading.set(false);
    userSearched.set(true);
  }
}