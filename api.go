package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type HTTPServer struct {
	addr string
	db   *sql.DB
	svc  *MonzoClient // TODO: define an interface for this so we're not tied to MonzoClient
}

func NewHTTPServer(addr string, db *sql.DB, svc *MonzoClient) *HTTPServer {
	return &HTTPServer{
		addr: addr,
		db:   db,
		svc:  svc,
	}
}

func (s *HTTPServer) Start() error {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		transactions, err := s.svc.ListTransactions(ctx, r.URL.Query().Get("account_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, transactions)
	})

	return http.ListenAndServe(s.addr, handler)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
