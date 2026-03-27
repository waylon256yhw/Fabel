<script lang="ts">
  import type { ChatState } from '../../stores/chat.svelte.ts'
  import { tick } from 'svelte'

  let { chat }: { chat: ChatState } = $props()

  let scrollEl: HTMLElement | undefined = $state()

  const _msgLen = $derived(chat.state.conversation?.messages.length ?? 0)
  const _streaming = $derived(chat.state.streamingContent)

  $effect(() => {
    void _msgLen
    void _streaming
    tick().then(() => {
      if (scrollEl) scrollEl.scrollTop = scrollEl.scrollHeight
    })
  })

  function escapeHtml(s: string): string {
    return s
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
  }

  function renderContent(text: string) {
    return escapeHtml(text)
      .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.+?)\*/g, '<em>$1</em>')
      .replace(/\n/g, '<br>')
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
      e.preventDefault()
      chat.send()
    }
  }

  function insertText(before: string, after: string) {
    chat.state.draft += before + after
  }
</script>

<!-- Messages -->
<div class="flex-1 overflow-y-auto scroll-smooth" bind:this={scrollEl}>
  {#if !chat.state.conversation}
    <div class="flex h-full items-center justify-center text-sm text-muted-foreground">
      Select a character and start a conversation.
    </div>
  {:else}
    <div class="prose-fabel mx-auto flex max-w-[800px] flex-col gap-12 px-12 py-16">
      {#each chat.state.conversation.messages as msg (msg.id)}
        <div class="flex flex-col gap-2">
          <div
            class="text-[0.85rem] font-semibold tracking-[0.05em]
              {msg.role === 'user' ? 'text-primary' : 'text-muted-foreground'}"
          >
            {msg.role === 'user' ? 'You' : chat.state.conversation.character.name}
          </div>
          <div>{@html renderContent(msg.content)}</div>
        </div>
      {/each}

      {#if chat.state.streaming && chat.state.streamingContent}
        <div class="flex flex-col gap-2">
          <div class="text-[0.85rem] font-semibold tracking-[0.05em] text-muted-foreground">
            {chat.state.conversation.character.name}
          </div>
          <div>{@html renderContent(chat.state.streamingContent)}<span class="inline-block animate-pulse text-primary">|</span></div>
        </div>
      {/if}

      {#if chat.state.error}
        <div class="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive">
          {chat.state.error}
        </div>
      {/if}
    </div>
  {/if}
</div>

<!-- Input area -->
<div class="shrink-0 border-t border-border px-6 py-1.5 pb-2">
  <div class="mx-auto max-w-[800px]">
    <!-- Plugin assistant bar -->
    <div class="flex gap-0.5 pl-10 pb-0.5">
      <button class="rounded px-1.5 py-0.5 font-story text-[0.8rem] text-muted-foreground transition-colors hover:bg-accent hover:text-foreground" onclick={() => insertText('「', '」')} title="Quote">「 」</button>
      <button class="rounded px-1.5 py-0.5 font-story text-[0.8rem] text-muted-foreground transition-colors hover:bg-accent hover:text-foreground" onclick={() => insertText('『', '』')} title="Double Quote">『 』</button>
      <button class="rounded px-1.5 py-0.5 font-story text-[0.8rem] text-muted-foreground transition-colors hover:bg-accent hover:text-foreground" onclick={() => insertText('*', '*')} title="Action">* *</button>
      <button class="rounded px-1.5 py-0.5 font-story text-[0.8rem] text-muted-foreground transition-colors hover:bg-accent hover:text-foreground" onclick={() => insertText('(', ')')} title="OOC">(  )</button>
    </div>

    <!-- Input row -->
    <div class="flex items-end gap-1.5">
      <button
        class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        title="Inspect prompt"
        onclick={chat.loadPrompt}
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
      </button>

      <textarea
        bind:value={chat.state.draft}
        onkeydown={onKeydown}
        placeholder="Continue the story..."
        rows="1"
        disabled={chat.state.streaming || !chat.state.conversation}
        class="flex-1 resize-none bg-transparent px-2 py-1 font-story text-base leading-relaxed text-foreground placeholder:font-sans placeholder:text-[0.95rem] placeholder:text-muted-foreground/50 focus:outline-none disabled:opacity-40"
      ></textarea>

      <button
        class="flex h-8 w-8 shrink-0 items-center justify-center rounded-[10px] border border-primary/30 bg-primary/15 text-primary transition-all hover:-translate-y-0.5 hover:bg-primary hover:text-primary-foreground hover:shadow-[0_4px_12px_rgba(196,164,106,0.2)] disabled:opacity-30 disabled:hover:translate-y-0"
        onclick={chat.send}
        disabled={chat.state.streaming || !chat.state.draft.trim() || !chat.state.conversation}
        title="Send"
      >
        {#if chat.state.streaming}
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><rect x="6" y="6" width="12" height="12"/></svg>
        {:else}
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22,2 15,22 11,13 2,9"/></svg>
        {/if}
      </button>
    </div>
  </div>
</div>

<!-- Prompt Inspector Drawer -->
{#if chat.state.showPromptDrawer}
  <div
    class="fixed inset-0 z-50 flex justify-end bg-black/60"
    onclick={() => (chat.state.showPromptDrawer = false)}
    role="presentation"
  >
    <div
      class="flex h-full w-full max-w-[680px] flex-col border-l border-border bg-card"
      onclick={(e) => e.stopPropagation()}
      role="presentation"
    >
      <div class="flex items-center justify-between border-b border-border px-5 py-4">
        <span class="text-sm font-semibold">Prompt Inspector</span>
        <button
          class="text-muted-foreground hover:text-foreground"
          onclick={() => (chat.state.showPromptDrawer = false)}
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
      <div class="flex flex-1 flex-col gap-3 overflow-y-auto p-4">
        {#each chat.state.promptMessages as msg}
          <div class="overflow-hidden rounded-md border border-border">
            <div class="bg-secondary px-3 py-1.5 text-[11px] font-semibold uppercase tracking-wider text-primary">
              {msg.role}
            </div>
            <pre class="whitespace-pre-wrap break-words p-3 font-mono text-xs leading-relaxed text-muted-foreground">{msg.content}</pre>
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}
