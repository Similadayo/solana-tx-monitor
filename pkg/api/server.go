package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/similadayo/solana-tx-monitor/pkg/monitor"
)

type ServerConfig struct {
	RPCEndpoint       string         `yaml:"rpc_endpoint"`
	WebSocketEndpoint string         `yaml:"websocket_endpoint"`
	Filters           monitor.Filter `yaml:"filters"`
	OutputType        string         `yaml:"output_type"`
	OutputFile        string         `yaml:"output_file"`
	TestMode          bool           `yaml:"test_mode"`
}

var defaultConfig = ServerConfig{
	RPCEndpoint:       "https://api.devnet.solana.com",
	WebSocketEndpoint: "wss://api.devnet.solana.com",
	OutputType:        "console",
	OutputFile:        "transactions.csv",
}

func SetDefaultConfig(cfg ServerConfig) {
	defaultConfig = cfg
}

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
	// Apply defaults from config file
	if cfg.RPCEndpoint == "" {
		cfg.RPCEndpoint = defaultConfig.RPCEndpoint
	}
	if cfg.WebSocketEndpoint == "" {
		cfg.WebSocketEndpoint = defaultConfig.WebSocketEndpoint
	}
	if cfg.OutputType == "" {
		cfg.OutputType = defaultConfig.OutputType
	}
	if cfg.OutputFile == "" {
		cfg.OutputFile = defaultConfig.OutputFile
	}
	if cfg.Filters.TokenProgramID == "" && defaultConfig.Filters.TokenProgramID != "" {
		cfg.Filters.TokenProgramID = defaultConfig.Filters.TokenProgramID
	}
	if cfg.Filters.MinSOL == 0 && defaultConfig.Filters.MinSOL != 0 {
		cfg.Filters.MinSOL = defaultConfig.Filters.MinSOL
	}
	cfg.TestMode = cfg.TestMode || defaultConfig.TestMode

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
	s.logs = []string{}
	go func() {
		mon.Start(func(msg string) {
			s.logs = append(s.logs, msg)
		})
		s.monitor = nil
	}()
	w.Write([]byte("Monitoring started"))
}

// handleStop and handleTransactions remain unchanged

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
