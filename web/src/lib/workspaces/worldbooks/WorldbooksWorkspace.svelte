<script lang="ts">
  import type { AppState } from '../../stores/app.svelte.ts'
  import Scene from '../../shell/Scene.svelte'

  let { app }: { app: AppState } = $props()
</script>

{#snippet stage()}
  <div class="flex h-full items-center justify-center">
    <div class="text-center">
      <div class="text-lg font-medium text-foreground">Worldbooks</div>
      <div class="mt-2 text-sm text-muted-foreground">Lorebook management coming in P1</div>
    </div>
  </div>
{/snippet}

<Scene {app} {stage} />
