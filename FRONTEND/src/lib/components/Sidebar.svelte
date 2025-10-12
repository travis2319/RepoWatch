<script lang="ts">
  import GithubInput from './GithubInput.svelte';
  import { githubUser, ownerName } from '$lib/store';
  export let sidebarOpen: boolean;
  export let closeSidebar: () => void;
  export let saveSettings: () => void;
  export let saveMessage: string;
</script>

<div class="fixed inset-y-0 left-0 z-50 w-80 transform bg-gray-900
            {sidebarOpen ? 'translate-x-0' : '-translate-x-full'} transition-transform duration-300 ease-in-out md:translate-x-0">
  <div class="flex h-full flex-col">
    <div class="flex items-center justify-between border-b border-gray-700 p-4">
      <h1 class="text-xl font-semibold text-white">GitHub Repo Checker</h1>
      <button on:click={closeSidebar} class="text-gray-400 md:hidden">
        ✕
      </button>
    </div>

    <div class="flex-1 space-y-6 p-4">
      <GithubInput label="GitHub User" store={githubUser} placeholder="Enter username or URL" />
      <GithubInput label="Repository Owner" store={ownerName} placeholder="Enter owner username or URL" />

      <button on:click={saveSettings} class="w-full rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700">
        Save
      </button>

      {#if saveMessage}
        <div class="rounded-lg border border-green-700 bg-green-900/30 p-3">
          <span class="text-green-300">{saveMessage}</span>
        </div>
      {/if}
    </div>
  </div>
</div>
