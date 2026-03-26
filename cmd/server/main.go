package main

import (
	"database/sql"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	fabeldb "fabel/db"
	"fabel/internal/api"
	"fabel/llm"
	"fabel/web"
)

func main() {
	// Open database.
	sqlDB, err := sql.Open("sqlite", "fabel.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	sqlDB.SetMaxOpenConns(1)

	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000",
	} {
		if _, err := sqlDB.Exec(pragma); err != nil {
			log.Fatal(err)
		}
	}

	// Run goose migrations (embedded).
	goose.SetBaseFS(fabeldb.Migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatal(err)
	}

	// Create server.
	llmClient := llm.New(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_BASE_URL"))
	srv := api.NewServer(sqlDB, llmClient)

	// Build router.
	r := chi.NewRouter()
	r.Use(cors)

	// SSE endpoint registered manually (not in OpenAPI spec).
	r.Post("/api/conversations/{id}/send", srv.HandleSendMessage)

	// Generated routes from OpenAPI spec.
	api.HandlerFromMux(srv, r)

	// SPA static file serving.
	sub, _ := fs.Sub(web.Dist, "dist")
	fileServer := http.FileServerFS(sub)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// Unknown API routes get a proper 404, not index.html.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		// SPA fallback: serve index.html for non-asset paths.
		if _, err := fs.Stat(sub, strings.TrimPrefix(r.URL.Path, "/")); err != nil {
			r = r.Clone(r.Context())
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
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
