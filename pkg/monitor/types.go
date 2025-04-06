package monitor

type Config struct {
	RPCEndpoint       string
	WebSocketEndpoint string
	Wallet            string
	Filters           Filter
	TestMode          bool   // New field for test mode
	OutputType        string // New: "console" or "csv"
	OutputFile        string // New: file path for CSV
}

type Filter struct {
	TokenProgramID string
	MinSOL         float64
}

type Monitor struct {
	Config Config
	RPC    *RPCClient
	WS     *WebSocketClient
}
