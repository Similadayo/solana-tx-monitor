# Solana Transaction Monitor

A lightweight tool to monitor Solana wallet transactions in real-time, built in Go—my first Web3 backend project.

## Features

- Tracks transactions via RPC polling and WebSocket.
- CLI: `solana-tx-monitor monitor --wallet <pubkey>` or `--config config.yaml`.
- API: `solana-tx-monitor server` with `/start`, `/stop`, `/transactions`.
- Filters: Token transfers, min SOL, specific token mints.
- Outputs: Console or CSV.
- Test mode: `--test` to check connections.

## Install

1. Clone: `git clone https://github.com/yourusername/solana-tx-monitor.git`
2. Build: `cd solana-tx-monitor && go build -o solana-tx-monitor`

## Usage

**CLI:**

```bash
./solana-tx-monitor monitor --config config.yaml
```

**Server:**

```bash
./solana-tx-monitor server --port 8080
curl -X POST -d '{"wallet":"YourAddress..."}' http://localhost:8080/start
curl http://localhost:8080/transactions
```

**Config Example (config.yaml):**

```yaml
wallet: "YourAddress"
output_type: "console"
```

## License

MIT License. See LICENSE for details.
