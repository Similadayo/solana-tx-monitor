package cmd

import (
	"io/ioutil"

	"github.com/similadayo/solana-tx-monitor/pkg/monitor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var wallet string
var tokenProgramID string
var minSOL float64
var testMode bool
var outputType string
var configFile string

func init() {
	monitorCmd.Flags().StringVar(&configFile, "config", "", "Path to config file (e.g., config.yaml)")
	monitorCmd.Flags().StringVarP(&wallet, "wallet", "w", "", "Wallet public key to monitor (overrides config)")
	monitorCmd.Flags().StringVar(&tokenProgramID, "filter-token", "", "Filter for SPL token transfers (overrides config)")
	monitorCmd.Flags().Float64Var(&minSOL, "min-sol", 0, "Minimum SOL amount to filter (overrides config)")
	monitorCmd.Flags().BoolVar(&testMode, "test", false, "Run in test mode (overrides config)")
	monitorCmd.Flags().StringVar(&outputType, "output", "console", "Output type: 'console' or 'csv' (overrides config)")
	// Removed: monitorCmd.MarkFlagRequired("wallet")
	rootCmd.AddCommand(monitorCmd)
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start monitoring transactions for a wallet",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := monitor.Config{
			RPCEndpoint:       "https://api.devnet.solana.com",
			WebSocketEndpoint: "wss://api.devnet.solana.com",
			OutputType:        "console",
			OutputFile:        "transactions.csv",
		}

		// Load config file if provided
		if configFile != "" {
			data, err := ioutil.ReadFile(configFile)
			if err != nil {
				panic(err)
			}
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				panic(err)
			}
		}

		// Override with CLI flags if set
		if wallet != "" {
			cfg.Wallet = wallet
		}
		if tokenProgramID != "" {
			cfg.Filters.TokenProgramID = tokenProgramID
		}
		if minSOL != 0 {
			cfg.Filters.MinSOL = minSOL
		}
		if testMode {
			cfg.TestMode = true
		}
		if outputType != "" {
			cfg.OutputType = outputType
		}

		// Validate wallet is set (either via config or flag)
		if cfg.Wallet == "" {
			panic("Wallet is required (set via --wallet or config file)")
		}

		mon, err := monitor.NewMonitor(cfg)
		if err != nil {
			panic(err)
		}
		mon.Start()
	},
}
