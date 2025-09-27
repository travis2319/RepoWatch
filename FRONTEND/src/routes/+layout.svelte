<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head><script lang="ts">
  import '../app.css';
  import favicon from '$lib/assets/favicon.svg';
  import { onMount } from 'svelte';
  
  let { children } = $props();
  let sidebarOpen = $state(false);
  let githubUser = $state('');
  let ownerName = $state('');
  
  // Load saved values from localStorage on mount
  onMount(() => {
    const savedUser = localStorage.getItem('githubUser');
    const savedOwner = localStorage.getItem('ownerName');
    if (savedUser) githubUser = savedUser;
    if (savedOwner) ownerName = savedOwner;
  });
  
  // Save values to localStorage when they change
  $effect(() => {
    if (githubUser) localStorage.setItem('githubUser', githubUser);
  });
  
  $effect(() => {
    if (ownerName) localStorage.setItem('ownerName', ownerName);
  });
  
  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
  }
  
  function closeSidebar() {
    sidebarOpen = false;
  }
  
  function extractUsername(input: string): string {
    if (input.includes('github.com/')) {
      const match = input.match(/github\.com\/([^\/]+)/);
      return match ? match[1] : input;
    }
    return input;
  }
  
  function handleUserInput(event: Event) {
    const target = event.target as HTMLInputElement;
    githubUser = extractUsername(target.value);
  }
  
  function handleOwnerInput(event: Event) {
    const target = event.target as HTMLInputElement;
    ownerName = extractUsername(target.value);
  }
</script>

<!-- Sidebar overlay -->
{#if sidebarOpen}
  <div 
    class="fixed inset-0 bg-black bg-opacity-50 z-40 md:hidden" 
    onclick={closeSidebar}
    onkeydown={(e) => e.key === 'Escape' && closeSidebar()}
    role="button"
    tabindex="0"
  ></div>
{/if}

<!-- Sidebar -->
<div class="fixed inset-y-0 left-0 z-50 w-80 bg-gray-900 transform {sidebarOpen ? 'translate-x-0' : '-translate-x-full'} transition-transform duration-300 ease-in-out md:translate-x-0">
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex items-center justify-between p-4 border-b border-gray-700">
      <h1 class="text-xl font-semibold text-white">GitHub Repo Checker</h1>
      <button 
        onclick={closeSidebar}
        class="md:hidden text-gray-400 hover:text-white transition-colors"
        aria-label="Close sidebar"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
    
    <!-- Settings Form -->
    <div class="flex-1 p-4 space-y-6">
      <div>
        <label for="githubUser" class="block text-sm font-medium text-gray-300 mb-2">
          GitHub User
        </label>
        <input
          id="githubUser"
          type="text"
          placeholder="Enter username or GitHub URL"
          value={githubUser}
          oninput={handleUserInput}
          class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <p class="text-xs text-gray-400 mt-1">
          You can enter a username or paste a GitHub profile URL
        </p>
      </div>
      
      <div>
        <label for="ownerName" class="block text-sm font-medium text-gray-300 mb-2">
          Repository Owner
        </label>
        <input
          id="ownerName"
          type="text"
          placeholder="Enter owner username or GitHub URL"
          value={ownerName}
          oninput={handleOwnerInput}
          class="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <p class="text-xs text-gray-400 mt-1">
          The owner whose repositories you want to check access for
        </p>
      </div>
      
      {#if githubUser && ownerName}
        <div class="p-3 bg-green-900/30 border border-green-700 rounded-lg">
          <div class="flex items-center">
            <svg class="w-4 h-4 text-green-400 mr-2" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span class="text-sm text-green-300">Configuration saved</span>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- Main content area -->
<div class="md:pl-80 min-h-screen bg-gray-50">
  <!-- Top bar -->
  <div class="bg-white shadow-sm border-b border-gray-200 px-4 py-3 md:px-6">
    <div class="flex items-center">
      <button 
        onclick={toggleSidebar}
        class="md:hidden text-gray-500 hover:text-gray-700 transition-colors mr-4"
        aria-label="Open sidebar"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
      
      <div class="flex items-center space-x-2">
        <svg class="w-6 h-6 text-gray-900" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
        </svg>
        <h1 class="text-xl font-semibold text-gray-900">Repository Access Checker</h1>
      </div>
    </div>
  </div>
  
  <!-- Page content -->
  <main class="p-4 md:p-6">
    {@render children?.()}
  </main>
</div>