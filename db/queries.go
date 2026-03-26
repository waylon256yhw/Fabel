package db

import "database/sql"

type Character struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Personality string `json:"personality"`
	Scenario    string `json:"scenario"`
	FirstMes    string `json:"first_mes"`
	MesExample  string `json:"mes_example"`
}

type Preset struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Model        string  `json:"model"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
	SystemPrompt string  `json:"system_prompt"`
}

type Conversation struct {
	ID          int    `json:"id"`
	CharacterID int    `json:"character_id"`
	PresetID    int    `json:"preset_id"`
	CreatedAt   string `json:"created_at"`
}

type Message struct {
	ID             int    `json:"id"`
	ConversationID int    `json:"conversation_id"`
	Role           string `json:"role"`
	Content        string `json:"content"`
	ParentID       *int   `json:"parent_id"`
	ActiveChildID  *int   `json:"active_child_id"`
	CreatedAt      string `json:"created_at"`
}

type ConversationDetail struct {
	Conversation
	Character Character `json:"character"`
	Preset    Preset    `json:"preset"`
	Messages  []Message `json:"messages"`
}

func (d *DB) ListCharacters() ([]Character, error) {
	rows, err := d.sql.Query(`SELECT id, name, description, personality, scenario, first_mes, mes_example FROM characters ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Character
	for rows.Next() {
		var c Character
		rows.Scan(&c.ID, &c.Name, &c.Description, &c.Personality, &c.Scenario, &c.FirstMes, &c.MesExample)
		out = append(out, c)
	}
	return out, rows.Err()
}

func (d *DB) ListPresets() ([]Preset, error) {
	rows, err := d.sql.Query(`SELECT id, name, model, temperature, max_tokens, system_prompt FROM presets ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Preset
	for rows.Next() {
		var p Preset
		rows.Scan(&p.ID, &p.Name, &p.Model, &p.Temperature, &p.MaxTokens, &p.SystemPrompt)
		out = append(out, p)
	}
	return out, rows.Err()
}

func (d *DB) GetCharacter(id int) (*Character, error) {
	var c Character
	err := d.sql.QueryRow(
		`SELECT id, name, description, personality, scenario, first_mes, mes_example FROM characters WHERE id = ?`, id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.Personality, &c.Scenario, &c.FirstMes, &c.MesExample)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *DB) CreateConversation(characterID, presetID int) (*Conversation, error) {
	res, err := d.sql.Exec(`INSERT INTO conversations (character_id, preset_id) VALUES (?, ?)`, characterID, presetID)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Conversation{ID: int(id), CharacterID: characterID, PresetID: presetID}, nil
}

func (d *DB) GetFirstConversation() (*ConversationDetail, error) {
	var id int
	err := d.sql.QueryRow(`SELECT id FROM conversations ORDER BY id LIMIT 1`).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return d.GetConversationWithDetails(id)
}

func (d *DB) GetConversationWithDetails(id int) (*ConversationDetail, error) {
	var c ConversationDetail
	err := d.sql.QueryRow(`
		SELECT c.id, c.character_id, c.preset_id, c.created_at,
		       ch.id, ch.name, ch.description, ch.personality, ch.scenario, ch.first_mes, ch.mes_example,
		       p.id, p.name, p.model, p.temperature, p.max_tokens, p.system_prompt
		FROM conversations c
		JOIN characters ch ON ch.id = c.character_id
		JOIN presets p     ON p.id  = c.preset_id
		WHERE c.id = ?`, id,
	).Scan(
		&c.ID, &c.CharacterID, &c.PresetID, &c.CreatedAt,
		&c.Character.ID, &c.Character.Name, &c.Character.Description,
		&c.Character.Personality, &c.Character.Scenario,
		&c.Character.FirstMes, &c.Character.MesExample,
		&c.Preset.ID, &c.Preset.Name, &c.Preset.Model,
		&c.Preset.Temperature, &c.Preset.MaxTokens, &c.Preset.SystemPrompt,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := d.GetMessages(id)
	if err != nil {
		return nil, err
	}
	c.Messages = msgs
	return &c, nil
}

func (d *DB) GetMessages(conversationID int) ([]Message, error) {
	rows, err := d.sql.Query(
		`SELECT id, conversation_id, role, content, parent_id, active_child_id, created_at
		 FROM messages WHERE conversation_id = ? ORDER BY id`,
		conversationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Message
	for rows.Next() {
		var m Message
		rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.ParentID, &m.ActiveChildID, &m.CreatedAt)
		out = append(out, m)
	}
	return out, rows.Err()
}

func (d *DB) AppendMessage(conversationID int, role, content string, parentID *int) error {
	_, err := d.sql.Exec(
		`INSERT INTO messages (conversation_id, role, content, parent_id) VALUES (?, ?, ?, ?)`,
		conversationID, role, content, parentID,
	)
	return err
}
