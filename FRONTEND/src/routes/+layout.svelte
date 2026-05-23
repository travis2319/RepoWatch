<script lang="ts">
  import '../app.css';
  import favicon from '$lib/assets/favicon.svg';
  import { onMount } from 'svelte';
  import { githubUser, ownerName } from '$lib/store';
  import Sidebar from '$lib/components/Sidebar.svelte';
  import TopBar from '$lib/components/TopBar.svelte';

  let { children } = $props();
  let sidebarOpen = $state(false);
  let saveMessage = $state('');

  onMount(() => {
    const savedUser = localStorage.getItem('githubUser');
    const savedOwner = localStorage.getItem('ownerName');
    if (savedUser) githubUser.set(savedUser);
    if (savedOwner) ownerName.set(savedOwner);
  });

  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
  }

  function closeSidebar() {
    sidebarOpen = false;
  }

  function saveToLocalStorage() {
    localStorage.setItem('githubUser', $githubUser);
    localStorage.setItem('ownerName', $ownerName);
    saveMessage = '✅ Settings saved successfully!';
    setTimeout(() => (saveMessage = ''), 3000);
  }
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<!-- Mobile sidebar overlay -->
{#if sidebarOpen}
  <div
    class="fixed inset-0 z-40 bg-black bg-opacity-50 md:hidden"
    onclick={closeSidebar}
    onkeydown={(e) => e.key === 'Escape' && closeSidebar()}
    role="button"
    tabindex="0"
  ></div>
{/if}

<!-- Sidebar -->
<Sidebar
  {sidebarOpen}
  {closeSidebar}
  saveSettings={saveToLocalStorage}
  {saveMessage}
/>

<!-- Main content area -->
<div class="min-h-screen bg-gray-50 md:pl-80">
  <TopBar {toggleSidebar} />

  <main class="p-4 md:p-6">
    {@render children?.()}
  </main>
</div>