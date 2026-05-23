<script lang="ts">
  import { onMount } from 'svelte';
  import { ownerName } from '$lib/store';

  const API = '/api/v1';
  const langColors: Record<string, string> = {
    TypeScript: '#3178c6', JavaScript: '#f7df1e', Go: '#00add8',
    Python: '#3572a5', Dart: '#00add8', Java: '#b07219'
  };

  type Panel = 'repos' | 'user' | 'token' | null;
  let activePanel = $state<Panel>(null);

  interface Repo { name: string; private: boolean; language?: string }
  let repos = $state<Repo[]>([]);
  let reposLoading = $state(false);
  let reposError = $state('');

  let userSearch = $state('');
  interface AccessRepo { repo: string; hasAccess: boolean; language?: string }
  let userRepos = $state<AccessRepo[]>([]);
  let userLoading = $state(false);
  let userSearched = $state(false);

  interface TokenStatus {
    expiresAt?: string;
    daysRemaining?: number;
    scopes?: string[];
    rateLimit?: { remaining: number; limit: number };
  }
  let token = $state<TokenStatus | null>(null);
  let tokenLoading = $state(false);
  let tokenError = $state('');

  onMount(() => {
    const saved = localStorage.getItem('ownerName');
    if (saved) ownerName.set(saved);
  });

  function saveOwner() {
    localStorage.setItem('ownerName', $ownerName);
  }

  async function openPanel(panel: Panel) {
    activePanel = panel;
    if (panel === 'repos') await loadRepos();
    if (panel === 'token') await loadToken();
  }

  async function loadRepos() {
    repos = []; reposError = '';
    reposLoading = true;
    try {
      const r = await fetch(`${API}/repos?owner=${$ownerName}`);
      if (!r.ok) throw new Error(`${r.status}`);
      const d = await r.json();
      repos = d.data ?? d ?? [];
    } catch (e) {
      reposError = 'Could not reach API';
    } finally {
      reposLoading = false;
    }
  }

  async function findUser() {
    if (!userSearch.trim()) return;
    userRepos = []; userSearched = false;
    userLoading = true;
    try {
      const r = await fetch(`${API}/check`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ owner: $ownerName, user: userSearch.trim() })
      });
      if (!r.ok) throw new Error(`${r.status}`);
      const d = await r.json();
      userRepos = (d.data ?? []).filter((x: AccessRepo) => x.hasAccess);
    } catch (e) {
      userRepos = [];
    } finally {
      userLoading = false;
      userSearched = true;
    }
  }

  async function loadToken() {
    token = null; tokenError = '';
    tokenLoading = true;
    try {
      const r = await fetch(`${API}/token-status`);
      if (!r.ok) throw new Error(`${r.status}`);
      const d = await r.json();
      token = d.data ?? d;
    } catch (e) {
      tokenError = 'Could not fetch token info';
    } finally {
      tokenLoading = false;
    }
  }

  function formatDate(iso?: string) {
    if (!iso) return 'Never';
    return new Date(iso).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' });
  }

  function tokenBadge(days?: number): { label: string; cls: string; msg: string } {
    if (days === undefined) return { label: 'No expiry', cls: 'badge-success', msg: 'Token does not expire' };
    if (days <= 0) return { label: 'Expired', cls: 'badge-danger', msg: 'Token has expired and needs renewal' };
    if (days <= 7) return { label: 'Expiring soon', cls: 'badge-warn', msg: `Renew within ${days} day(s)` };
    return { label: 'Active', cls: 'badge-success', msg: 'Token is valid' };
  }
</script>

<svelte:head>
  <title>RepoWatch Dashboard</title>
</svelte:head>

<div class="dashboard">

  <div class="welcome">
    <span class="welcome-label">Welcome,</span>
    <div class="owner-wrap">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z"/></svg>
      <input
        class="owner-input"
        bind:value={$ownerName}
        onchange={saveOwner}
        placeholder="owner username"
        aria-label="Repository owner"
      />
      <svg class="icon icon-sm muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931z"/></svg>
    </div>
  </div>

  <div class="card-grid">
    <button class="card {activePanel === 'repos' ? 'active' : ''}" onclick={() => openPanel('repos')}>
      <svg class="card-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25"/></svg>
      <div class="card-title">List repositories</div>
      <div class="card-sub">Browse all public &amp; private repos</div>
    </button>

    <button class="card {activePanel === 'user' ? 'active' : ''}" onclick={() => openPanel('user')}>
      <svg class="card-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 15.75l-2.489-2.489m0 0a3.375 3.375 0 10-4.773-4.773 3.375 3.375 0 004.774 4.774zM21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
      <div class="card-title">Find user in repos</div>
      <div class="card-sub">Check collaborator access</div>
    </button>

    <button class="card {activePanel === 'token' ? 'active' : ''}" onclick={() => openPanel('token')}>
      <svg class="card-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 5.25a3 3 0 013 3m3 0a6 6 0 01-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1121.75 8.25z"/></svg>
      <div class="card-title">Token status</div>
      <div class="card-sub">Expiry &amp; scope details</div>
    </button>
  </div>

  {#if activePanel === 'repos'}
    <div class="panel">
      <div class="panel-header">
        <h2 class="panel-title">
          Repositories
          <span class="badge badge-info">{$ownerName}</span>
        </h2>
        <button class="close-btn" onclick={() => activePanel = null} aria-label="Close">✕</button>
      </div>
      <div class="panel-body">
        {#if reposLoading}
          <div class="loading-row"><span class="spinner"></span> Loading repos…</div>
        {:else if reposError}
          <div class="empty">{reposError}</div>
        {:else if repos.length === 0}
          <div class="empty">No repositories found</div>
        {:else}
          <div class="repo-list">
            {#each repos as repo}
              <div class="repo-item">
                <div class="repo-left">
                  <div class="repo-name">
                    {#if repo.private}
                      <svg class="icon icon-sm muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z"/></svg>
                    {:else}
                      <svg class="icon icon-sm muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5a17.92 17.92 0 01-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418"/></svg>
                    {/if}
                    {repo.name}
                  </div>
                  {#if repo.language}
                    <div class="repo-meta">
                      <span class="lang-dot" style="background:{langColors[repo.language] ?? '#888'}"></span>
                      {repo.language}
                    </div>
                  {/if}
                </div>
                <span class="badge {repo.private ? 'badge-warn' : 'badge-success'}">{repo.private ? 'private' : 'public'}</span>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  {/if}

  {#if activePanel === 'user'}
    <div class="panel">
      <div class="panel-header">
        <h2 class="panel-title">Find user in repos</h2>
        <button class="close-btn" onclick={() => activePanel = null} aria-label="Close">✕</button>
      </div>
      <div class="panel-body">
        <div class="search-row">
          <input
            type="text"
            bind:value={userSearch}
            placeholder="GitHub username to check…"
            onkeydown={(e) => e.key === 'Enter' && findUser()}
          />
          <button onclick={findUser} disabled={userLoading || !userSearch.trim()}>
            Search
            <svg style="width:16px;height:16px;vertical-align:-3px;margin-left:4px" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
          </button>
        </div>
        {#if userLoading}
          <div class="loading-row"><span class="spinner"></span> Checking access…</div>
        {:else if userSearched}
          {#if userRepos.length === 0}
            <div class="empty">No accessible repos found for <strong>{userSearch}</strong></div>
          {:else}
            <div class="repo-list">
              {#each userRepos as repo}
                <div class="repo-item">
                  <div class="repo-left">
                    <div class="repo-name">
                      <svg class="icon icon-sm" style="color:#16a34a" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5"/></svg>
                      {repo.repo}
                    </div>
                    {#if repo.language}
                      <div class="repo-meta">
                        <span class="lang-dot" style="background:{langColors[repo.language] ?? '#888'}"></span>
                        {repo.language}
                      </div>
                    {/if}
                  </div>
                  <span class="badge badge-success">has access</span>
                </div>
              {/each}
            </div>
          {/if}
        {/if}
      </div>
    </div>
  {/if}

  {#if activePanel === 'token'}
    <div class="panel">
      <div class="panel-header">
        <h2 class="panel-title">Token status</h2>
        <button class="close-btn" onclick={() => activePanel = null} aria-label="Close">✕</button>
      </div>
      <div class="panel-body">
        {#if tokenLoading}
          <div class="loading-row"><span class="spinner"></span> Checking token…</div>
        {:else if tokenError || !token}
          <div class="empty">{tokenError || 'No token data available'}</div>
        {:else}
          {@const tb = tokenBadge(token.daysRemaining)}
          <div class="token-status-row">
            <span class="badge {tb.cls}">{tb.label}</span>
            <span class="token-msg">{tb.msg}</span>
          </div>
          <div class="token-grid">
            <div class="token-stat">
              <div class="stat-label">Expires at</div>
              <div class="stat-value">{formatDate(token.expiresAt)}</div>
            </div>
            <div class="token-stat">
              <div class="stat-label">Days remaining</div>
              <div class="stat-value">{token.daysRemaining ?? '∞'}</div>
            </div>
            <div class="token-stat">
              <div class="stat-label">Scopes</div>
              <div class="stat-value scopes">{token.scopes?.join(', ') ?? '—'}</div>
            </div>
            <div class="token-stat">
              <div class="stat-label">Rate limit</div>
              <div class="stat-value">
                {#if token.rateLimit}
                  {token.rateLimit.remaining} / {token.rateLimit.limit}
                {:else}—{/if}
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}

</div>

<style>
  .dashboard { max-width: 860px; }
  .welcome { display: flex; align-items: center; gap: 12px; margin-bottom: 2rem; flex-wrap: wrap; }
  .welcome-label { font-size: 22px; font-weight: 500; }
  .owner-wrap { display: flex; align-items: center; gap: 8px; background: #f3f4f6; border: 1px solid #e5e7eb; border-radius: 10px; padding: 6px 12px; }
  .owner-input { background: transparent; border: none; outline: none; font-size: 22px; font-weight: 500; color: inherit; font-family: inherit; width: 200px; }
  .card-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px; margin-bottom: 2rem; }
  .card { background: white; border: 0.5px solid #e5e7eb; border-radius: 12px; padding: 1.5rem 1.25rem; cursor: pointer; text-align: left; transition: border-color .15s, box-shadow .15s; }
  .card:hover { border-color: #d1d5db; box-shadow: 0 1px 4px rgba(0,0,0,.06); }
  .card.active { border: 2px solid #3b82f6; }
  .card-icon { width: 32px; height: 32px; color: #6b7280; margin-bottom: .75rem; }
  .card-title { font-size: 15px; font-weight: 500; margin-bottom: 4px; }
  .card-sub { font-size: 13px; color: #6b7280; }
  .panel { background: white; border: 0.5px solid #e5e7eb; border-radius: 12px; overflow: hidden; }
  .panel-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-bottom: 0.5px solid #e5e7eb; }
  .panel-title { font-size: 16px; font-weight: 500; display: flex; align-items: center; gap: 10px; }
  .panel-body { padding: 1.25rem; }
  .close-btn { background: none; border: none; cursor: pointer; font-size: 16px; color: #6b7280; padding: 4px 8px; border-radius: 6px; }
  .close-btn:hover { background: #f3f4f6; }
  .search-row { display: flex; gap: 8px; margin-bottom: 1rem; }
  .search-row input { flex: 1; border: 1px solid #e5e7eb; border-radius: 8px; padding: 8px 12px; font-size: 14px; outline: none; font-family: inherit; }
  .search-row input:focus { border-color: #3b82f6; }
  .search-row button { background: #1d4ed8; color: white; border: none; border-radius: 8px; padding: 8px 16px; font-size: 14px; cursor: pointer; white-space: nowrap; font-family: inherit; }
  .search-row button:disabled { opacity: .5; cursor: not-allowed; }
  .repo-list { display: flex; flex-direction: column; gap: 8px; }
  .repo-item { display: flex; align-items: center; justify-content: space-between; padding: .75rem 1rem; background: #f9fafb; border-radius: 8px; }
  .repo-left { display: flex; flex-direction: column; gap: 3px; }
  .repo-name { display: flex; align-items: center; gap: 6px; font-size: 14px; font-weight: 500; }
  .repo-meta { display: flex; align-items: center; gap: 5px; font-size: 12px; color: #6b7280; }
  .lang-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
  .badge { display: inline-block; font-size: 11px; padding: 2px 8px; border-radius: 6px; font-weight: 500; }
  .badge-info { background: #dbeafe; color: #1e40af; }
  .badge-success { background: #dcfce7; color: #166534; }
  .badge-warn { background: #fef9c3; color: #854d0e; }
  .badge-danger { background: #fee2e2; color: #991b1b; }
  .token-status-row { display: flex; align-items: center; gap: 10px; margin-bottom: 1.25rem; }
  .token-msg { font-size: 14px; color: #6b7280; }
  .token-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .token-stat { background: #f9fafb; border-radius: 8px; padding: 1rem; }
  .stat-label { font-size: 12px; color: #6b7280; margin-bottom: 4px; }
  .stat-value { font-size: 18px; font-weight: 500; }
  .stat-value.scopes { font-size: 13px; line-height: 1.6; }
  .icon { width: 20px; height: 20px; flex-shrink: 0; }
  .icon-sm { width: 15px; height: 15px; }
  .muted { color: #9ca3af; }
  .loading-row { color: #6b7280; font-size: 14px; display: flex; align-items: center; gap: 8px; padding: 1rem 0; }
  .empty { text-align: center; padding: 2rem; color: #9ca3af; font-size: 14px; }
  .spinner { display: inline-block; width: 16px; height: 16px; border: 2px solid #e5e7eb; border-top-color: #3b82f6; border-radius: 50%; animation: spin .7s linear infinite; flex-shrink: 0; }
  @keyframes spin { to { transform: rotate(360deg); } }
</style>