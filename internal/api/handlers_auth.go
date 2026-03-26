package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"fabel/internal/auth"
	"fabel/internal/dbq"
	"fabel/internal/resource"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr(w, err, 400)
		return
	}

	// Validate username.
	username := strings.TrimSpace(req.Username)
	if utf8.RuneCountInString(username) < 2 || utf8.RuneCountInString(username) > 32 {
		http.Error(w, "username must be 2-32 characters", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, "password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Determine role: first user becomes admin.
	count, err := s.q.CountUsers(ctx)
	if err != nil {
		httpErr(w, err, 500)
		return
	}

	role := "user"
	if count == 0 {
		role = "admin"
	} else {
		// Check if registration is allowed.
		val, err := s.q.GetSetting(ctx, "allow_registration")
		if err != nil || val != "true" {
			http.Error(w, "registration is disabled", http.StatusForbidden)
			return
		}
	}

	// Check username uniqueness.
	_, err = s.q.GetUserByUsername(ctx, username)
	if err == nil {
		http.Error(w, "username already taken", http.StatusBadRequest)
		return
	}
	if err != sql.ErrNoRows {
		httpErr(w, err, 500)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		httpErr(w, err, 500)
		return
	}

	displayName := ""
	if req.DisplayName != nil {
		displayName = *req.DisplayName
	}

	userID := auth.NewID()
	if err := s.q.CreateUser(ctx, dbq.CreateUserParams{
		ID:           userID,
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: hash,
		Role:         role,
	}); err != nil {
		if isUniqueViolation(err) {
			http.Error(w, "username already taken", http.StatusBadRequest)
			return
		}
		httpErr(w, err, 500)
		return
	}

	// Backfill orphaned seed data to the first user and index into resource layer.
	if role == "admin" {
		uid := userIDParam(userID)
		_ = s.q.BackfillCharacters(ctx, uid)
		_ = s.q.BackfillPresets(ctx, uid)
		_ = s.q.BackfillConversations(ctx, uid)

		// Index backfilled seed data.
		if chars, err := s.q.ListCharacters(ctx, uid); err == nil {
			for _, c := range chars {
				_ = s.idx.Upsert(ctx, resource.Document{
					ID:      fmt.Sprintf("character:%d", c.ID),
					UserID:  userID,
					Kind:    resource.KindCharacter,
					Title:   c.Name,
					Content: c.Description + "\n" + c.Personality + "\n" + c.Scenario,
				})
			}
		}
		if presets, err := s.q.ListPresets(ctx, uid); err == nil {
			for _, p := range presets {
				_ = s.idx.Upsert(ctx, resource.Document{
					ID:      fmt.Sprintf("preset:%d", p.ID),
					UserID:  userID,
					Kind:    resource.KindPreset,
					Title:   p.Name,
					Content: p.SystemPrompt,
				})
			}
		}
	}

	// Auto-login: renew token to prevent session fixation, then set user_id.
	if err := s.sm.RenewToken(ctx); err != nil {
		httpErr(w, err, 500)
		return
	}
	s.sm.Put(ctx, "user_id", userID)

	user, err := s.q.GetUserByID(ctx, userID)
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	writeJSON(w, 201, dbUserToAPI(user))
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr(w, err, 400)
		return
	}

	ctx := r.Context()

	user, err := s.q.GetUserByUsername(ctx, req.Username)
	if err != nil {
		// Generic error message: don't reveal whether username exists.
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	// Renew session token to prevent session fixation.
	if err := s.sm.RenewToken(ctx); err != nil {
		httpErr(w, err, 500)
		return
	}
	s.sm.Put(ctx, "user_id", user.ID)

	writeJSON(w, 200, dbUserToAPI(user))
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	if err := s.sm.Destroy(r.Context()); err != nil {
		httpErr(w, err, 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	writeJSON(w, 200, dbUserToAPI(*user))
}
