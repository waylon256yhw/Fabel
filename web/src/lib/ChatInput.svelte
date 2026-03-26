<script lang="ts">
  import type { ChatState } from './chat.svelte.ts'

  let { chat }: { chat: ChatState } = $props()

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      chat.send()
    }
  }
</script>

<div class="bar">
  <div class="tools">
    <button class="tool" onclick={chat.loadPrompt} title="Inspect prompt">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
        <polyline points="14,2 14,8 20,8"/>
        <line x1="16" y1="13" x2="8" y2="13"/>
        <line x1="16" y1="17" x2="8" y2="17"/>
      </svg>
      Prompt
    </button>
  </div>

  <div class="input-row">
    <textarea
      bind:value={chat.state.draft}
      onkeydown={onKeydown}
      placeholder="Type a message… (Enter to send, Shift+Enter for newline)"
      rows="3"
      disabled={chat.state.streaming || !chat.state.conversation}
    ></textarea>

    <button
      class="send"
      onclick={chat.send}
      disabled={chat.state.streaming || !chat.state.draft.trim() || !chat.state.conversation}
    >
      {#if chat.state.streaming}
        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
          <rect x="6" y="6" width="12" height="12"/>
        </svg>
      {:else}
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="22" y1="2" x2="11" y2="13"/>
          <polygon points="22,2 15,22 11,13 2,9"/>
        </svg>
      {/if}
    </button>
  </div>
</div>

<style>
  .bar {
    border-top: 1px solid #222;
    padding: 12px 24px 16px;
    background: #0f0f11;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .tools { display: flex; gap: 8px; }
  .tool {
    display: flex;
    align-items: center;
    gap: 5px;
    background: none;
    border: 1px solid #2a2a2e;
    border-radius: 5px;
    color: #666;
    font-size: 12px;
    padding: 4px 8px;
    cursor: pointer;
  }
  .tool:hover { color: #999; border-color: #3a3a3e; }
  .input-row {
    display: flex;
    gap: 10px;
    align-items: flex-end;
  }
  textarea {
    flex: 1;
    background: #1a1a1d;
    border: 1px solid #2a2a2e;
    border-radius: 8px;
    color: #e2e0db;
    font-size: 14px;
    line-height: 1.5;
    padding: 10px 14px;
    resize: none;
    font-family: inherit;
  }
  textarea:focus { outline: none; border-color: #3a4a5e; }
  textarea:disabled { opacity: 0.5; }
  .send {
    width: 42px;
    height: 42px;
    border-radius: 8px;
    background: #c9a86c;
    border: none;
    color: #0d0d0f;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
  .send:hover { background: #d4b87a; }
  .send:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
