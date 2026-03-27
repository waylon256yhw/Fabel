<script lang="ts">
  import type { ChatState } from '../../stores/chat.svelte.ts'

  let { chat }: { chat: ChatState } = $props()
</script>

<div class="flex h-full flex-col overflow-y-auto p-5">
  {#if chat.state.conversation}
    {@const char = chat.state.conversation.character}
    {@const conv = chat.state.conversation}

    <!-- Profile hero -->
    <div class="flex flex-col items-center gap-4 pb-6 text-center">
      <div
        class="flex items-center justify-center overflow-hidden rounded-2xl border border-border bg-background font-story text-5xl text-primary"
        style="width: 140px; height: 140px; box-shadow: 0 4px 16px rgba(0,0,0,0.3);"
      >
        {char.name.charAt(0)}
      </div>
      <div class="text-[1.3rem] font-semibold text-foreground">{char.name}</div>
      {#if char.description}
        <div class="font-story text-[0.9rem] italic text-muted-foreground">
          "{char.description.length > 60 ? char.description.slice(0, 60) + '...' : char.description}"
        </div>
      {/if}
    </div>

    <!-- Stats -->
    <div class="flex justify-center gap-4 border-b border-border pb-6">
      <div class="flex flex-col items-center gap-1">
        <span class="text-base font-semibold text-foreground">{conv.messages.length}</span>
        <span class="text-[0.7rem] uppercase tracking-[0.05em] text-muted-foreground">Messages</span>
      </div>
      <div class="flex flex-col items-center gap-1">
        <span class="text-base font-semibold text-foreground">{conv.preset.model}</span>
        <span class="text-[0.7rem] uppercase tracking-[0.05em] text-muted-foreground">Model</span>
      </div>
    </div>

    <!-- Sections -->
    <div class="flex flex-col gap-6 pt-6">
      {#if char.description}
        <div class="flex flex-col gap-3">
          <span class="text-[0.8rem] font-semibold uppercase tracking-[0.08em] text-muted-foreground">
            Persona / Description
          </span>
          <div class="rounded-lg border border-border bg-background p-4 font-story text-[0.9rem] leading-relaxed text-foreground">
            {char.description}
          </div>
        </div>
      {/if}

      {#if char.first_mes}
        <div class="flex flex-col gap-3">
          <span class="text-[0.8rem] font-semibold uppercase tracking-[0.08em] text-muted-foreground">
            First Message
          </span>
          <div class="rounded-lg border border-border bg-background p-4 font-story text-[0.9rem] italic leading-relaxed text-muted-foreground">
            {char.first_mes.length > 300 ? char.first_mes.slice(0, 300) + '...' : char.first_mes}
          </div>
        </div>
      {/if}
    </div>
  {:else}
    <div class="flex h-full items-center justify-center text-sm text-muted-foreground">
      No active conversation
    </div>
  {/if}
</div>
