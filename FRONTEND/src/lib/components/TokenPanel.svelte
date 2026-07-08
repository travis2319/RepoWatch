<script lang="ts">
  import { validate, validateLoading, validateError } from '$lib/stores/token.store';
  import { formatExpiry, daysUntil, expiryBadge, resetTime } from '$lib/utils/format';

  const { onClose }: { onClose: () => void } = $props();
</script>

<div class="panel">
  <div class="panel-header">
    <h2 class="panel-title">
      Token status
      {#if $validate?.valid}
        <span class="badge badge-success">valid</span>
      {/if}
    </h2>
    <button class="close-btn" onclick={onClose} aria-label="Close">✕</button>
  </div>

  <div class="panel-body">
    {#if $validateLoading}
      <div class="loading-row"><span class="spinner"></span> Validating token…</div>
    {:else if $validateError || !$validate}
      <div class="empty">{$validateError || 'No data available'}</div>
    {:else}
      {@const v = $validate}
      {@const days = daysUntil(v.token_expiry)}
      {@const eb = expiryBadge(days)}
      {@const rlLimit = parseInt(v.rate_limit.limit)}
      {@const rlUsed = parseInt(v.rate_limit.used)}
      {@const rlRemaining = parseInt(v.rate_limit.remaining)}
      {@const rlPct = Math.round((rlRemaining / rlLimit) * 100)}

      <div class="user-identity">
        <img class="avatar" src={v.user.avatar_url} alt={v.user.login} />
        <div class="user-info">
          <div class="user-name">{v.user.name}</div>
          <a class="user-login" href={v.user.html_url} target="_blank" rel="noopener">
            @{v.user.login}
            <svg style="width:12px;height:12px;vertical-align:-1px;margin-left:3px" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25"/>
            </svg>
          </a>
        </div>
        <div class="user-id">ID #{v.user.id}</div>
      </div>

      <div class="token-grid">
        <div class="token-stat">
          <div class="stat-label">Token expires</div>
          <div class="stat-value-row">
            <div class="stat-value">{formatExpiry(v.token_expiry)}</div>
            <span class="badge {eb.cls}">{eb.label}</span>
          </div>
        </div>

        <div class="token-stat">
          <div class="stat-label">API rate limit — resets {resetTime(v.rate_limit.reset)}</div>
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
  </div>
</div>