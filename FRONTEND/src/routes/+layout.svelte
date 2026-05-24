<script lang="ts">
  import '../app.css';
  import favicon from '$lib/assets/favicon.svg';
  import { onMount } from 'svelte';
  import { githubUser, ownerName } from '$lib/store';
  import TopBar from '$lib/components/TopBar.svelte';

  let { children } = $props();

  onMount(() => {
    const savedUser = localStorage.getItem('githubUser');
    const savedOwner = localStorage.getItem('ownerName');
    if (savedUser) githubUser.set(savedUser);
    if (savedOwner) ownerName.set(savedOwner);
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-base">
  <TopBar />

  <main class="p-4 md:p-8 max-w-5xl mx-auto">
    {@render children?.()}
  </main>
</div>