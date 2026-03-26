package app

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"fabel/db"
	"fabel/llm"
)

type App struct {
	db  *db.DB
	llm *llm.Client
	web fs.FS
}

func New(database *db.DB, llmClient *llm.Client, webDist embed.FS) *App {
	sub, _ := fs.Sub(webDist, "dist")
	return &App{db: database, llm: llmClient, web: sub}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/bootstrap", a.handleBootstrap)
	mux.HandleFunc("POST /api/conversations", a.handleCreateConversation)
	mux.HandleFunc("GET /api/conversations/{id}", a.handleGetConversation)
	mux.HandleFunc("POST /api/conversations/{id}/send", a.handleSendMessageStream)
	mux.HandleFunc("GET /api/conversations/{id}/prompt", a.handleGetPrompt)

	// Serve built frontend; in dev mode dist only has placeholder.txt so this is a no-op.
	fileServer := http.FileServerFS(a.web)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// SPA fallback: serve index.html for non-asset paths
		if _, err := fs.Stat(a.web, strings.TrimPrefix(r.URL.Path, "/")); err != nil {
			r = r.Clone(r.Context())
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
