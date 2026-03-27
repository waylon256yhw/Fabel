<script lang="ts">
  import { onMount, setContext } from 'svelte'
  import { AppState } from './lib/stores/app.svelte.ts'
  import AuthGate from './lib/auth/AuthGate.svelte'
  import AppShell from './lib/shell/AppShell.svelte'
  import ChatWorkspace from './lib/workspaces/chat/ChatWorkspace.svelte'
  import CharactersWorkspace from './lib/workspaces/characters/CharactersWorkspace.svelte'
  import WorldbooksWorkspace from './lib/workspaces/worldbooks/WorldbooksWorkspace.svelte'
  import SearchWorkspace from './lib/workspaces/search/SearchWorkspace.svelte'

  const app = new AppState()
  setContext('app', app)

  onMount(() => app.checkAuth())
</script>

<AuthGate {app}>
  <AppShell {app}>
    {#if app.activeWorkspace === 'chat'}
      <ChatWorkspace {app} />
    {:else if app.activeWorkspace === 'characters'}
      <CharactersWorkspace {app} />
    {:else if app.activeWorkspace === 'worldbooks'}
      <WorldbooksWorkspace {app} />
    {:else if app.activeWorkspace === 'search'}
      <SearchWorkspace {app} />
    {/if}
  </AppShell>
</AuthGate>
