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
		console.log('checkAccess called with:', { githubUser, ownerName });

		if (!githubUser || !ownerName) {
			console.warn('Please configure both GitHub user and repository owner in the sidebar');
			return;
		}

		loading = true;
		hasSearched = false;
		console.log('Loading state set to true');

		try {
			console.log('Making API request...');
			const response = await fetch('http://localhost:3000/api/v1/check', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					owner: ownerName,
					user: githubUser
				})
			});

			console.log('Response received:', response.status, response.statusText);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const result = await response.json();
			console.log('API response:', result);

			// Filter repositories where user has access
			const accessibleRepos = result.data.filter((repo: { repo: any; hasAccess: boolean }) => {
				console.log(`Checking repo ${repo.repo}: hasAccess = ${repo.hasAccess}`);
				return repo.hasAccess === true;
			});

			console.log('Accessible repos after filtering:', accessibleRepos);

			// Transform the API response to match the Repository interface exactly
			repositories = accessibleRepos.map(
				(repo: any, index: number): Repository => ({
					id: index + 1, // Generate unique ID for each repository
					name: repo.repo,
					description: repo.description || `Repository: ${repo.repo}`, // Fallback description
					private: false, // Default to false, you might want to fetch this from GitHub API
					collaboratorRole: 'write' as CollaboratorRole, // Default role, you might want to determine this from API
					lastUpdated: repo.checkedAt, // Use checkedAt as lastUpdated
					language: repo.language || 'Unknown' // Default language
				})
			);

			console.log('Final repositories array:', repositories);
			console.log(
				`Found ${repositories.length} accessible repositories out of ${result.data.length} total repositories`
			);
		} catch (error) {
			console.error('Error checking repository access:', error);
			repositories = [];
			// You might want to show an error message to the user here
		} finally {
			console.log('Setting hasSearched to true and loading to false');
			hasSearched = true;
			loading = false;
		}
	}
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
			case 'admin':
				return 'bg-red-100 text-red-800';
			case 'write':
				return 'bg-green-100 text-green-800';
			case 'read':
				return 'bg-blue-100 text-blue-800';
			// A default case is good practice, though with strict types, it's less likely to be hit.
			default:
				return 'bg-gray-100 text-gray-800';
		}
	}

	/**
	 * Returns a Tailwind CSS background color for a given programming language.
	 * @param language The programming language name.
	 * @returns A string containing a background color class.
	 */
	function getLanguageColor(language: string): string {
		const colors: Record<string, string> = {
			TypeScript: 'bg-blue-500',
			JavaScript: 'bg-yellow-400',
			Python: 'bg-green-600',
			Dart: 'bg-blue-400',
			Java: 'bg-orange-500',
			Go: 'bg-cyan-500'
		};
		return colors[language] || 'bg-gray-500';
	}
</script>

<svelte:head>
	<title>GitHub Repository Access Checker</title>
</svelte:head>

<div class="mx-auto max-w-6xl p-4 md:p-6">
	<!-- Header section -->
	<header class="mb-8">
		<h1 class="mb-2 text-3xl font-bold text-gray-900">Repository Access Checker</h1>
		<p class="text-gray-600">
			Check which repositories a GitHub user has access to as a collaborator.
		</p>
	</header>

	<!-- Configuration status -->
	<section class="mb-6 rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
		<h2 class="mb-4 text-xl font-semibold text-gray-900">Current Configuration</h2>

		{#if githubUser && ownerName}
			<div class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-2">
				<!-- GitHub User Info -->
				<div class="flex items-center rounded-lg border border-gray-200 bg-gray-50 p-4">
					<div class="flex-shrink-0">
						<svg
							class="h-8 w-8 text-gray-700"
							fill="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
							/>
						</svg>
					</div>
					<div class="ml-4">
						<p class="text-sm font-medium text-gray-600">GitHub User</p>
						<p class="text-lg font-semibold text-gray-800">{githubUser}</p>
					</div>
				</div>

				<!-- Repository Owner Info -->
				<div class="flex items-center rounded-lg border border-gray-200 bg-gray-50 p-4">
					<div class="flex-shrink-0">
						<svg
							class="h-8 w-8 text-gray-700"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
							aria-hidden="true"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
							/>
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
				class="inline-flex w-full items-center justify-center rounded-md border border-transparent bg-blue-600 px-6 py-3 text-base font-medium text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 sm:w-auto"
			>
				{#if loading}
					<svg
						class="mr-3 -ml-1 h-5 w-5 animate-spin text-white"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
					>
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
						></circle>
						<path
							class="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						></path>
					</svg>
					Checking Access...
				{:else}
					<svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
						/>
					</svg>
					Check Repository Access
				{/if}
			</button>
		{:else}
			<!-- Prompt to Configure -->
			<div class="py-8 text-center">
				<svg
					class="mx-auto h-12 w-12 text-gray-400"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					aria-hidden="true"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
					/>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
					/>
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
		<section class="rounded-lg border border-gray-200 bg-white shadow-sm">
			<header class="border-b border-gray-200 px-6 py-4">
				<h2 class="text-xl font-semibold text-gray-900">
					Accessible Repositories
					<span
						class="ml-2 inline-flex items-center rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-medium text-blue-800"
					>
						{repositories.length} found
					</span>
				</h2>
				<p class="mt-1 text-sm text-gray-500">
					Showing repositories where <strong>{githubUser}</strong> has collaborator access in
					<strong>{ownerName}</strong>'s account.
				</p>
			</header>

			{#if repositories.length > 0}
				<div class="divide-y divide-gray-200">
					{#each repositories as repo (repo.id)}
						<article class="px-6 py-5 transition-colors hover:bg-gray-50">
							<div class="flex flex-wrap items-start justify-between gap-4">
								<div class="min-w-0 flex-1">
									<div class="mb-2 flex flex-wrap items-center gap-x-3">
										<h3 class="truncate text-lg font-semibold text-blue-700">{repo.name}</h3>
										{#if repo.private}
											<span
												class="inline-flex items-center rounded border border-gray-300 bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-800"
											>
												<svg
													class="mr-1.5 h-3 w-3"
													fill="currentColor"
													viewBox="0 0 20 20"
													aria-hidden="true"
												>
													<path
														fill-rule="evenodd"
														d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
														clip-rule="evenodd"
													/>
												</svg>
												Private
											</span>
										{/if}
										<span
											class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium {getRoleColor(
												repo.collaboratorRole
											)}"
										>
											{repo.collaboratorRole}
										</span>
									</div>

									<p class="mb-3 text-sm text-gray-600">{repo.description}</p>

									<div class="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-gray-500">
										<div class="flex items-center">
											<div
												class="h-3 w-3 rounded-full {getLanguageColor(repo.language)} mr-1.5"
											></div>
											{repo.language}
										</div>
										<div class="flex items-center">
											<svg
												class="mr-1.5 h-4 w-4"
												fill="none"
												stroke="currentColor"
												viewBox="0 0 24 24"
												aria-hidden="true"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
												/>
											</svg>
											Updated {formatDate(repo.lastUpdated)}
										</div>
									</div>
								</div>

								<!-- svelte-ignore a11y_invalid_attribute -->
								<a
									href="#"
									class="ml-4 inline-flex flex-shrink-0 items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm transition-colors hover:bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none"
								>
									<svg
										class="mr-2 h-4 w-4"
										fill="none"
										stroke="currentColor"
										viewBox="0 0 24 24"
										aria-hidden="true"
									>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
										/>
									</svg>
									View on GitHub
								</a>
							</div>
						</article>
					{/each}
				</div>
			{:else}
				<!-- No Results Found -->
				<div class="px-6 py-16 text-center">
					<svg
						class="mx-auto h-12 w-12 text-gray-400"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						aria-hidden="true"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
						/>
					</svg>
					<h3 class="mt-4 text-lg font-medium text-gray-900">No Accessible Repositories Found</h3>
					<p class="mt-2 text-sm text-gray-500">
						The user <strong>{githubUser}</strong> doesn't appear to have collaborator access to any
						repositories owned by <strong>{ownerName}</strong>.
					</p>
				</div>
			{/if}
		</section>
	{/if}
</div>
