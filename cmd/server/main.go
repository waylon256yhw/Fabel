package main

import (
	"database/sql"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	fabeldb "fabel/db"
	"fabel/internal/api"
	"fabel/internal/auth"
	"fabel/internal/dbq"
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

	// Session manager.
	store := auth.NewSQLiteStore(sqlDB)
	store.StartCleanup(30 * time.Minute)

	sm := scs.New()
	sm.Store = store
	sm.Lifetime = 7 * 24 * time.Hour
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode

	// Create server.
	llmClient := llm.New(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_BASE_URL"))
	srv := api.NewServer(sqlDB, llmClient, sm)

	// Wrapper for OpenAPI parameter extraction.
	wrapper := api.ServerInterfaceWrapper{
		Handler: srv,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	q := dbq.New(sqlDB)
	requireAuth := auth.RequireAuth(sm, q)

	// Build router.
	r := chi.NewRouter()
	r.Use(sm.LoadAndSave)

	// Public auth routes.
	r.Post("/api/auth/register", wrapper.Register)
	r.Post("/api/auth/login", wrapper.Login)

	// Authenticated routes.
	r.Group(func(r chi.Router) {
		r.Use(requireAuth)

		r.Post("/api/auth/logout", wrapper.Logout)
		r.Get("/api/auth/me", wrapper.GetMe)

		r.Get("/api/bootstrap", wrapper.GetBootstrap)
		r.Post("/api/conversations", wrapper.CreateConversation)
		r.Get("/api/conversations/{id}", wrapper.GetConversation)
		r.Get("/api/conversations/{id}/prompt", wrapper.GetConversationPrompt)
		r.Post("/api/conversations/{id}/send", srv.HandleSendMessage)
	})

	// Admin routes.
	r.Group(func(r chi.Router) {
		r.Use(requireAuth)
		r.Use(auth.RequireAdmin)

		r.Get("/api/admin/users", wrapper.ListUsers)
		r.Patch("/api/admin/users/{id}", wrapper.UpdateUser)
		r.Get("/api/admin/settings", wrapper.GetSettings)
		r.Patch("/api/admin/settings", wrapper.UpdateSettings)
	})

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
