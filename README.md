# nucleus-sdk-go

## Overview

The Nucleus SDK is a Go-based library designed to interact with blockchain networks. It provides tools for managing calldata queues, interacting with smart contracts, and handling Ethereum RPC clients.

## Features

- **RPC Registry**: Manage Ethereum RPC clients for multiple chains.
- **Calldata Queue**: Queue, build calldata and execute transactions with strategist.
- **Strategist Transactor**: Handle signing and transaction execution for strategists.
- **Client Integration**: Interact with the Nucleus API for address book and merkle proof data.

## Installation

To use the Nucleus SDK in your project:

```bash
go get github.com/Tempest-Finance/nucleus-sdk-go
```

## Usage
### Initializing the RPC Registry
```go
import (
	"github.com/Tempest-Finance/nucleus-sdk-go/pkg/rpcregistry"
)

config := map[int64]rpcregistry.Config{
	1: {HTTP: "url"},
}

rpcRegistry, err := rpcregistry.NewRpcRegistry(config)
if err != nil {
	log.Fatalf("Failed to initialize RPC registry: %v", err)
}
```

### Creating a Calldata Queue
```go
import (
"github.com/Tempest-Finance/nucleus-sdk-go/pkg/nucleus"
)

client, _ := nucleus.NewClient("YOUR_API_KEY", "https://api.nucleusearn.io")
rpcRegistry := // Initialize RPC registry as shown above

calldataQueue, err := nucleus.NewCalldataQueue(
  1, // Chain ID
  "0xStrategistAddress",
  "VaultSymbol",
  client,
  rpcRegistry,
)
if err != nil {
  log.Fatalf("Failed to create calldata queue: %v", err)
}

```

### Adding a Transaction to the Queue
```go
import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
)

calldataQueue.AddCall(
  common.HexToAddress("0xTargetAddress"),
  []byte("0xCalldata"),
  big.NewInt(0),
)

txHash, err := calldataQueue.Execute(context.Background(), transactors)
if err != nil {
  log.Fatalf("Failed to execute transaction: %v", err)
}

log.Printf("Transaction executed successfully: %s", txHash)
```


