package cmd

import (
	"github.com/similadayo/solana-tx-monitor/monitor"
	"github.com/spf13/cobra"
)

var wallet string

func init() {
	monitorCmd.Flags().StringVarP(&wallet, "wallet", "w", "", "Wallet address to monitor")
	monitorCmd.MarkFlagRequired("wallet")
	rootCmd.AddCommand(monitorCmd)
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor Solana transactions",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := monitor.Config{
			RPCEndpoint:       "https://api.devnet.solana.com",
			WebSocketEndpoint: "wss://api.devnet.solana.com",
			Wallet:            wallet,
		}
		mon, err := monitor.NewMonitor(cfg)
		if err != nil {
			panic(err)
		}
		mon.Start()
	},
}
