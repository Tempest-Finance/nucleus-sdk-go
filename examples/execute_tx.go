package examples

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Tempest-Finance/nucleus-sdk-go/pkg/nucleus"
	"github.com/Tempest-Finance/nucleus-sdk-go/pkg/rpcregistry"
	"github.com/ethereum/go-ethereum/common"
)

func ExampleExecuteTx(ctx context.Context, apiKey, baseUrl string, chainId int64, strategistAddress, symbol string, rpcRegistry rpcregistry.IRegistry, target string, calldata []byte, transactors nucleus.IStrategistTransactor) {
	client, err := nucleus.NewClient(apiKey, baseUrl)
	if err != nil {
		panic(err)
	}

	calldataQueue, err := nucleus.NewCalldataQueue(chainId, strategistAddress, symbol, client, rpcRegistry)
	if err != nil {
		panic(err)
	}

	calldataQueue.AddCall(common.HexToAddress(target), calldata, big.NewInt(0))

	txHash, err := calldataQueue.Execute(ctx, transactors)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction hash: %s\n", txHash)
}
