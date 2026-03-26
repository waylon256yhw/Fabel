package api

import (
	"database/sql"

	"fabel/internal/dbq"
	"fabel/llm"

	"github.com/alexedwards/scs/v2"
)

type Server struct {
	q   *dbq.Queries
	db  *sql.DB
	llm *llm.Client
	sm  *scs.SessionManager
}

func NewServer(sqlDB *sql.DB, llmClient *llm.Client, sm *scs.SessionManager) *Server {
	return &Server{
		q:   dbq.New(sqlDB),
		db:  sqlDB,
		llm: llmClient,
		sm:  sm,
	}
}
