package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"fabel/internal/dbq"
)

type sendMessageRequest struct {
	Content string `json:"content"`
}

type sseDelta struct {
	Delta string `json:"delta"`
}

func (s *Server) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpErr(w, err, 400)
		return
	}

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr(w, err, 400)
		return
	}
	if strings.TrimSpace(req.Content) == "" {
		httpErr(w, fmt.Errorf("content required"), 400)
		return
	}

	// Get conversation state BEFORE appending user message.
	detail, err := s.getConversationDetail(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			httpErr(w, err, 404)
		} else {
			httpErr(w, err, 500)
		}
		return
	}

	promptMsgs := buildPrompt(detail, req.Content)

	// Append user message.
	if err := s.q.AppendMessage(ctx, dbq.AppendMessageParams{
		ConversationID: id,
		Role:           "user",
		Content:        req.Content,
	}); err != nil {
		httpErr(w, err, 500)
		return
	}

	// Stream LLM response.
	ch, err := s.llm.Stream(ctx, promptMsgs, detail.Preset.Model, detail.Preset.Temperature, detail.Preset.MaxTokens)
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
		data, _ := json.Marshal(sseDelta{Delta: delta})
		fmt.Fprintf(w, "data: %s\n\n", data)
		if flusher != nil {
			flusher.Flush()
		}
	}

	// Append assistant message.
	if err := s.q.AppendMessage(ctx, dbq.AppendMessageParams{
		ConversationID: id,
		Role:           "assistant",
		Content:        full.String(),
	}); err != nil {
		errData, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(w, "data: %s\n\n", errData)
		if flusher != nil {
			flusher.Flush()
		}
		return
	}
	fmt.Fprintf(w, "data: [DONE]\n\n")
	if flusher != nil {
		flusher.Flush()
	}
}
