package rpcregistry

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

// IRegistry is a collection of clients for different chains
type IRegistry interface {
	GetClient(chainID int64) (*ethclient.Client, error)
}
