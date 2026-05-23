<script lang="ts">
  import type { Writable } from 'svelte/store';

  export let label: string;
  export let store: Writable<string>;
  export let placeholder: string;
  export let hint: string = '';

  function extractUsername(input: string): string {
    if (input.includes('github.com/')) {
      const match = input.match(/github\.com\/([^\/]+)/);
      return match ? match[1] : input;
    }
    return input;
  }

  function handleInput(event: Event) {
    const target = event.target as HTMLInputElement;
    store.set(extractUsername(target.value));
  }
</script>

<div>
  <label for={label} class="mb-2 block text-sm font-medium text-gray-300">{label}</label>
  <input
    id={label}
    type="text"
    {placeholder}
    bind:value={$store}
    oninput={handleInput}
    class="w-full rounded-lg border border-gray-600 bg-gray-800 px-3 py-2 text-white placeholder-gray-400 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-blue-500"
  />
  {#if hint}
    <p class="mt-1 text-xs text-gray-400">{hint}</p>
  {/if}
</div>