package api

import (
	"database/sql"

	"fabel/internal/dbq"
	"fabel/internal/resource"
	"fabel/llm"

	"github.com/alexedwards/scs/v2"
)

type Server struct {
	q   *dbq.Queries
	db  *sql.DB
	llm *llm.Client
	sm  *scs.SessionManager
	idx *resource.SQLiteIndexer
}

func NewServer(sqlDB *sql.DB, llmClient *llm.Client, sm *scs.SessionManager) *Server {
	q := dbq.New(sqlDB)
	return &Server{
		q:   q,
		db:  sqlDB,
		llm: llmClient,
		sm:  sm,
		idx: resource.NewIndexer(q),
	}
}
