<script lang="ts">
  import type { ChatState } from './chat.svelte.ts'
  import { tick } from 'svelte'

  let { chat }: { chat: ChatState } = $props()

  let listEl: HTMLElement

  $effect(() => {
    // Scroll to bottom on new messages or streaming update
    chat.state.conversation?.messages.length
    chat.state.streamingContent
    tick().then(() => {
      if (listEl) listEl.scrollTop = listEl.scrollHeight
    })
  })

  function renderContent(text: string) {
    // Basic markdown: *italic* and **bold**
    return text
      .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.+?)\*/g, '<em>$1</em>')
      .replace(/\n/g, '<br>')
  }
</script>

<div class="list" bind:this={listEl}>
  {#if !chat.state.conversation}
    <div class="empty">Select a character and start a conversation.</div>
  {:else}
    {#each chat.state.conversation.messages as msg (msg.id)}
      <div class="msg {msg.role}">
        <div class="label">{msg.role === 'user' ? 'You' : chat.state.conversation.character.name}</div>
        <div class="content">{@html renderContent(msg.content)}</div>
      </div>
    {/each}

    {#if chat.state.streaming && chat.state.streamingContent}
      <div class="msg assistant streaming">
        <div class="label">{chat.state.conversation.character.name}</div>
        <div class="content">{@html renderContent(chat.state.streamingContent)}<span class="cursor">▍</span></div>
      </div>
    {/if}

    {#if chat.state.error}
      <div class="error">{chat.state.error}</div>
    {/if}
  {/if}
</div>

<style>
  .list {
    flex: 1;
    overflow-y: auto;
    padding: 24px 32px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    scroll-behavior: smooth;
  }
  .empty {
    color: #555;
    text-align: center;
    margin-top: 40px;
    font-size: 14px;
  }
  .msg { max-width: 720px; width: 100%; }
  .msg.user { align-self: flex-end; }
  .msg.assistant { align-self: flex-start; }
  .label {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: #555;
    margin-bottom: 4px;
  }
  .msg.user .label { text-align: right; }
  .content {
    background: #1a1a1d;
    border: 1px solid #2a2a2e;
    border-radius: 10px;
    padding: 12px 16px;
    font-size: 14px;
    line-height: 1.65;
    color: #d8d6d0;
  }
  .msg.user .content {
    background: #1e2a3a;
    border-color: #2a3d55;
    color: #c5daf2;
  }
  .msg.streaming .content { border-color: #3a3020; }
  .cursor {
    display: inline-block;
    animation: blink 0.8s step-end infinite;
    color: #c9a86c;
  }
  @keyframes blink { 50% { opacity: 0 } }
  .error {
    color: #e87070;
    font-size: 13px;
    padding: 8px 12px;
    background: #2a1a1a;
    border-radius: 6px;
    border: 1px solid #4a2020;
  }
</style>
