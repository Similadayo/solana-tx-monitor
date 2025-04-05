package monitor

type Config struct {
	RPCEndpoint       string
	WebSocketEndpoint string
	Wallet            string
	Filter            Filter
}

type Filter struct {
	TokenProgramID string  // Filter for SPL token transfers (empty = no filter)
	MinSOL         float64 // Minimum SOL amount (0 = no filter)
}

type Monitor struct {
	Config Config
	RPC    *RPCClient
	WS     *WebSocketClient
}
