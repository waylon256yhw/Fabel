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
</script>

<div class="flex h-full flex-col gap-4 p-4">
  <div class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
    Conversations
  </div>

  <div class="flex flex-col gap-3">
    <div class="flex flex-col gap-1.5">
      <label for="sel-char" class="text-[11px] uppercase tracking-wider text-muted-foreground">Character</label>
      <select
        id="sel-char"
        bind:value={selectedChar}
        class="rounded-md border border-border bg-background px-2 py-1.5 text-sm text-foreground"
      >
        {#each chat.state.characters as c}
          <option value={c.id}>{c.name}</option>
        {/each}
      </select>
    </div>

    <div class="flex flex-col gap-1.5">
      <label for="sel-preset" class="text-[11px] uppercase tracking-wider text-muted-foreground">Preset</label>
      <select
        id="sel-preset"
        bind:value={selectedPreset}
        class="rounded-md border border-border bg-background px-2 py-1.5 text-sm text-foreground"
      >
        {#each chat.state.presets as p}
          <option value={p.id}>{p.name}</option>
        {/each}
      </select>
    </div>
  </div>

  <button
    class="rounded-md border border-border bg-secondary px-3 py-2 text-left text-sm text-foreground transition-colors hover:bg-accent"
    onclick={() => chat.newConversation(selectedChar, selectedPreset)}
  >
    + New conversation
  </button>

  <div class="flex-1"></div>

  <button
    class="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-left text-xs text-destructive transition-colors hover:bg-destructive/20"
    onclick={chat.reset}
  >
    Reset app
  </button>
</div>
