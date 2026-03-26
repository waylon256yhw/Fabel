package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"fabel/db"
	"fabel/llm"
)

func (a *App) handleBootstrap(w http.ResponseWriter, r *http.Request) {
	characters, err := a.db.ListCharacters()
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	presets, err := a.db.ListPresets()
	if err != nil {
		httpErr(w, err, 500)
		return
	}

	// Return the first seeded conversation so the UI has something to open immediately.
	conv, _ := a.db.GetFirstConversation()

	writeJSON(w, BootstrapResponse{
		Characters: characters,
		Presets:    presets,
		Seeded:     conv,
	})
}

func (a *App) handleCreateConversation(w http.ResponseWriter, r *http.Request) {
	var req CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr(w, err, 400)
		return
	}

	conv, err := a.db.CreateConversation(req.CharacterID, req.PresetID)
	if err != nil {
		httpErr(w, err, 500)
		return
	}

	char, err := a.db.GetCharacter(req.CharacterID)
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	if err := a.db.AppendMessage(conv.ID, "assistant", char.FirstMes, nil); err != nil {
		httpErr(w, err, 500)
		return
	}

	detail, err := a.db.GetConversationWithDetails(conv.ID)
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	writeJSON(w, GetConversationResponse{Conversation: detail})
}

func (a *App) handleGetConversation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpErr(w, err, 400)
		return
	}
	detail, err := a.db.GetConversationWithDetails(id)
	if err != nil {
		httpErr(w, err, 404)
		return
	}
	writeJSON(w, GetConversationResponse{Conversation: detail})
}

func (a *App) handleSendMessageStream(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpErr(w, err, 400)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr(w, err, 400)
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		httpErr(w, fmt.Errorf("content required"), 400)
		return
	}

	// Get state BEFORE appending the new user message.
	conv, err := a.db.GetConversationWithDetails(id)
	if err != nil {
		httpErr(w, err, 404)
		return
	}

	promptMsgs := a.buildPrompt(conv, req.Content)

	if err := a.db.AppendMessage(id, "user", req.Content, nil); err != nil {
		httpErr(w, err, 500)
		return
	}

	ch, err := a.llm.Stream(r.Context(), promptMsgs, conv.Preset.Model, conv.Preset.Temperature, conv.Preset.MaxTokens)
	if err != nil {
		httpErr(w, err, 502)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, _ := w.(http.Flusher)
	var full strings.Builder

	for delta := range ch {
		full.WriteString(delta)
		data, _ := json.Marshal(SSEDelta{Delta: delta})
		fmt.Fprintf(w, "data: %s\n\n", data)
		if flusher != nil {
			flusher.Flush()
		}
	}

	a.db.AppendMessage(id, "assistant", full.String(), nil)
	fmt.Fprintf(w, "data: [DONE]\n\n")
	if flusher != nil {
		flusher.Flush()
	}
}

func (a *App) handleGetPrompt(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpErr(w, err, 400)
		return
	}
	conv, err := a.db.GetConversationWithDetails(id)
	if err != nil {
		httpErr(w, err, 404)
		return
	}
	msgs := a.buildPrompt(conv, "(next user message would appear here)")
	out := make([]map[string]string, len(msgs))
	for i, m := range msgs {
		out[i] = map[string]string{"role": m.Role, "content": m.Content}
	}
	writeJSON(w, PromptResponse{Messages: out})
}

func (a *App) buildPrompt(conv *db.ConversationDetail, userContent string) []llm.Message {
	replacer := strings.NewReplacer("{{char}}", conv.Character.Name, "{{user}}", "User")

	system := replacer.Replace(conv.Preset.SystemPrompt)
	if conv.Character.Description != "" {
		system += "\n\n[Character: " + replacer.Replace(conv.Character.Description) + "]"
	}
	if conv.Character.Personality != "" {
		system += "\n[Personality: " + replacer.Replace(conv.Character.Personality) + "]"
	}
	if conv.Character.Scenario != "" {
		system += "\n[Scenario: " + replacer.Replace(conv.Character.Scenario) + "]"
	}

	messages := []llm.Message{{Role: "system", Content: system}}
	for _, msg := range conv.Messages {
		messages = append(messages, llm.Message{Role: msg.Role, Content: msg.Content})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userContent})
	return messages
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func httpErr(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}
