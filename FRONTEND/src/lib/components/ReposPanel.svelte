<script lang="ts">
  import { repos, reposLoading, reposError } from '$lib/stores/repos.store';
  import { ownerName } from '$lib/store';
  import { LANG_COLORS } from '$lib/constants';

  export let onClose: () => void;
</script>

<div class="panel">
  <div class="panel-header">
    <h2 class="panel-title">
      Repositories
      <span class="badge badge-info">{$ownerName}</span>
    </h2>
    <button class="close-btn" onclick={onClose} aria-label="Close">✕</button>
  </div>

  <div class="panel-body">
    {#if $reposLoading}
      <div class="loading-row"><span class="spinner"></span> Loading repos…</div>
    {:else if $reposError}
      <div class="empty">{$reposError}</div>
    {:else if $repos.length === 0}
      <div class="empty">No repositories found</div>
    {:else}
      <div class="repo-list">
        {#each $repos as repo}
          <div class="repo-item">
            <div class="repo-left">
              <div class="repo-name">
                {#if repo.private}
                  <svg class="icon icon-sm muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z"/>
                  </svg>
                {:else}
                  <svg class="icon icon-sm muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5a17.92 17.92 0 01-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418"/>
                  </svg>
                {/if}
                {repo.name}
              </div>
              {#if repo.language}
                <div class="repo-meta">
                  <span class="lang-dot" style="background:{LANG_COLORS[repo.language] ?? '#888'}"></span>
                  {repo.language}
                </div>
              {/if}
            </div>
            <span class="badge {repo.private ? 'badge-warn' : 'badge-success'}">
              {repo.private ? 'private' : 'public'}
            </span>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>