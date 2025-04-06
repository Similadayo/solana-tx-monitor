package cmd

import (
	"io/ioutil"

	"github.com/similadayo/solana-tx-monitor/pkg/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var port string
var serverConfigFile string

func init() {
	serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the API server on")
	serverCmd.Flags().StringVar(&serverConfigFile, "config", "", "Path to config file (e.g., config.yaml)")
	rootCmd.AddCommand(serverCmd) // Ensure this is present
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the transaction monitor as an HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		if serverConfigFile != "" {
			data, err := ioutil.ReadFile(serverConfigFile)
			if err != nil {
				panic(err)
			}
			var cfg api.ServerConfig
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				panic(err)
			}
			api.SetDefaultConfig(cfg) // Pass config to API
		}
		api.StartServer(port)
	},
}
