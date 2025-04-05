package monitor

type Config struct {
	RPCEndpoint       string
	WebSocketEndpoint string
	Wallet            string
}

type Monitor struct {
	Config Config
	RPC    *RPCClient
	WS     *WebSocketClient
}
