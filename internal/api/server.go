package api

import (
	"database/sql"

	"fabel/internal/dbq"
	"fabel/llm"
)

type Server struct {
	q   *dbq.Queries
	db  *sql.DB
	llm *llm.Client
}

func NewServer(sqlDB *sql.DB, llmClient *llm.Client) *Server {
	return &Server{
		q:   dbq.New(sqlDB),
		db:  sqlDB,
		llm: llmClient,
	}
}
