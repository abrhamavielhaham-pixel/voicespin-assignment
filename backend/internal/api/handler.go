package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"voicespin/backend/internal/conversation"
)

type Server struct {
	store *conversation.Store
	log   *slog.Logger
}

func NewServer(store *conversation.Store) *Server {
	return &Server{store: store, log: slog.Default()}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/conversations", s.listConversations)
	mux.HandleFunc("GET /api/conversations/{id}", s.getConversation)
	mux.HandleFunc("PATCH /api/conversations/{id}", s.patchConversation)
	return withCORS(mux)
}

func (s *Server) listConversations(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilter(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, s.store.List(filter))
}

func (s *Server) getConversation(w http.ResponseWriter, r *http.Request) {
	c, err := s.store.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "conversation not found")
		return
	}

	writeJSON(w, http.StatusOK, c)
}

func (s *Server) patchConversation(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status   *string `json:"status"`
		Priority *string `json:"priority"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if body.Status == nil && body.Priority == nil {
		writeError(w, http.StatusBadRequest, "at least one of status or priority is required")
		return
	}

	var update conversation.Update
	if body.Status != nil {
		status := conversation.Status(*body.Status)
		if !status.Valid() {
			writeError(w, http.StatusBadRequest, "invalid status value")
			return
		}
		update.Status = &status
	}
	if body.Priority != nil {
		priority := conversation.Priority(*body.Priority)
		if !priority.Valid() {
			writeError(w, http.StatusBadRequest, "invalid priority value")
			return
		}
		update.Priority = &priority
	}

	updated, err := s.store.Update(r.PathValue("id"), update)
	if err != nil {
		writeError(w, http.StatusNotFound, "conversation not found")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func parseFilter(r *http.Request) (conversation.Filter, error) {
	q := r.URL.Query()
	var filter conversation.Filter

	if raw := strings.TrimSpace(q.Get("status")); raw != "" {
		status := conversation.Status(strings.ToUpper(raw))
		if !status.Valid() {
			return filter, errors.New("invalid status filter")
		}
		filter.Status = status
	}
	if raw := strings.TrimSpace(q.Get("priority")); raw != "" {
		priority := conversation.Priority(strings.ToUpper(raw))
		if !priority.Valid() {
			return filter, errors.New("invalid priority filter")
		}
		filter.Priority = priority
	}
	filter.Search = strings.TrimSpace(q.Get("search"))

	return filter, nil
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encoding response", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
