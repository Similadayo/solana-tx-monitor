package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "solana-tx-monitor",
	Short: "Real-time Solana transaction monitoring",
	Long:  "A powerful tool to monitor Solana wallets and programs with RPC, WebSocket, and extensible outputs.",
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
