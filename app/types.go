package app

import "fabel/db"

type BootstrapResponse struct {
	Characters  []db.Character `json:"characters"`
	Presets     []db.Preset    `json:"presets"`
	Seeded      *db.ConversationDetail `json:"seeded_conversation,omitempty"`
}

type CreateConversationRequest struct {
	CharacterID int `json:"character_id"`
	PresetID    int `json:"preset_id"`
}

type GetConversationResponse struct {
	Conversation *db.ConversationDetail `json:"conversation"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}

type SSEDelta struct {
	Delta string `json:"delta"`
}

type PromptResponse struct {
	Messages []map[string]string `json:"messages"`
}
