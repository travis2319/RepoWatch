<script lang="ts">
  import { onMount } from 'svelte';

  // --- Type Definitions for Stronger Safety ---

  /**
   * Defines the possible roles a collaborator can have.
   * Using a literal type prevents typos and ensures only valid roles are used.
   */
  type CollaboratorRole = 'admin' | 'write' | 'read';

  /**
   * Defines the shape of a GitHub repository object.
   * This ensures that every repository object in our app is consistent.
   */
  interface Repository {
    id: number;
    name: string;
    description: string;
    private: boolean;
    collaboratorRole: CollaboratorRole;
    lastUpdated: string; // ISO date string
    language: string;
  }

  // --- Component State ---

  let githubUser = $state('');
  let ownerName = $state('');
  let loading = $state(false);
  // Explicitly type the repositories state as an array of Repository objects.
  let repositories: Repository[] = $state([]);
  let hasSearched = $state(false);

  // Mock repository data now strictly follows the Repository interface.
  const mockRepositories: Repository[] = [
    {
      id: 1,
      name: 'awesome-project',
      description: 'An awesome project with great features',
      private: false,
      collaboratorRole: 'write',
      lastUpdated: '2024-12-15T10:30:00Z',
      language: 'TypeScript'
    },
    {
      id: 2,
      name: 'secret-repo',
      description: 'A private repository for internal use',
      private: true,
      collaboratorRole: 'read',
      lastUpdated: '2024-12-14T15:45:00Z',
      language: 'Python'
    },
    {
      id: 3,
      name: 'web-app',
      description: 'Modern web application built with Svelte',
      private: false,
      collaboratorRole: 'admin',
      lastUpdated: '2024-12-16T09:20:00Z',
      language: 'JavaScript'
    },
    {
      id: 4,
      name: 'mobile-client',
      description: 'Cross-platform mobile application',
      private: true,
      collaboratorRole: 'write',
      lastUpdated: '2024-12-13T14:10:00Z',
      language: 'Dart'
    }
  ];

  onMount(() => {
    // Load saved values from localStorage
    const savedUser = localStorage.getItem('githubUser');
    const savedOwner = localStorage.getItem('ownerName');
    if (savedUser) githubUser = savedUser;
    if (savedOwner) ownerName = savedOwner;
  });

  // This effect will react to changes in localStorage if needed,
  // but onMount is sufficient for initial load.
  $effect(() => {
    if (typeof window !== 'undefined') {
      githubUser = localStorage.getItem('githubUser') || '';
      ownerName = localStorage.getItem('ownerName') || '';
    }
  });

  async function checkAccess() {
    if (!githubUser || !ownerName) {
      // Switched from alert to a more modern notification approach (if this were a real app)
      console.warn('Please configure both GitHub user and repository owner in the sidebar');
      return;
    }

    loading = true;
    hasSearched = false;

    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 1500));

    // In a real app, you would fetch data and validate it against the Repository interface.
    repositories = mockRepositories;
    hasSearched = true;
    loading = false;
  }

  // --- Utility Functions ---

  /**
   * Formats an ISO date string into a more readable format.
   * @param dateString The ISO date string to format.
   * @returns A formatted date string (e.g., "Dec 15, 2024").
   */
  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }

  /**
   * Returns Tailwind CSS classes based on the collaborator's role.
   * The parameter is now typed as CollaboratorRole for better safety.
   * @param role The collaborator role ('admin', 'write', or 'read').
   * @returns A string of CSS classes for styling.
   */
  function getRoleColor(role: CollaboratorRole): string {
    switch (role) {
      case 'admin': return 'bg-red-100 text-red-800';
      case 'write': return 'bg-green-100 text-green-800';
      case 'read': return 'bg-blue-100 text-blue-800';
      // A default case is good practice, though with strict types, it's less likely to be hit.
      default: return 'bg-gray-100 text-gray-800';
    }
  }
  
  /**
   * Returns a Tailwind CSS background color for a given programming language.
   * @param language The programming language name.
   * @returns A string containing a background color class.
   */
  function getLanguageColor(language: string): string {
    const colors: Record<string, string> = {
      'TypeScript': 'bg-blue-500',
      'JavaScript': 'bg-yellow-400',
      'Python': 'bg-green-600',
      'Dart': 'bg-blue-400',
      'Java': 'bg-orange-500',
      'Go': 'bg-cyan-500'
    };
    return colors[language] || 'bg-gray-500';
  }
</script>

<svelte:head>
  <title>GitHub Repository Access Checker</title>
</svelte:head>

<div class="max-w-6xl mx-auto p-4 md:p-6">
  <!-- Header section -->
  <header class="mb-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-2">Repository Access Checker</h1>
    <p class="text-gray-600">
      Check which repositories a GitHub user has access to as a collaborator.
    </p>
  </header>
  
  <!-- Configuration status -->
  <section class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
    <h2 class="text-xl font-semibold text-gray-900 mb-4">Current Configuration</h2>
    
    {#if githubUser && ownerName}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <!-- GitHub User Info -->
        <div class="flex items-center p-4 bg-gray-50 rounded-lg border border-gray-200">
          <div class="flex-shrink-0">
            <svg class="w-8 h-8 text-gray-700" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">GitHub User</p>
            <p class="text-lg font-semibold text-gray-800">{githubUser}</p>
          </div>
        </div>
        
        <!-- Repository Owner Info -->
        <div class="flex items-center p-4 bg-gray-50 rounded-lg border border-gray-200">
          <div class="flex-shrink-0">
            <svg class="w-8 h-8 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Repository Owner</p>
            <p class="text-lg font-semibold text-gray-800">{ownerName}</p>
          </div>
        </div>
      </div>
      
      <!-- Action Button -->
      <button 
        onclick={checkAccess}
        disabled={loading}
        class="w-full sm:w-auto inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {#if loading}
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Checking Access...
        {:else}
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          Check Repository Access
        {/if}
      </button>
    {:else}
      <!-- Prompt to Configure -->
      <div class="text-center py-8">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900">Configuration Required</h3>
        <p class="mt-2 text-sm text-gray-500">
          Please set the GitHub user and repository owner in the sidebar to begin.
        </p>
      </div>
    {/if}
  </section>
  
  <!-- Results section -->
  {#if hasSearched}
    <section class="bg-white rounded-lg shadow-sm border border-gray-200">
      <header class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-xl font-semibold text-gray-900">
          Accessible Repositories 
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 ml-2">
            {repositories.length} found
          </span>
        </h2>
        <p class="text-sm text-gray-500 mt-1">
          Showing repositories where <strong>{githubUser}</strong> has collaborator access in <strong>{ownerName}</strong>'s account.
        </p>
      </header>
      
      {#if repositories.length > 0}
        <div class="divide-y divide-gray-200">
          {#each repositories as repo (repo.id)}
            <article class="px-6 py-5 hover:bg-gray-50 transition-colors">
              <div class="flex items-start justify-between flex-wrap gap-4">
                <div class="flex-1 min-w-0">
                  <div class="flex items-center flex-wrap gap-x-3 mb-2">
                    <h3 class="text-lg font-semibold text-blue-700 truncate">{repo.name}</h3>
                    {#if repo.private}
                      <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800 border border-gray-300">
                        <svg class="w-3 h-3 mr-1.5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                          <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
                        </svg>
                        Private
                      </span>
                    {/if}
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getRoleColor(repo.collaboratorRole)}">
                      {repo.collaboratorRole}
                    </span>
                  </div>
                  
                  <p class="text-gray-600 mb-3 text-sm">{repo.description}</p>
                  
                  <div class="flex items-center flex-wrap gap-x-4 gap-y-2 text-sm text-gray-500">
                    <div class="flex items-center">
                      <div class="w-3 h-3 rounded-full {getLanguageColor(repo.language)} mr-1.5"></div>
                      {repo.language}
                    </div>
                    <div class="flex items-center">
                      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      Updated {formatDate(repo.lastUpdated)}
                    </div>
                  </div>
                </div>
                
                <!-- svelte-ignore a11y_invalid_attribute -->
                <a href="#" class="ml-4 flex-shrink-0 inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors">
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                  View on GitHub
                </a>
              </div>
            </article>
          {/each}
        </div>
      {:else}
        <!-- No Results Found -->
        <div class="text-center py-16 px-6">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h3 class="mt-4 text-lg font-medium text-gray-900">No Accessible Repositories Found</h3>
          <p class="mt-2 text-sm text-gray-500">
            The user <strong>{githubUser}</strong> doesn't appear to have collaborator access to any repositories owned by <strong>{ownerName}</strong>.
          </p>
        </div>
      {/if}
    </section>
  {/if}
</div>
