import { writable } from 'svelte/store';
import { fetchValidate } from '$lib/services/github';
import type { ValidateData } from '$lib/types/github';

export const validate = writable<ValidateData | null>(null);
export const validateLoading = writable(false);
export const validateError = writable('');

export async function loadValidate(): Promise<void> {
  validate.set(null);
  validateError.set('');
  validateLoading.set(true);
  try {
    validate.set(await fetchValidate());
  } catch {
    validateError.set('Could not reach API');
  } finally {
    validateLoading.set(false);
  }
}