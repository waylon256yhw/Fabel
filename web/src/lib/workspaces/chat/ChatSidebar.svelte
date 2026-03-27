<script lang="ts">
  import type { ChatState } from '../../stores/chat.svelte.ts'

  let { chat }: { chat: ChatState } = $props()

  let selectedChar = $state(0)
  let selectedPreset = $state(0)

  $effect(() => {
    if (chat.state.characters.length && !selectedChar) {
      selectedChar = chat.state.characters[0].id
    }
    if (chat.state.presets.length && !selectedPreset) {
      selectedPreset = chat.state.presets[0].id
    }
  })

  function startConversation() {
    chat.newConversation(selectedChar, selectedPreset)
  }
</script>

<div class="flex h-full flex-col">
  <!-- Header -->
  <div class="flex items-center justify-between px-5 pb-3 pt-5">
    <span class="text-[0.75rem] font-semibold uppercase tracking-[0.08em] text-muted-foreground">Active Threads</span>
    <button
      class="flex h-6 w-6 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
      title="New Chat"
      onclick={startConversation}
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
    </button>
  </div>

  <!-- Thread list -->
  <div class="flex flex-1 flex-col gap-2 overflow-y-auto px-3 pb-3">
    {#if chat.state.conversation}
      <!-- Active conversation card -->
      <button
        class="flex items-center gap-3 rounded-[10px] border border-border bg-secondary p-3 text-left shadow-sm transition-colors"
      >
        <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-[10px] border border-border bg-background font-story text-lg font-semibold text-primary">
          {chat.state.conversation.character.name.charAt(0)}
        </div>
        <div class="flex min-w-0 flex-1 flex-col gap-1">
          <span class="truncate text-[0.95rem] font-semibold text-foreground">
            {chat.state.conversation.character.name}
          </span>
          <span class="truncate text-[0.8rem] text-muted-foreground">
            {chat.state.conversation.messages.at(-1)?.content.slice(0, 40) ?? '...'}
          </span>
        </div>
      </button>
    {/if}

    <!-- Other characters as potential threads -->
    {#each chat.state.characters as c}
      {#if c.id !== chat.state.conversation?.character.id}
        <button
          class="flex items-center gap-3 rounded-[10px] border border-transparent p-3 text-left transition-colors hover:bg-secondary"
          onclick={() => chat.newConversation(c.id, selectedPreset)}
        >
          <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-[10px] border border-border bg-background font-story text-lg text-muted-foreground">
            {c.name.charAt(0)}
          </div>
          <div class="flex min-w-0 flex-1 flex-col gap-1">
            <span class="truncate text-[0.95rem] font-semibold text-foreground">{c.name}</span>
            <span class="truncate text-[0.8rem] text-muted-foreground">
              {c.first_mes?.slice(0, 40) ?? 'Start a conversation...'}
            </span>
          </div>
        </button>
      {/if}
    {/each}
  </div>

  <!-- Preset selector (compact, bottom) -->
  <div class="border-t border-border px-4 py-3">
    <div class="flex items-center gap-2">
      <label for="sel-preset" class="text-[10px] uppercase tracking-wider text-muted-foreground">Preset</label>
      <select
        id="sel-preset"
        bind:value={selectedPreset}
        class="flex-1 truncate rounded-md border border-border bg-background px-2 py-1 text-xs text-foreground"
      >
        {#each chat.state.presets as p}
          <option value={p.id}>{p.name}</option>
        {/each}
      </select>
    </div>
  </div>
</div>
