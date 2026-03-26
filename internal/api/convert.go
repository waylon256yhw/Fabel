package api

import (
	"database/sql"

	"fabel/internal/dbq"
)

func nullInt64ToPtr(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int64)
	return &v
}

func nullTimeToString(n sql.NullTime) string {
	if !n.Valid {
		return ""
	}
	return n.Time.Format("2006-01-02T15:04:05Z")
}

func dbCharacterToAPI(c dbq.ListCharactersRow) Character {
	return Character{
		Id:          int(c.ID),
		Name:        c.Name,
		Description: c.Description,
		Personality: c.Personality,
		Scenario:    c.Scenario,
		FirstMes:    c.FirstMes,
		MesExample:  c.MesExample,
	}
}

func dbGetCharacterToAPI(c dbq.GetCharacterRow) Character {
	return Character{
		Id:          int(c.ID),
		Name:        c.Name,
		Description: c.Description,
		Personality: c.Personality,
		Scenario:    c.Scenario,
		FirstMes:    c.FirstMes,
		MesExample:  c.MesExample,
	}
}

func dbPresetToAPI(p dbq.ListPresetsRow) Preset {
	return Preset{
		Id:           int(p.ID),
		Name:         p.Name,
		Model:        p.Model,
		Temperature:  p.Temperature,
		MaxTokens:    int(p.MaxTokens),
		SystemPrompt: p.SystemPrompt,
	}
}

func dbGetPresetToAPI(p dbq.GetPresetRow) Preset {
	return Preset{
		Id:           int(p.ID),
		Name:         p.Name,
		Model:        p.Model,
		Temperature:  p.Temperature,
		MaxTokens:    int(p.MaxTokens),
		SystemPrompt: p.SystemPrompt,
	}
}

func dbMessageToAPI(m dbq.Message) Message {
	return Message{
		Id:            int(m.ID),
		ConversationId: int(m.ConversationID),
		Role:          MessageRole(m.Role),
		Content:       m.Content,
		ParentId:      nullInt64ToPtr(m.ParentID),
		ActiveChildId: nullInt64ToPtr(m.ActiveChildID),
		CreatedAt:     nullTimeToString(m.CreatedAt),
	}
}

func dbMessagesToAPI(msgs []dbq.Message) []Message {
	out := make([]Message, len(msgs))
	for i, m := range msgs {
		out[i] = dbMessageToAPI(m)
	}
	return out
}

func dbUserToAPI(u dbq.User) User {
	return User{
		Id:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Role:        UserRole(u.Role),
		CreatedAt:   u.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// userIDParam converts a user ID string to sql.NullString for use in queries.
func userIDParam(id string) sql.NullString {
	return sql.NullString{String: id, Valid: true}
}
