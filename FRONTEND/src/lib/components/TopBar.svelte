<script lang="ts">
  export let toggleSidebar: () => void;
</script>

<div class="border-b border-gray-200 bg-white px-4 py-3 shadow-sm md:px-6">
  <div class="flex items-center">
    <button on:click={toggleSidebar} class="mr-4 text-gray-500 md:hidden">☰</button>
    <h1 class="text-xl font-semibold text-gray-900">Repository Access Checker</h1>
  </div>
</div>
