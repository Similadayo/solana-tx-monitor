package cmd

import (
	"github.com/similadayo/solana-tx-monitor/pkg/api"
	"github.com/spf13/cobra"
)

var port string

func init() {
	serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the API server on")
	rootCmd.AddCommand(serverCmd) // This should register the command
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the transaction monitor as an HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		api.StartServer(port)
	},
}
