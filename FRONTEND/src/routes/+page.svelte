<script lang="ts">
  import { onMount } from 'svelte';
  import { ownerName } from '$lib/store';

  const API = '/api/v1';
  const langColors: Record<string, string> = {
    TypeScript: '#3178c6', JavaScript: '#f7df1e', Go: '#00d4aa',
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

  interface GithubUser {
    login: string;
    name: string;
    avatar_url: string;
    html_url: string;
    id: number;
  }
  interface RateLimit {
    limit: string;
    remaining: string;
    used: string;
    reset: string;
    resource: string;
  }
  interface ValidateData {
    valid: boolean;
    user: GithubUser;
    rate_limit: RateLimit;
    token_expiry: string;
  }
  let validate = $state<ValidateData | null>(null);
  let validateLoading = $state(false);
  let validateError = $state('');

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
    if (panel === 'token') await loadValidate();
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

  async function loadValidate() {
    validate = null; validateError = '';
    validateLoading = true;
    try {
      const r = await fetch(`${API}/github/validate`);
      if (!r.ok) throw new Error(`${r.status}`);
      const d = await r.json();
      validate = d.data ?? d;
    } catch (e) {
      validateError = 'Could not reach API';
    } finally {
      validateLoading = false;
    }
  }

  function formatExpiry(iso?: string) {
    if (!iso) return 'Never';
    return new Date(iso).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' });
  }

  function daysUntil(iso?: string): number | null {
    if (!iso) return null;
    const diff = new Date(iso).getTime() - Date.now();
    return Math.ceil(diff / (1000 * 60 * 60 * 24));
  }

  function expiryBadge(days: number | null): { label: string; cls: string } {
    if (days === null) return { label: 'No expiry', cls: 'badge-success' };
    if (days <= 0)    return { label: 'Expired',       cls: 'badge-danger' };
    if (days <= 7)    return { label: `${days}d left`,  cls: 'badge-warn' };
    return { label: `${days}d left`, cls: 'badge-success' };
  }

  function resetTime(unixStr: string) {
    return new Date(parseInt(unixStr) * 1000).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' });
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
      <div class="card-sub">Validate &amp; inspect GitHub token</div>
    </button>
  </div>

  <!-- Repos panel -->
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

  <!-- User panel -->
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
                      <svg class="icon icon-sm" style="color:var(--clr-accent)" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5"/></svg>
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

  <!-- Token / validate panel -->
  {#if activePanel === 'token'}
    <div class="panel">
      <div class="panel-header">
        <h2 class="panel-title">
          Token status
          {#if validate?.valid}
            <span class="badge badge-success">valid</span>
          {/if}
        </h2>
        <button class="close-btn" onclick={() => activePanel = null} aria-label="Close">✕</button>
      </div>
      <div class="panel-body">
        {#if validateLoading}
          <div class="loading-row"><span class="spinner"></span> Validating token…</div>
        {:else if validateError || !validate}
          <div class="empty">{validateError || 'No data available'}</div>
        {:else}
          <!-- User identity row -->
          <div class="user-identity">
            <img class="avatar" src={validate.user.avatar_url} alt={validate.user.login} />
            <div class="user-info">
              <div class="user-name">{validate.user.name}</div>
              <a class="user-login" href={validate.user.html_url} target="_blank" rel="noopener">
                @{validate.user.login}
                <svg style="width:12px;height:12px;vertical-align:-1px;margin-left:3px" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25"/></svg>
              </a>
            </div>
            <div class="user-id">ID #{validate.user.id}</div>
          </div>

          {#if validate}
            {@const days = daysUntil(validate.token_expiry)}
            {@const eb = expiryBadge(days)}
            {@const rlUsed = parseInt(validate.rate_limit.used)}
            {@const rlLimit = parseInt(validate.rate_limit.limit)}
            {@const rlRemaining = parseInt(validate.rate_limit.remaining)}
            {@const rlPct = Math.round((rlRemaining / rlLimit) * 100)}
            <div class="token-grid">
              <!-- Token expiry -->
              <div class="token-stat">
                <div class="stat-label">Token expires</div>
                <div class="stat-value-row">
                  <div class="stat-value">{formatExpiry(validate.token_expiry)}</div>
                  <span class="badge {eb.cls}">{eb.label}</span>
                </div>
              </div>

              <!-- Rate limit -->
              <div class="token-stat">
                <div class="stat-label">API rate limit — resets {resetTime(validate.rate_limit.reset)}</div>
                <div class="stat-value-row">
                  <div class="stat-value">{rlRemaining} <span class="stat-of">/ {rlLimit}</span></div>
                  <span class="badge {rlPct > 50 ? 'badge-success' : rlPct > 20 ? 'badge-warn' : 'badge-danger'}">{rlPct}% left</span>
                </div>
                <div class="rate-bar-track">
                  <div
                    class="rate-bar-fill"
                    style="width:{rlPct}%; background:{rlPct > 50 ? 'var(--clr-accent)' : rlPct > 20 ? 'var(--clr-warn)' : 'var(--clr-danger)'}"
                  ></div>
                </div>
                <div class="rate-sub">{rlUsed} requests used this window</div>
              </div>
            </div>
          {/if}
        {/if}
      </div>
    </div>
  {/if}

</div>

<style>
  .dashboard { max-width: 860px; }

  .welcome { display: flex; align-items: center; gap: 12px; margin-bottom: 2.5rem; flex-wrap: wrap; }
  .welcome-label { font-size: 22px; font-weight: 500; color: var(--clr-muted); }
  .owner-wrap {
    display: flex; align-items: center; gap: 8px;
    background: var(--clr-surface);
    border: 1px solid var(--clr-border);
    border-radius: 10px; padding: 6px 12px;
  }
  .owner-wrap:focus-within { border-color: var(--clr-accent); }
  .owner-input {
    background: transparent; border: none; outline: none;
    font-size: 22px; font-weight: 500;
    color: var(--clr-accent);
    font-family: inherit; width: 200px;
  }
  .owner-input::placeholder { color: var(--clr-muted); }

  .card-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 14px; margin-bottom: 2rem; }
  .card {
    background: var(--clr-surface); border: 1px solid var(--clr-border);
    border-radius: 12px; padding: 1.5rem 1.25rem;
    cursor: pointer; text-align: left;
    transition: border-color .15s, box-shadow .15s;
  }
  .card:hover { border-color: var(--clr-accent); box-shadow: 0 0 0 1px color-mix(in srgb, var(--clr-accent) 20%, transparent); }
  .card.active { border: 2px solid var(--clr-accent); }
  .card-icon { width: 30px; height: 30px; color: var(--clr-accent); margin-bottom: .75rem; }
  .card-title { font-size: 14px; font-weight: 500; margin-bottom: 4px; color: var(--clr-text); }
  .card-sub { font-size: 12px; color: var(--clr-muted); }

  .panel { background: var(--clr-surface); border: 1px solid var(--clr-border); border-radius: 12px; overflow: hidden; }
  .panel-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 1rem 1.25rem; border-bottom: 1px solid var(--clr-border);
    background: var(--clr-surface2);
  }
  .panel-title { font-size: 15px; font-weight: 500; display: flex; align-items: center; gap: 10px; color: var(--clr-text); }
  .panel-body { padding: 1.25rem; }
  .close-btn { background: none; border: none; cursor: pointer; font-size: 16px; color: var(--clr-muted); padding: 4px 8px; border-radius: 6px; }
  .close-btn:hover { background: var(--clr-surface2); color: var(--clr-text); }

  .search-row { display: flex; gap: 8px; margin-bottom: 1rem; }
  .search-row input {
    flex: 1; background: var(--clr-bg); border: 1px solid var(--clr-border);
    border-radius: 8px; padding: 8px 12px; font-size: 14px; outline: none;
    font-family: inherit; color: var(--clr-text);
  }
  .search-row input::placeholder { color: var(--clr-muted); }
  .search-row input:focus { border-color: var(--clr-accent); }
  .search-row button {
    background: var(--clr-accent); color: #0d1117; border: none; border-radius: 8px;
    padding: 8px 16px; font-size: 14px; cursor: pointer; white-space: nowrap;
    font-family: inherit; font-weight: 500;
  }
  .search-row button:hover { background: var(--clr-accent-dk); }
  .search-row button:disabled { opacity: .4; cursor: not-allowed; }

  .repo-list { display: flex; flex-direction: column; gap: 6px; }
  .repo-item {
    display: flex; align-items: center; justify-content: space-between;
    padding: .7rem 1rem; background: var(--clr-surface2);
    border: 1px solid var(--clr-border); border-radius: 8px;
  }
  .repo-left { display: flex; flex-direction: column; gap: 3px; }
  .repo-name { display: flex; align-items: center; gap: 6px; font-size: 13px; font-weight: 500; color: var(--clr-text); }
  .repo-meta { display: flex; align-items: center; gap: 5px; font-size: 12px; color: var(--clr-muted); }
  .lang-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }

  .badge { display: inline-block; font-size: 11px; padding: 2px 8px; border-radius: 5px; font-weight: 500; }
  .badge-info    { background: #0e2a4a; color: #60a5fa; border: 1px solid #1e3a5f; }
  .badge-success { background: #052e16; color: var(--clr-accent); border: 1px solid #064e32; }
  .badge-warn    { background: #2d1a00; color: var(--clr-warn);   border: 1px solid #4a2f00; }
  .badge-danger  { background: #2d0a0a; color: var(--clr-danger); border: 1px solid #4a1515; }

  /* User identity */
  .user-identity {
    display: flex; align-items: center; gap: 14px;
    padding: 1rem; margin-bottom: 1.25rem;
    background: var(--clr-surface2); border: 1px solid var(--clr-border);
    border-radius: 10px;
  }
  .avatar { width: 52px; height: 52px; border-radius: 50%; border: 2px solid var(--clr-border); flex-shrink: 0; }
  .user-info { flex: 1; display: flex; flex-direction: column; gap: 3px; }
  .user-name { font-size: 16px; font-weight: 500; color: var(--clr-text); }
  .user-login { font-size: 13px; color: var(--clr-accent); text-decoration: none; }
  .user-login:hover { text-decoration: underline; }
  .user-id { font-size: 11px; color: var(--clr-muted); white-space: nowrap; }

  /* Token grid */
  .token-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
  .token-stat { background: var(--clr-surface2); border: 1px solid var(--clr-border); border-radius: 8px; padding: 1rem; }
  .stat-label { font-size: 11px; color: var(--clr-muted); margin-bottom: 8px; text-transform: uppercase; letter-spacing: .05em; }
  .stat-value-row { display: flex; align-items: baseline; justify-content: space-between; gap: 8px; margin-bottom: 8px; }
  .stat-value { font-size: 15px; font-weight: 500; color: var(--clr-text); }
  .stat-of { font-size: 12px; color: var(--clr-muted); font-weight: 400; }

  /* Rate limit bar */
  .rate-bar-track { height: 4px; background: var(--clr-border); border-radius: 2px; overflow: hidden; margin-bottom: 6px; }
  .rate-bar-fill { height: 100%; border-radius: 2px; transition: width .4s ease; }
  .rate-sub { font-size: 11px; color: var(--clr-muted); }

  .icon { width: 20px; height: 20px; flex-shrink: 0; }
  .icon-sm { width: 15px; height: 15px; }
  .muted { color: var(--clr-muted); }
  .loading-row { color: var(--clr-muted); font-size: 14px; display: flex; align-items: center; gap: 8px; padding: 1rem 0; }
  .empty { text-align: center; padding: 2rem; color: var(--clr-muted); font-size: 13px; }
  .spinner {
    display: inline-block; width: 16px; height: 16px;
    border: 2px solid var(--clr-border); border-top-color: var(--clr-accent);
    border-radius: 50%; animation: spin .7s linear infinite; flex-shrink: 0;
  }
  @keyframes spin { to { transform: rotate(360deg); } }
</style>