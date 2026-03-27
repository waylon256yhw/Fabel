<script lang="ts">
  import type { AppState, WorkspaceId } from '../stores/app.svelte.ts'

  let { app }: { app: AppState } = $props()

  const items: { id: WorkspaceId; label: string; pushBottom?: boolean }[] = [
    { id: 'chat', label: 'Chat' },
    { id: 'characters', label: 'Characters' },
    { id: 'worldbooks', label: 'Worldbooks' },
    { id: 'search', label: 'Search', pushBottom: true },
  ]
</script>

<nav class="flex w-14 shrink-0 flex-col items-center border-r border-border bg-card py-4 gap-2">
  {#each items as item}
    {#if item.pushBottom}
      <div class="mt-auto"></div>
    {/if}
    <button
      class="flex h-10 w-10 items-center justify-center rounded-lg transition-colors
        {app.activeWorkspace === item.id
          ? 'bg-primary/15 text-primary'
          : 'text-muted-foreground hover:text-foreground hover:bg-accent'}"
      title={item.label}
      onclick={() => app.navigate(item.id)}
    >
      {#if item.id === 'chat'}
        <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
      {:else if item.id === 'characters'}
        <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
      {:else if item.id === 'worldbooks'}
        <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>
      {:else if item.id === 'search'}
        <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
      {/if}
    </button>
  {/each}
</nav>
