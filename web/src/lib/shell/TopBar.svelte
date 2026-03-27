<script lang="ts">
  import type { AppState } from '../stores/app.svelte.ts'

  let { app }: { app: AppState } = $props()
</script>

<header class="flex h-13 shrink-0 items-center justify-between border-b border-border bg-card px-4">
  <!-- Left: sidebar toggle + logo -->
  <div class="flex items-center gap-3 min-w-[200px]">
    <button
      class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
      title="Toggle sidebar"
      onclick={() => app.toggleSidebar()}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="9" y1="3" x2="9" y2="21"/></svg>
    </button>
    <div class="flex items-center gap-2">
      <span class="text-lg font-medium tracking-wide" style="color: var(--foreground);">Fabel</span>
    </div>
  </div>

  <!-- Center: workspace-provided title -->
  <div class="flex flex-col items-center gap-0.5">
    {#if app.topBarTitle}
      <span class="text-sm font-semibold text-foreground">{app.topBarTitle}</span>
    {/if}
    {#if app.topBarSubtitle}
      <span class="rounded-full border border-border bg-secondary px-2 py-0.5 font-mono text-[11px] text-muted-foreground">
        {app.topBarSubtitle}
      </span>
    {/if}
  </div>

  <!-- Right: global actions -->
  <div class="flex items-center gap-2 min-w-[200px] justify-end">
    <button
      class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
      title="Toggle inspector"
      onclick={() => app.toggleInspector()}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="15" y1="3" x2="15" y2="21"/></svg>
    </button>

    <!-- User menu -->
    {#if app.user}
      <button
        class="flex h-8 items-center gap-2 rounded-md px-2 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        title="Logout"
        onclick={() => app.logout()}
      >
        <span class="max-w-[100px] truncate">{app.user.display_name || app.user.username}</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
      </button>
    {/if}
  </div>
</header>
