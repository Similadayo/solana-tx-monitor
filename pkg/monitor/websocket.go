package monitor

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

type WebSocketClient struct {
	Client *ws.Client
	Wallet solana.PublicKey
}

func NewWebSocketClient(endpoint, wallet string) (*WebSocketClient, error) {
	client, err := ws.Connect(context.Background(), endpoint)
	if err != nil {
		return nil, err
	}
	pubKey, err := solana.PublicKeyFromBase58(wallet)
	if err != nil {
		return nil, err
	}
	return &WebSocketClient{Client: client, Wallet: pubKey}, nil
}

func (w *WebSocketClient) Subscribe(ctx context.Context, out chan<- string) {
	fmt.Println("Connecting to WebSocket...")
	sub, err := w.Client.AccountSubscribe(w.Wallet, "")
	if err != nil {
		fmt.Printf("WS subscription error: %v\n", err)
		return
	}
	defer sub.Unsubscribe()
	defer w.Client.Close()

	fmt.Printf("Subscribed to wallet: %s\n", w.Wallet.String())
	out <- "WebSocket connection is live" // Confirm connection

	for {
		select {
		case <-ctx.Done():
			fmt.Println("WS subscription stopped")
			return
		case data := <-sub.Response():
			out <- fmt.Sprintf("Account Update (WS): %v", data.Value)
		case err := <-sub.Err():
			fmt.Printf("WS error: %v\n", err)
			return
		}
	}
}
