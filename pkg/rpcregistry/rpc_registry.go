package rpcregistry

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

// RpcRegistry is a collection of clients for different chains
type RpcRegistry struct {
	clientByChainID map[int64]*ethclient.Client
}

func NewRpcRegistry(config Config) (*RpcRegistry, error) {
	clientByChainID := make(map[int64]*ethclient.Client)

	for chainID, cfg := range config {
		client, err := ethclient.Dial(cfg.HTTP)
		if err != nil {
			return nil, fmt.Errorf("failed to dial RPC for chainID %d: %w", chainID, err)
		}

		clientByChainID[chainID] = client
	}

	return &RpcRegistry{
		clientByChainID: clientByChainID,
	}, nil
}

func (h *RpcRegistry) GetClient(chainID int64) (*ethclient.Client, error) {
	client, ok := h.clientByChainID[chainID]
	if !ok {
		return nil, fmt.Errorf("no client found for chainID %d", chainID)
	}

	return client, nil
}
