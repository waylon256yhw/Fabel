package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"fabel/internal/auth"
	"fabel/internal/dbq"
)

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := s.q.ListUsers(r.Context())
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	users := make([]User, len(rows))
	for i, row := range rows {
		users[i] = User{
			Id:          row.ID,
			Username:    row.Username,
			DisplayName: row.DisplayName,
			Role:        UserRole(row.Role),
			CreatedAt:   row.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	writeJSON(w, 200, UserList{Users: users})
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr(w, err, 400)
		return
	}

	ctx := r.Context()
	currentUser := auth.UserFromContext(ctx)

	// Check target user exists.
	target, err := s.q.GetUserByID(ctx, id)
	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		httpErr(w, err, 500)
		return
	}

	if req.Role != nil {
		newRole := string(*req.Role)
		// Prevent admin from demoting themselves.
		if currentUser.ID == target.ID && newRole != "admin" {
			http.Error(w, "cannot demote yourself", http.StatusBadRequest)
			return
		}
		if err := s.q.UpdateUserRole(ctx, dbq.UpdateUserRoleParams{
			Role: newRole,
			ID:   id,
		}); err != nil {
			httpErr(w, err, 500)
			return
		}
	}

	// Re-fetch updated user.
	updated, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	writeJSON(w, 200, dbUserToAPI(updated))
}

func (s *Server) GetSettings(w http.ResponseWriter, r *http.Request) {
	rows, err := s.q.ListSettings(r.Context())
	if err != nil {
		httpErr(w, err, 500)
		return
	}
	settings := ServerSettings{}
	for _, row := range rows {
		if row.Key == "allow_registration" {
			v := row.Value == "true"
			settings.AllowRegistration = &v
		}
	}
	writeJSON(w, 200, settings)
}

func (s *Server) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req ServerSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpErr(w, err, 400)
		return
	}

	ctx := r.Context()
	if req.AllowRegistration != nil {
		val := "false"
		if *req.AllowRegistration {
			val = "true"
		}
		if err := s.q.UpsertSetting(ctx, dbq.UpsertSettingParams{
			Key:   "allow_registration",
			Value: val,
		}); err != nil {
			httpErr(w, err, 500)
			return
		}
	}

	// Return current settings.
	s.GetSettings(w, r)
}
