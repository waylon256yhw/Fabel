<script lang="ts">
  import type { Snippet } from 'svelte'
  import type { AppState } from '../stores/app.svelte.ts'

  let { app, sidebar, stage, inspector }: {
    app: AppState
    sidebar?: Snippet
    stage: Snippet
    inspector?: Snippet
  } = $props()
</script>

<div class="flex flex-1 overflow-hidden">
  {#if sidebar && app.sidebarOpen}
    <aside class="flex w-80 shrink-0 flex-col overflow-hidden border-r border-border bg-card">
      {@render sidebar()}
    </aside>
  {/if}

  <main class="flex min-w-0 flex-1 flex-col overflow-hidden">
    {@render stage()}
  </main>

  {#if inspector && app.inspectorOpen}
    <aside class="flex w-80 shrink-0 flex-col overflow-hidden border-l border-border bg-card">
      {@render inspector()}
    </aside>
  {/if}
</div>
