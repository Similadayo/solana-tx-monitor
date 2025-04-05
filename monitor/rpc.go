package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type RPCClient struct {
	Client *rpc.Client
	Wallet solana.PublicKey
}

func NewRPCClient(endpoint, wallet string) (*RPCClient, error) {
	client := rpc.New(endpoint)
	pubkey, err := solana.PublicKeyFromBase58(wallet)
	if err != nil {
		return nil, err
	}
	return &RPCClient{
		Client: client,
		Wallet: pubkey,
	}, nil
}

func (r *RPCClient) Poll(ctx context.Context, out chan<- string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("Starting RPC polling...")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("RPC polling stopped")
			return
		case <-ticker.C:
			fmt.Printf("Polling for wallet: %s\n", r.Wallet.String())
			sigs, err := r.Client.GetSignaturesForAddress(ctx, r.Wallet)
			if err != nil {
				fmt.Printf("RPC error: %v\n", err)
				continue
			}
			if len(sigs) == 0 {
				out <- fmt.Sprintf("No recent transactions found for wallet %s", r.Wallet.String())
				continue
			}
			for _, sig := range sigs {
				out <- fmt.Sprintf("Tx Signature (RPC): %s", sig.Signature)
			}
		}
	}
}
