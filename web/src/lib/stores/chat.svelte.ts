export interface Character {
  id: number;
  name: string;
  description: string;
  first_mes: string;
}

export interface Preset {
  id: number;
  name: string;
  model: string;
}

export interface Message {
  id: number;
  role: "user" | "assistant" | "system";
  content: string;
  created_at: string;
}

export interface ConversationDetail {
  id: number;
  character: Character;
  preset: Preset;
  messages: Message[];
}

export function createChatState() {
  const state = $state({
    characters: [] as Character[],
    presets: [] as Preset[],
    conversation: null as ConversationDetail | null,
    draft: "",
    streaming: false,
    streamingContent: "",
    error: null as string | null,
    showPromptDrawer: false,
    promptMessages: [] as { role: string; content: string }[],
  });

  async function bootstrap() {
    const res = await fetch("/api/bootstrap");
    const data = await res.json();
    state.characters = data.characters ?? [];
    state.presets = data.presets ?? [];
    if (data.seeded_conversation) {
      state.conversation = data.seeded_conversation;
    }
  }

  async function newConversation(characterId: number, presetId: number) {
    const res = await fetch("/api/conversations", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ character_id: characterId, preset_id: presetId }),
    });
    const data = await res.json();
    state.conversation = data.conversation;
  }

  async function send() {
    if (!state.conversation || !state.draft.trim() || state.streaming) return;

    const content = state.draft.trim();
    state.draft = "";
    state.streaming = true;
    state.streamingContent = "";
    state.error = null;

    // Optimistically add user message
    state.conversation.messages.push({
      id: Date.now(),
      role: "user",
      content,
      created_at: new Date().toISOString(),
    });

    try {
      const res = await fetch(`/api/conversations/${state.conversation.id}/send`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content }),
      });

      const reader = res.body!.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");
        buffer = lines.pop() ?? "";

        for (const line of lines) {
          if (!line.startsWith("data: ")) continue;
          const data = line.slice(6);
          if (data === "[DONE]") {
            state.conversation!.messages.push({
              id: Date.now() + 1,
              role: "assistant",
              content: state.streamingContent,
              created_at: new Date().toISOString(),
            });
            state.streamingContent = "";
            break;
          }
          try {
            const { delta } = JSON.parse(data);
            state.streamingContent += delta;
          } catch {}
        }
      }
    } catch (e) {
      state.error = String(e);
    } finally {
      state.streaming = false;
    }
  }

  async function loadPrompt() {
    if (!state.conversation) return;
    const res = await fetch(`/api/conversations/${state.conversation.id}/prompt`);
    const data = await res.json();
    state.promptMessages = data.messages ?? [];
    state.showPromptDrawer = true;
  }

  async function reset() {
    await fetch("/api/conversations", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        character_id: state.characters[0]?.id ?? 1,
        preset_id: state.presets[0]?.id ?? 1,
      }),
    })
      .then((r) => r.json())
      .then((data) => {
        state.conversation = data.conversation;
      });
  }

  return { state, bootstrap, newConversation, send, loadPrompt, reset };
}

export type ChatState = ReturnType<typeof createChatState>;
