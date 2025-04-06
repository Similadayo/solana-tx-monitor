package monitor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/similadayo/solana-tx-monitor/pkg/output"
)

func NewMonitor(cfg Config) (*Monitor, error) {
	rpcClient, err := NewRPCClient(cfg.RPCEndpoint, cfg.Wallet, cfg.Filters, cfg.TestMode)
	if err != nil {
		return nil, err
	}
	wsClient, err := NewWebSocketClient(cfg.WebSocketEndpoint, cfg.Wallet)
	if err != nil {
		return nil, err
	}
	return &Monitor{Config: cfg, RPC: rpcClient, WS: wsClient}, nil
}

func (m *Monitor) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var outHandler interface {
		Write(string) error
	}
	switch m.Config.OutputType {
	case "csv":
		csvOut, err := output.NewCSVOutput(m.Config.OutputFile)
		if err != nil {
			fmt.Printf("Failed to initialize CSV output: %v\n", err)
			return
		}
		defer csvOut.Close()
		outHandler = csvOut
	default: // "console" or empty
		outHandler = output.NewConsoleOutput()
	}

	out := make(chan string, 100)
	go m.RPC.Poll(ctx, out)
	go m.WS.Subscribe(ctx, out)

	fmt.Println("Monitoring started for wallet:", m.Config.Wallet)
	outHandler.Write("Monitoring started for wallet: " + m.Config.Wallet)

	inactivityTimeout := time.NewTimer(30 * time.Second)
	defer inactivityTimeout.Stop()

	for {
		select {
		case msg := <-out:
			outHandler.Write(msg)
			if strings.HasPrefix(msg, "Tx Signature (RPC):") || strings.HasPrefix(msg, "Account Update (WS):") || m.Config.TestMode {
				if !inactivityTimeout.Stop() {
					<-inactivityTimeout.C
				}
				inactivityTimeout.Reset(30 * time.Second)
			}
		case <-ctx.Done():
			outHandler.Write("Monitoring stopped")
			return
		case <-inactivityTimeout.C:
			if m.Config.TestMode {
				outHandler.Write("Test mode completed successfully")
			} else {
				outHandler.Write(fmt.Sprintf("No activity detected for wallet %s in the last 30 seconds on Devnet. Closing connections and exiting.", m.Config.Wallet))
			}
			cancel()
			time.Sleep(1 * time.Second)
			return
		}
	}
}
