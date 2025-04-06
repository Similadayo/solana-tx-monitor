package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/similadayo/solana-tx-monitor/pkg/monitor"
)

type Server struct {
	monitor *monitor.Monitor
	logs    []string
}

func StartServer(port string) {
	s := &Server{logs: make([]string, 0)}
	router := mux.NewRouter()
	router.HandleFunc("/start", s.handleStart).Methods("POST")
	router.HandleFunc("/stop", s.handleStop).Methods("POST")
	router.HandleFunc("/transactions", s.handleTransactions).Methods("GET")

	fmt.Printf("Starting API server on :%s\n", port)
	http.ListenAndServe(":"+port, router)
}

func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	var cfg monitor.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if cfg.Wallet == "" {
		http.Error(w, "Wallet is required", http.StatusBadRequest)
		return
	}
	if cfg.RPCEndpoint == "" {
		cfg.RPCEndpoint = "https://api.devnet.solana.com"
	}
	if cfg.WebSocketEndpoint == "" {
		cfg.WebSocketEndpoint = "wss://api.devnet.solana.com"
	}
	if cfg.OutputType == "" {
		cfg.OutputType = "console"
	}

	if s.monitor != nil {
		http.Error(w, "Monitoring already in progress", http.StatusConflict)
		return
	}

	mon, err := monitor.NewMonitor(cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start monitor: %v", err), http.StatusInternalServerError)
		return
	}
	s.monitor = mon
	s.logs = []string{} // Reset logs
	go func() {
		mon.Start(func(msg string) {
			s.logs = append(s.logs, msg)
		})
		s.monitor = nil // Clear monitor when done
	}()
	w.Write([]byte("Monitoring started"))
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	if s.monitor == nil {
		http.Error(w, "No monitoring in progress", http.StatusBadRequest)
		return
	}
	s.monitor.Stop()
	s.monitor = nil
	w.Write([]byte("Monitoring stopped"))
}

func (s *Server) handleTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.logs)
}
