package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DB struct {
	sql *sql.DB
}

func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)

	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000",
	} {
		if _, err := conn.Exec(pragma); err != nil {
			return nil, err
		}
	}

	d := &DB{sql: conn}
	if err := d.migrate(); err != nil {
		return nil, err
	}
	if err := d.seed(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DB) Close() error {
	return d.sql.Close()
}

func (d *DB) migrate() error {
	_, err := d.sql.Exec(`
	CREATE TABLE IF NOT EXISTS characters (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		name         TEXT    NOT NULL,
		description  TEXT    NOT NULL DEFAULT '',
		personality  TEXT    NOT NULL DEFAULT '',
		scenario     TEXT    NOT NULL DEFAULT '',
		first_mes    TEXT    NOT NULL DEFAULT '',
		mes_example  TEXT    NOT NULL DEFAULT '',
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS presets (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		name          TEXT    NOT NULL,
		model         TEXT    NOT NULL DEFAULT 'gpt-4o-mini',
		temperature   REAL    NOT NULL DEFAULT 1.0,
		max_tokens    INTEGER NOT NULL DEFAULT 1000,
		system_prompt TEXT    NOT NULL DEFAULT '',
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS conversations (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		character_id INTEGER NOT NULL REFERENCES characters(id),
		preset_id    INTEGER NOT NULL REFERENCES presets(id),
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS messages (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id INTEGER NOT NULL REFERENCES conversations(id),
		role            TEXT    NOT NULL CHECK(role IN ('user', 'assistant', 'system')),
		content         TEXT    NOT NULL,
		parent_id       INTEGER REFERENCES messages(id),
		active_child_id INTEGER REFERENCES messages(id),
		created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)
	return err
}

func (d *DB) seed() error {
	var count int
	d.sql.QueryRow("SELECT COUNT(*) FROM characters").Scan(&count)
	if count > 0 {
		return nil
	}

	res, err := d.sql.Exec(
		`INSERT INTO characters (name, description, personality, scenario, first_mes) VALUES (?, ?, ?, ?, ?)`,
		"Seraphina",
		"Seraphina is a mystical forest guardian with flowing silver hair and amber eyes that glow softly with forest magic. Her ethereal gown shimmers as she moves.",
		"Gentle, wise, and compassionate. She speaks with warmth and quiet dignity, though she carries deep sadness about the Shadowfangs corrupting Eldoria.",
		"You have wandered into Eldoria's enchanted forest and found Seraphina's glade — the last sanctuary of peace in the darkening woods.",
		`*A figure steps from between the ancient trees, her silver hair catching the dappled light. She turns with amber eyes that hold centuries of forest wisdom, her gown shimmering with quiet magic.* "Welcome, traveler. You have found the last safe glade in Eldoria." *She gestures at the peaceful clearing around you.* "I am Seraphina, guardian of this forest. Few find their way here — though I sense you needed to."`,
	)
	if err != nil {
		return err
	}
	charID, _ := res.LastInsertId()

	res, err = d.sql.Exec(
		`INSERT INTO presets (name, model, temperature, max_tokens, system_prompt) VALUES (?, ?, ?, ?, ?)`,
		"Default",
		"gpt-4o-mini",
		1.0,
		1000,
		`Write {{char}}'s next reply in a fictional roleplay between {{char}} and {{user}}. Write 1 reply only in internet RP style, italicize actions, and avoid quotation marks. Use markdown. Be proactive, creative, and drive the plot and conversation forward. Write at least 1 paragraph, up to 4. Always stay in character and avoid repetition.`,
	)
	if err != nil {
		return err
	}
	presetID, _ := res.LastInsertId()

	res, err = d.sql.Exec(
		`INSERT INTO conversations (character_id, preset_id) VALUES (?, ?)`,
		charID, presetID,
	)
	if err != nil {
		return err
	}
	convID, _ := res.LastInsertId()

	_, err = d.sql.Exec(
		`INSERT INTO messages (conversation_id, role, content) VALUES (?, ?, ?)`,
		convID, "assistant",
		`*A figure steps from between the ancient trees, her silver hair catching the dappled light. She turns with amber eyes that hold centuries of forest wisdom, her gown shimmering with quiet magic.* "Welcome, traveler. You have found the last safe glade in Eldoria." *She gestures at the peaceful clearing around you.* "I am Seraphina, guardian of this forest. Few find their way here — though I sense you needed to."`,
	)
	return err
}
