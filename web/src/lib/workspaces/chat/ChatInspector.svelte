<script lang="ts">
  import type { ChatState } from '../../stores/chat.svelte.ts'

  let { chat }: { chat: ChatState } = $props()
</script>

<div class="flex h-full flex-col gap-6 overflow-y-auto p-5">
  {#if chat.state.conversation}
    {@const char = chat.state.conversation.character}

    <!-- Character hero -->
    <div class="flex flex-col items-center gap-3 text-center">
      <div class="flex h-24 w-24 items-center justify-center rounded-2xl border border-border bg-background font-story text-3xl text-primary shadow-md">
        {char.name.charAt(0)}
      </div>
      <div class="text-lg font-semibold text-foreground">{char.name}</div>
    </div>

    <!-- Description -->
    {#if char.description}
      <div class="flex flex-col gap-2">
        <span class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">Description</span>
        <div class="rounded-md border border-border bg-background p-3 font-story text-sm leading-relaxed text-foreground">
          {char.description}
        </div>
      </div>
    {/if}

    <!-- First message preview -->
    {#if char.first_mes}
      <div class="flex flex-col gap-2">
        <span class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">First Message</span>
        <div class="rounded-md border border-border bg-background p-3 font-story text-sm italic leading-relaxed text-muted-foreground">
          {char.first_mes.length > 200 ? char.first_mes.slice(0, 200) + '...' : char.first_mes}
        </div>
      </div>
    {/if}

    <!-- Stats -->
    <div class="flex justify-center gap-6 border-t border-border pt-4">
      <div class="flex flex-col items-center gap-1">
        <span class="text-sm font-semibold text-foreground">{chat.state.conversation.messages.length}</span>
        <span class="text-[10px] uppercase tracking-wider text-muted-foreground">Messages</span>
      </div>
      <div class="flex flex-col items-center gap-1">
        <span class="text-sm font-semibold text-foreground">{chat.state.conversation.preset.model}</span>
        <span class="text-[10px] uppercase tracking-wider text-muted-foreground">Model</span>
      </div>
    </div>
  {:else}
    <div class="flex h-full items-center justify-center text-sm text-muted-foreground">
      No active conversation
    </div>
  {/if}
</div>
