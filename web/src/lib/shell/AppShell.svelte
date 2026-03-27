<script lang="ts">
  import type { Snippet } from 'svelte'
  import TopBar from './TopBar.svelte'
  import Rail from './Rail.svelte'
  import type { AppState } from '../stores/app.svelte.ts'

  let { app, children }: { app: AppState; children: Snippet } = $props()
</script>

<div class="flex h-screen flex-col overflow-hidden bg-background text-foreground">
  <TopBar {app} />
  <div class="flex flex-1 overflow-hidden">
    <Rail {app} />
    {@render children()}
  </div>
</div>
