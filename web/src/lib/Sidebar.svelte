<script lang="ts">
  import type { ChatState } from './chat.svelte.ts'

  let { chat }: { chat: ChatState } = $props()

  let selectedChar = $derived(chat.state.characters[0]?.id ?? 0)
  let selectedPreset = $derived(chat.state.presets[0]?.id ?? 0)
</script>

<aside>
  <div class="brand">Fabel</div>

  <section>
    <label for="sel-char">Character</label>
    <select id="sel-char" bind:value={selectedChar}>
      {#each chat.state.characters as c}
        <option value={c.id}>{c.name}</option>
      {/each}
    </select>
  </section>

  <section>
    <label for="sel-preset">Preset</label>
    <select id="sel-preset" bind:value={selectedPreset}>
      {#each chat.state.presets as p}
        <option value={p.id}>{p.name}</option>
      {/each}
    </select>
  </section>

  <button onclick={() => chat.newConversation(selectedChar, selectedPreset)}>
    + New conversation
  </button>

  <div class="spacer"></div>

  <button class="danger" onclick={chat.reset}>Reset app</button>
</aside>

<style>
  aside {
    width: 220px;
    flex-shrink: 0;
    height: 100%;
    background: #121214;
    border-right: 1px solid #222;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .brand {
    font-size: 18px;
    font-weight: 700;
    letter-spacing: 0.05em;
    color: #c9a86c;
    padding-bottom: 8px;
    border-bottom: 1px solid #222;
  }
  section { display: flex; flex-direction: column; gap: 6px; }
  label { font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em; color: #666; }
  select {
    background: #1a1a1d;
    border: 1px solid #333;
    border-radius: 6px;
    color: #e2e0db;
    padding: 6px 8px;
    font-size: 13px;
  }
  button {
    background: #1e2a3a;
    border: 1px solid #2a3d55;
    border-radius: 6px;
    color: #7eb8e8;
    padding: 8px 12px;
    font-size: 13px;
    cursor: pointer;
    text-align: left;
  }
  button:hover { background: #253245; }
  .spacer { flex: 1; }
  .danger {
    background: #2a1a1a;
    border-color: #4a2020;
    color: #e87070;
    font-size: 12px;
  }
  .danger:hover { background: #341f1f; }
</style>
