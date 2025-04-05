package cmd

import (
	"github.com/similadayo/solana-tx-monitor/monitor"
	"github.com/spf13/cobra"
)

var wallet string
var tokenProgramID string
var minSOL float64

func init() {
	monitorCmd.Flags().StringVarP(&wallet, "wallet", "w", "", "Wallet public key to monitor (required)")
	monitorCmd.Flags().StringVar(&tokenProgramID, "filter-token", "", "Filter for SPL token transfers (e.g., TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA)")
	monitorCmd.Flags().Float64Var(&minSOL, "min-sol", 0, "Minimum SOL amount to filter (e.g., 1.0)")
	monitorCmd.MarkFlagRequired("wallet")
	rootCmd.AddCommand(monitorCmd)
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start monitoring transactions for a wallet",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := monitor.Config{
			RPCEndpoint:       "https://api.devnet.solana.com",
			WebSocketEndpoint: "wss://api.devnet.solana.com",
			Wallet:            wallet,
			Filter: monitor.Filter{
				TokenProgramID: tokenProgramID,
				MinSOL:         minSOL,
			},
		}
		mon, err := monitor.NewMonitor(cfg)
		if err != nil {
			panic(err)
		}
		mon.Start()
	},
}
