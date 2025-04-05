package monitor

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func NewMonitor(cfg Config) (*Monitor, error) {
	rpcClient, err := NewRPCClient(cfg.RPCEndpoint, cfg.Wallet, cfg.Filter) // Pass filters
	if err != nil {
		return nil, err
	}
	wsClient, err := NewWebSocketClient(cfg.WebSocketEndpoint, cfg.Wallet)
	if err != nil {
		return nil, err
	}
	return &Monitor{Config: cfg, RPC: rpcClient, WS: wsClient}, nil
}

// Rest of the file remains unchanged

func (m *Monitor) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := make(chan string, 100)
	go m.RPC.Poll(ctx, out)
	go m.WS.Subscribe(ctx, out)

	fmt.Println("Monitoring started for wallet:", m.Config.Wallet)

	inactivityTimeout := time.NewTimer(30 * time.Second)
	defer inactivityTimeout.Stop()

	for {
		select {
		case msg := <-out:
			fmt.Println(msg)
			// Only reset timer on actual activity (tx signatures or WS updates)
			if strings.HasPrefix(msg, "Tx Signature (RPC):") || strings.HasPrefix(msg, "Account Update (WS):") {
				if !inactivityTimeout.Stop() {
					<-inactivityTimeout.C
				}
				inactivityTimeout.Reset(30 * time.Second)
			}
		case <-ctx.Done():
			fmt.Println("Monitoring stopped")
			return
		case <-inactivityTimeout.C:
			fmt.Printf("No activity detected for wallet %s in the last 30 seconds on Devnet. Closing connections and exiting.\n", m.Config.Wallet)
			cancel()
			time.Sleep(1 * time.Second) // Allow goroutines to finish
			return
		}
	}
}
