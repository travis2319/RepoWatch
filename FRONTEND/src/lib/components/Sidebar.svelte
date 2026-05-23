<script lang="ts">
  import GithubInput from './GithubInput.svelte';
  import { githubUser, ownerName } from '$lib/store';

  export let sidebarOpen: boolean;
  export let closeSidebar: () => void;
  export let saveSettings: () => void;
  export let saveMessage: string;
</script>

<div
  class="fixed inset-y-0 left-0 z-50 w-80 transform bg-gray-900
    {sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
    transition-transform duration-300 ease-in-out md:translate-x-0"
>
  <div class="flex h-full flex-col">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-gray-700 p-4">
      <h1 class="text-xl font-semibold text-white">REPO WATCH</h1>
      <button
        onclick={closeSidebar}
        class="text-gray-400 transition-colors hover:text-white md:hidden"
        aria-label="Close sidebar"
      >
        <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Settings Form -->
    <div class="flex-1 space-y-6 p-4">
      <GithubInput
        label="GitHub User"
        store={githubUser}
        placeholder="Enter username or GitHub URL"
        hint="You can enter a username or paste a GitHub profile URL"
      />
      <GithubInput
        label="Repository Owner"
        store={ownerName}
        placeholder="Enter owner username or GitHub URL"
        hint="The owner whose repositories you want to check access for"
      />

      <!-- Save button -->
      <button
        onclick={saveSettings}
        class="w-full rounded-lg bg-blue-600 px-4 py-2 font-medium text-white transition hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-400"
      >
        Save
      </button>

      <!-- Save confirmation -->
      {#if saveMessage}
        <div class="rounded-lg border border-green-700 bg-green-900/30 p-3">
          <div class="flex items-center">
            <svg class="mr-2 h-4 w-4 text-green-400" fill="currentColor" viewBox="0 0 20 20">
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clip-rule="evenodd"
              />
            </svg>
            <span class="text-sm text-green-300">{saveMessage}</span>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>