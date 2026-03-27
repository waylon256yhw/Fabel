<script context="module" lang="ts">
  import { createChatState } from '../../stores/chat.svelte.ts'

  // Module-level singleton — survives component mount/unmount cycles
  const chat = createChatState()
  let bootstrapped = false
</script>

<script lang="ts">
  import { setContext, onMount } from 'svelte'
  import type { AppState } from '../../stores/app.svelte.ts'
  import Scene from '../../shell/Scene.svelte'
  import ChatSidebar from './ChatSidebar.svelte'
  import ChatStage from './ChatStage.svelte'
  import ChatInspector from './ChatInspector.svelte'

  let { app }: { app: AppState } = $props()

  setContext('chat', chat)

  // Update TopBar center reactively
  $effect(() => {
    const conv = chat.state.conversation
    app.topBarTitle = conv?.character.name ?? ""
    app.topBarSubtitle = conv?.preset.model ?? ""
  })

  onMount(() => {
    if (!bootstrapped) {
      bootstrapped = true
      chat.bootstrap()
    }
  })
</script>

{#snippet sidebar()}
  <ChatSidebar {chat} />
{/snippet}

{#snippet stage()}
  <ChatStage {chat} />
{/snippet}

{#snippet inspector()}
  <ChatInspector {chat} />
{/snippet}

<Scene {app} {sidebar} {stage} {inspector} />
