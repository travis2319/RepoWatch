export function formatExpiry(iso?: string): string {
  if (!iso) return 'Never';
  return new Date(iso).toLocaleDateString('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function daysUntil(iso?: string): number | null {
  if (!iso) return null;
  const diff = new Date(iso).getTime() - Date.now();
  return Math.ceil(diff / (1000 * 60 * 60 * 24));
}

export function expiryBadge(days: number | null): { label: string; cls: string } {
  if (days === null) return { label: 'No expiry', cls: 'badge-success' };
  if (days <= 0)     return { label: 'Expired',      cls: 'badge-danger' };
  if (days <= 7)     return { label: `${days}d left`, cls: 'badge-warn' };
  return              { label: `${days}d left`,        cls: 'badge-success' };
}

export function resetTime(unixStr: string): string {
  return new Date(parseInt(unixStr) * 1000).toLocaleTimeString('en-GB', {
    hour: '2-digit',
    minute: '2-digit',
  });
}