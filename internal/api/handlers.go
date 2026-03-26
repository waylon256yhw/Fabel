package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"fabel/internal/auth"
	"fabel/internal/dbq"
)

func (s *Server) GetBootstrap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := userIDParam(auth.UserFromContext(ctx).ID)

	chars, err := s.q.ListCharacters(ctx, uid)
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	presets, err := s.q.ListPresets(ctx, uid)
	if err != nil {
		httpErr(w, err, 500)
		return
	}

	apiChars := make([]Character, len(chars))
	for i, c := range chars {
		apiChars[i] = dbCharacterToAPI(c)
	}
	apiPresets := make([]Preset, len(presets))
	for i, p := range presets {
		apiPresets[i] = dbPresetToAPI(p)
	}

	resp := BootstrapResponse{
		Characters: apiChars,
		Presets:    apiPresets,
	}

	// Load first conversation if exists.
	firstID, err := s.q.GetFirstConversationID(ctx, uid)
	if err == nil {
		detail, err := s.getConversationDetail(ctx, firstID, uid)
		if err == nil {
			resp.SeededConversation = detail
		}
	}

	writeJSON(w, 200, resp)
}

func (s *Server) CreateConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := userIDParam(auth.UserFromContext(ctx).ID)

	var req CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr(w, err, 400)
		return
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	defer tx.Rollback()

	qtx := s.q.WithTx(tx)

	convID, err := qtx.CreateConversation(ctx, dbq.CreateConversationParams{
		CharacterID: int64(req.CharacterId),
		PresetID:    int64(req.PresetId),
		UserID:      uid,
	})
	if err != nil {
		httpErr(w, err, 500)
		return
	}

	// Append greeting message.
	char, err := qtx.GetCharacter(ctx, dbq.GetCharacterParams{
		ID:     int64(req.CharacterId),
		UserID: uid,
	})
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	if err := qtx.AppendMessage(ctx, dbq.AppendMessageParams{
		ConversationID: convID,
		Role:           "assistant",
		Content:        char.FirstMes,
	}); err != nil {
		httpErr(w, err, 500)
		return
	}

	if err := tx.Commit(); err != nil {
		httpErr(w, err, 500)
		return
	}

	detail, err := s.getConversationDetail(ctx, convID, uid)
	if err != nil {
		httpErr(w, err, 500)
		return
	}

	writeJSON(w, 201, GetConversationResponse{Conversation: *detail})
}

func (s *Server) GetConversation(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	uid := userIDParam(auth.UserFromContext(ctx).ID)

	detail, err := s.getConversationDetail(ctx, int64(id), uid)
	if err != nil {
		if err == sql.ErrNoRows {
			httpErr(w, err, 404)
		} else {
			httpErr(w, err, 500)
		}
		return
	}
	writeJSON(w, 200, GetConversationResponse{Conversation: *detail})
}

func (s *Server) GetConversationPrompt(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	uid := userIDParam(auth.UserFromContext(ctx).ID)

	detail, err := s.getConversationDetail(ctx, int64(id), uid)
	if err != nil {
		if err == sql.ErrNoRows {
			httpErr(w, err, 404)
		} else {
			httpErr(w, err, 500)
		}
		return
	}

	msgs := buildPrompt(detail, "(next user message would appear here)")
	out := make([]PromptMessage, len(msgs))
	for i, m := range msgs {
		out[i] = PromptMessage{Role: m.Role, Content: m.Content}
	}
	writeJSON(w, 200, PromptResponse{Messages: out})
}

func (s *Server) getConversationDetail(ctx context.Context, id int64, uid sql.NullString) (*ConversationDetail, error) {
	conv, err := s.q.GetConversation(ctx, dbq.GetConversationParams{ID: id, UserID: uid})
	if err != nil {
		return nil, err
	}
	char, err := s.q.GetCharacter(ctx, dbq.GetCharacterParams{ID: conv.CharacterID, UserID: uid})
	if err != nil {
		return nil, err
	}
	preset, err := s.q.GetPreset(ctx, dbq.GetPresetParams{ID: conv.PresetID, UserID: uid})
	if err != nil {
		return nil, err
	}
	msgs, err := s.q.GetMessagesByConversation(ctx, id)
	if err != nil {
		return nil, err
	}

	detail := &ConversationDetail{
		Id:          int(conv.ID),
		CharacterId: int(conv.CharacterID),
		PresetId:    int(conv.PresetID),
		CreatedAt:   nullTimeToString(conv.CreatedAt),
		Character:   dbGetCharacterToAPI(char),
		Preset:      dbGetPresetToAPI(preset),
		Messages:    dbMessagesToAPI(msgs),
	}
	return detail, nil
}
