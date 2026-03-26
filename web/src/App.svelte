<script lang="ts">
  import { setContext, onMount } from 'svelte'
  import { createChatState } from './lib/chat.svelte.ts'
  import Sidebar from './lib/Sidebar.svelte'
  import MessageList from './lib/MessageList.svelte'
  import ChatInput from './lib/ChatInput.svelte'

  const chat = createChatState()
  setContext('chat', chat)

  const { state, bootstrap, loadPrompt } = chat

  onMount(() => bootstrap())
</script>

<div class="shell">
  <Sidebar {chat} />

  <div class="main">
    {#if state.conversation}
      <header>
        <span class="char-name">{state.conversation.character.name}</span>
        <span class="model">{state.conversation.preset.model}</span>
      </header>
    {/if}

    <MessageList {chat} />
    <ChatInput {chat} />
  </div>
</div>

<!-- Prompt Inspector Drawer -->
{#if state.showPromptDrawer}
  <div class="overlay" onclick={() => (state.showPromptDrawer = false)}>
    <div class="drawer" onclick={(e) => e.stopPropagation()}>
      <div class="drawer-header">
        <span>Prompt Inspector</span>
        <button onclick={() => (state.showPromptDrawer = false)}>✕</button>
      </div>
      <div class="drawer-body">
        {#each state.promptMessages as msg}
          <div class="pm">
            <div class="pm-role">{msg.role}</div>
            <pre class="pm-content">{msg.content}</pre>
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  .shell {
    height: 100vh;
    display: flex;
    overflow: hidden;
  }
  .main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  header {
    padding: 12px 24px;
    border-bottom: 1px solid #1e1e22;
    display: flex;
    align-items: center;
    gap: 12px;
    background: #0f0f11;
    flex-shrink: 0;
  }
  .char-name { font-weight: 600; font-size: 15px; }
  .model { font-size: 11px; color: #555; font-family: monospace; }

  /* Prompt drawer */
  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.6);
    z-index: 50;
    display: flex;
    justify-content: flex-end;
  }
  .drawer {
    width: min(680px, 90vw);
    height: 100%;
    background: #111113;
    border-left: 1px solid #222;
    display: flex;
    flex-direction: column;
  }
  .drawer-header {
    padding: 16px 20px;
    border-bottom: 1px solid #222;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
    font-size: 14px;
  }
  .drawer-header button {
    background: none;
    border: none;
    color: #666;
    font-size: 16px;
    cursor: pointer;
  }
  .drawer-header button:hover { color: #aaa; }
  .drawer-body {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .pm { border: 1px solid #222; border-radius: 6px; overflow: hidden; }
  .pm-role {
    background: #1a1a1d;
    padding: 6px 12px;
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: #c9a86c;
  }
  .pm-content {
    padding: 12px;
    font-size: 12px;
    line-height: 1.6;
    white-space: pre-wrap;
    word-break: break-word;
    color: #b0aea8;
    font-family: 'Menlo', 'Monaco', monospace;
  }
</style>
