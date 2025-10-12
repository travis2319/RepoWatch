<script lang="ts">
	import { onMount } from 'svelte';
	import { githubUser, ownerName } from '$lib/store';

	// --- Type Definitions for Stronger Safety ---
	type CollaboratorRole = 'admin' | 'write' | 'read';

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
	let loading = $state(false);
	let repositories: Repository[] = $state([]);
	let hasSearched = $state(false);

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
		const savedUser = localStorage.getItem('githubUser');
		const savedOwner = localStorage.getItem('ownerName');
		if (savedUser) githubUser.set(savedUser);
		if (savedOwner) ownerName.set(savedOwner);
	});

	$effect(() => {
		if (typeof window !== 'undefined') {
			githubUser.set(localStorage.getItem('githubUser') || '');
			ownerName.set(localStorage.getItem('ownerName') || '');
		}
	});

	async function checkAccess() {
		console.log('checkAccess called with:', { githubUser: $githubUser, ownerName: $ownerName });

		if (!$githubUser || !$ownerName) {
			console.warn('Please configure both GitHub user and repository owner in the sidebar');
			return;
		}

		loading = true;
		hasSearched = false;
		console.log('Loading state set to true');

		try {
			console.log('Making API request...');
			const response = await fetch('http://localhost:4000/api/v1/check', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					owner: $ownerName,
					user: $githubUser
				})
			});

			console.log('Response received:', response.status, response.statusText);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const result = await response.json();
			console.log('API response:', result);

			const accessibleRepos = result.data.filter((repo: { repo: any; hasAccess: boolean }) => repo.hasAccess);

			repositories = accessibleRepos.map(
				(repo: any, index: number): Repository => ({
					id: index + 1,
					name: repo.repo,
					description: repo.description || `Repository: ${repo.repo}`,
					private: false,
					collaboratorRole: 'write',
					lastUpdated: repo.checkedAt,
					language: repo.language || 'Unknown'
				})
			);

			console.log('Final repositories array:', repositories);
		} catch (error) {
			console.error('Error checking repository access:', error);
			repositories = [];
		} finally {
			hasSearched = true;
			loading = false;
		}
	}

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	function getRoleColor(role: CollaboratorRole): string {
		switch (role) {
			case 'admin':
				return 'bg-red-100 text-red-800';
			case 'write':
				return 'bg-green-100 text-green-800';
			case 'read':
				return 'bg-blue-100 text-blue-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	}

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
	<!-- Header -->
	<header class="mb-8">
		<h1 class="mb-2 text-3xl font-bold text-gray-900">Repository Access Checker</h1>
		<p class="text-gray-600">Check which repositories a GitHub user has access to as a collaborator.</p>
	</header>

	<!-- Configuration -->
	<section class="mb-6 rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
		<h2 class="mb-4 text-xl font-semibold text-gray-900">Current Configuration</h2>

		{#if $githubUser && $ownerName}
			<div class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-2">
				<!-- GitHub User -->
				<div class="flex items-center rounded-lg border border-gray-200 bg-gray-50 p-4">
					<div class="flex-shrink-0">
						<svg class="h-8 w-8 text-gray-700" fill="currentColor" viewBox="0 0 24 24">
							<path
								d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
							/>
						</svg>
					</div>
					<div class="ml-4">
						<p class="text-sm font-medium text-gray-600">GitHub User</p>
						<p class="text-lg font-semibold text-gray-800">{$githubUser}</p>
					</div>
				</div>

				<!-- Repository Owner -->
				<div class="flex items-center rounded-lg border border-gray-200 bg-gray-50 p-4">
					<div class="flex-shrink-0">
						<svg class="h-8 w-8 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
						<p class="text-lg font-semibold text-gray-800">{$ownerName}</p>
					</div>
				</div>
			</div>

			<!-- Action Button -->
			<button
				onclick={checkAccess}
				disabled={loading}
				class="inline-flex w-full items-center justify-center rounded-md bg-blue-600 px-6 py-3 text-base font-medium text-white hover:bg-blue-700 disabled:opacity-50 sm:w-auto"
			>
				{#if loading}
					<svg class="mr-3 -ml-1 h-5 w-5 animate-spin text-white" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
						<path
							class="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373..."
						/>
					</svg>
					Checking Access...
				{:else}
					<svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
					</svg>
					Check Repository Access
				{/if}
			</button>
		{:else}
			<div class="py-8 text-center text-gray-500">
				<h3 class="text-lg font-medium text-gray-900">Configuration Required</h3>
				<p>Please set the GitHub user and repository owner in the sidebar to begin.</p>
			</div>
		{/if}
	</section>

	<!-- Results -->
	{#if hasSearched}
		<section class="rounded-lg border border-gray-200 bg-white shadow-sm">
			<header class="border-b border-gray-200 px-6 py-4">
				<h2 class="text-xl font-semibold text-gray-900">
					Accessible Repositories
					<span class="ml-2 rounded-full bg-blue-100 px-2.5 py-0.5 text-xs text-blue-800">
						{repositories.length} found
					</span>
				</h2>
				<p class="mt-1 text-sm text-gray-500">
					Showing repositories where <strong>{$githubUser}</strong> has collaborator access in
					<strong>{$ownerName}</strong>'s account.
				</p>
			</header>

			{#if repositories.length > 0}
				<div class="divide-y divide-gray-200">
					{#each repositories as repo (repo.id)}
						<article class="px-6 py-5 hover:bg-gray-50 transition-colors">
							<div class="flex flex-wrap justify-between gap-4">
								<div class="min-w-0 flex-1">
									<h3 class="truncate text-lg font-semibold text-blue-700">{repo.name}</h3>
									<p class="mb-3 text-sm text-gray-600">{repo.description}</p>
									<div class="flex flex-wrap gap-x-4 text-sm text-gray-500">
										<div class="flex items-center">
											<div class="h-3 w-3 rounded-full {getLanguageColor(repo.language)} mr-1.5"></div>
											{repo.language}
										</div>
										<div class="flex items-center">
											Updated {formatDate(repo.lastUpdated)}
										</div>
									</div>
								</div>
							</div>
						</article>
					{/each}
				</div>
			{:else}
				<div class="px-6 py-16 text-center text-gray-500">
					<h3 class="text-lg font-medium text-gray-900">No Accessible Repositories Found</h3>
					<p>
						The user <strong>{$githubUser}</strong> doesn’t appear to have collaborator access to any
						repositories owned by <strong>{$ownerName}</strong>.
					</p>
				</div>
			{/if}
		</section>
	{/if}
</div>
