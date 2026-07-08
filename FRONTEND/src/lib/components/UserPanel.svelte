<script lang="ts">
  import { userRepos, userLoading, userSearched, findUser } from '$lib/stores/user.store';
  import { LANG_COLORS } from '$lib/constants/index';

  const { onClose, ownerName }: { onClose: () => void; ownerName: string } = $props();

  let userSearch = $state('');

  function handleSearch() {
    findUser(ownerName, userSearch);
  }
</script>

<div class="panel">
  <div class="panel-header">
    <h2 class="panel-title">Find user in repos</h2>
    <button class="close-btn" onclick={onClose} aria-label="Close">✕</button>
  </div>

  <div class="panel-body">
    <div class="search-row">
      <input
        type="text"
        bind:value={userSearch}
        placeholder="GitHub username to check…"
        onkeydown={(e) => e.key === 'Enter' && handleSearch()}
      />
      <button onclick={handleSearch} disabled={$userLoading || !userSearch.trim()}>
        Search
        <svg style="width:16px;height:16px;vertical-align:-3px;margin-left:4px" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
      </button>
    </div>

    {#if $userLoading}
      <div class="loading-row"><span class="spinner"></span> Checking access…</div>
    {:else if $userSearched}
      {#if $userRepos.length === 0}
        <div class="empty">No accessible repos found for <strong>{userSearch}</strong></div>
      {:else}
        <div class="repo-list">
          {#each $userRepos as repo}
            <div class="repo-item">
              <div class="repo-left">
                <div class="repo-name">
                  <svg class="icon icon-sm" style="color:var(--clr-accent)" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5"/>
                  </svg>
                  {repo.repo}
                </div>
                {#if repo.language}
                  <div class="repo-meta">
                    <span class="lang-dot" style="background:{LANG_COLORS[repo.language] ?? '#888'}"></span>
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