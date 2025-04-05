package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type RPCClient struct {
	Client   *rpc.Client
	Wallet   solana.PublicKey
	Filter   Filter
	TestMode bool // Add test mode flag
}

func NewRPCClient(endpoint, wallet string, filter Filter, testMode bool) (*RPCClient, error) {
	client := rpc.New(endpoint)
	pubKey, err := solana.PublicKeyFromBase58(wallet)
	if err != nil {
		return nil, err
	}
	return &RPCClient{Client: client, Wallet: pubKey, Filter: filter, TestMode: testMode}, nil
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
			if r.TestMode {
				// Simulate a network check in test mode
				_, err := r.Client.GetVersion(ctx)
				if err != nil {
					out <- fmt.Sprintf("Test mode: RPC connection failed: %v", err)
				} else {
					out <- "Test mode: RPC connection is live"
				}
				continue
			}
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
				tx, err := r.Client.GetTransaction(ctx, sig.Signature, &rpc.GetTransactionOpts{
					Encoding: solana.EncodingBase64,
				})
				if err != nil {
					fmt.Printf("Error fetching tx %s: %v\n", sig.Signature, err)
					continue
				}
				if tx == nil || tx.Meta == nil {
					continue
				}
				if r.applyFilters(tx) {
					out <- fmt.Sprintf("Tx Signature (RPC): %s", sig.Signature)
				}
			}
		}
	}
}

func (r *RPCClient) applyFilters(tx *rpc.GetTransactionResult) bool {
	if tx == nil || tx.Meta == nil || tx.Transaction == nil {
		return false
	}

	// Decode the transaction using GetTransaction()
	decodedTx, err := tx.Transaction.GetTransaction()
	if err != nil {
		fmt.Printf("Error decoding transaction: %v\n", err)
		return false
	}
	if decodedTx == nil {
		return false
	}

	// Check SOL amount
	if r.Filter.MinSOL > 0 {
		preBalance := float64(tx.Meta.PreBalances[0]) / 1e9 // Lamports to SOL
		postBalance := float64(tx.Meta.PostBalances[0]) / 1e9
		amount := preBalance - postBalance
		if amount <= 0 || amount < r.Filter.MinSOL {
			return false
		}
	}

	// Check for token transfers
	if r.Filter.TokenProgramID != "" {
		tokenProgramID, _ := solana.PublicKeyFromBase58(r.Filter.TokenProgramID)
		for _, instruction := range decodedTx.Message.Instructions {
			programID := decodedTx.Message.AccountKeys[instruction.ProgramIDIndex]
			if programID.Equals(tokenProgramID) {
				return true // Match found
			}
		}
		return false // No token transfer found
	}

	return true // No filters applied or passed
}
