package auth

import (
	"context"
	"database/sql"
	"time"

	"fabel/internal/dbq"
)

// SQLiteStore implements the alexedwards/scs Store interface backed by sqlc queries.
type SQLiteStore struct {
	q  *dbq.Queries
	db *sql.DB
}

// NewSQLiteStore creates a new session store.
func NewSQLiteStore(sqlDB *sql.DB) *SQLiteStore {
	return &SQLiteStore{
		q:  dbq.New(sqlDB),
		db: sqlDB,
	}
}

func (s *SQLiteStore) Find(token string) ([]byte, bool, error) {
	row, err := s.q.GetSession(context.Background(), token)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if time.Now().After(time.Unix(int64(row.Expiry), 0)) {
		return nil, false, nil
	}
	return row.Data, true, nil
}

func (s *SQLiteStore) Commit(token string, data []byte, expiry time.Time) error {
	return s.q.UpsertSession(context.Background(), dbq.UpsertSessionParams{
		Token:  token,
		Data:   data,
		Expiry: float64(expiry.Unix()),
	})
}

func (s *SQLiteStore) Delete(token string) error {
	return s.q.DeleteSession(context.Background(), token)
}

func (s *SQLiteStore) All() (map[string][]byte, error) {
	rows, err := s.q.GetAllSessions(context.Background(), float64(time.Now().Unix()))
	if err != nil {
		return nil, err
	}
	m := make(map[string][]byte, len(rows))
	for _, r := range rows {
		m[r.Token] = r.Data
	}
	return m, nil
}

// StartCleanup periodically removes expired sessions.
func (s *SQLiteStore) StartCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			_ = s.q.DeleteExpiredSessions(context.Background(), float64(time.Now().Unix()))
		}
	}()
}
