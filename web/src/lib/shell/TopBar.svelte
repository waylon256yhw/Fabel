<script lang="ts">
  import type { AppState } from '../stores/app.svelte.ts'

  let { app }: { app: AppState } = $props()
</script>

<header class="flex h-[52px] shrink-0 items-center justify-between border-b border-border bg-card px-4">
  <!-- Left: sidebar toggle + logo -->
  <div class="flex items-center gap-3 min-w-[220px]">
    <button
      class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
      title="Toggle sidebar"
      onclick={() => app.toggleSidebar()}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="9" y1="3" x2="9" y2="21"/></svg>
    </button>
    <div class="flex items-center gap-2.5">
      <!-- Fabel Logo SVG -->
      <svg xmlns="http://www.w3.org/2000/svg" width="26" height="26" viewBox="0 0 192 192">
        <path d="m163.4 39.45h-11.78v89.49c0 13.12-10.79 23.68-23.18 23.68h-88.75v10.15c0 13.16 10.01 23.76 23.54 23.76h100.2c13.05 0 23.31-11.18 23.31-23.64v-99.8c0-12.46-10.26-23.64-23.31-23.64z" fill="var(--muted-foreground)"/>
        <path d="m128.8 5.81h-100.9c-12.62 0-22.25 11.5-22.25 23.31v100.4c0 12.22 10.04 23.13 23.05 23.13h11.02v-91.76c0-11.65 8.2-21.23 20.18-21.23h50.4c5.05 0 7.58 2.23 7.58 9.32v7.17c0 6.79-2.28 10.11-8.61 10.11h-27.42c-7.27 0-9.31 1.68-9.31 9.86v3.32c0 7.46 2.07 9.77 8.33 9.77h20.21c6.32 0 7.57 3.2 7.57 9.44v7.02c0 5.82-2.18 8.13-8.5 8.13h-18.24c-7.14 0-9.3 2.14-9.3 11.54v26.88h56.11c12.89 0 22.93-11.33 22.93-23.38v-99.23c0-12.05-9.93-23.77-22.82-23.77z" fill="var(--foreground)"/>
      </svg>
      <span class="font-story text-[1.1rem] font-medium tracking-wide text-foreground">Fabel</span>
    </div>
  </div>

  <!-- Center: workspace-provided title -->
  <div class="flex flex-col items-center gap-0.5">
    {#if app.topBarTitle}
      <span class="text-[0.95rem] font-semibold text-foreground">{app.topBarTitle}</span>
    {/if}
    {#if app.topBarSubtitle}
      <span class="rounded-full border border-border bg-secondary px-2 py-0.5 font-mono text-[0.7rem] text-muted-foreground">
        {app.topBarSubtitle}
      </span>
    {/if}
  </div>

  <!-- Right: global actions -->
  <div class="flex items-center gap-2 min-w-[220px] justify-end">
    <button
      class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
      title="Toggle inspector"
      onclick={() => app.toggleInspector()}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="15" y1="3" x2="15" y2="21"/></svg>
    </button>

    {#if app.user}
      <button
        class="flex h-8 items-center gap-2 rounded-md px-2 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        title="Logout"
        onclick={() => app.logout()}
      >
        <span class="max-w-[100px] truncate text-[13px]">{app.user.display_name || app.user.username}</span>
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
      </button>
    {/if}
  </div>
</header>
