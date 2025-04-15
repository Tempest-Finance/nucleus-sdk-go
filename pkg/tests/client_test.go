package tests

import (
	"context"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/Tempest-Finance/nucleus-sdk-go/pkg/nucleus"
	"github.com/Tempest-Finance/nucleus-sdk-go/pkg/rpcregistry"
	"github.com/Tempest-Finance/nucleus-sdk-go/pkg/tests/abis/erc20"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	zero = big.NewInt(0)
)

func TestClient(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err, "failed to load .env file")
	baseUrl := os.Getenv("BASE_URL")
	nucleusApiKey := os.Getenv("API_KEY")
	rpcUrl := os.Getenv("RPC_URL")
	strategistAddress := os.Getenv("STRATEGIST_ADDRESS")
	invalidStrategistAddress := os.Getenv("INVALID_STRATEGIST_ADDRESS")
	symbol := os.Getenv("SYMBOL")
	erc20Address := common.HexToAddress(os.Getenv("ERC20_ADDRESS"))
	chainIdInt, err := strconv.Atoi(os.Getenv("CHAIN_ID"))
	chainId := int64(chainIdInt)
	assert.NoError(t, err, "chainId should be a number")
	rpcRegistry, err := rpcregistry.NewRpcRegistry(map[int64]rpcregistry.ChainConfig{
		chainId: {
			HTTP: rpcUrl,
		},
	})
	assert.NoError(t, err, "failed to create RPC registry")

	strategistConfig := []nucleus.StrategistSignerConfig{
		{
			ChainId: chainId,
			Name:    "STRATEGIST_PRIVATE_KEY",
			Address: strategistAddress,
		},
	}

	invalidStrategistConfig := []nucleus.StrategistSignerConfig{
		{
			ChainId: chainId,
			Name:    "INVALID_STRATEGIST_PRIVATE_KEY",
			Address: strategistAddress,
		},
	}

	transactors, err := nucleus.NewStrategist(strategistConfig)
	assert.NoError(t, err, "failed to create strategist transactor")

	invalidTransactors, err := nucleus.NewStrategist(invalidStrategistConfig)
	assert.NoError(t, err, "failed to create invalid strategist transactor")

	client, err := nucleus.NewClient(nucleusApiKey, baseUrl)
	assert.NoError(t, err, "failed to create client")

	erc20Abi, err := erc20.Erc20MetaData.GetAbi()
	assert.NoError(t, err, "failed to get erc20 abi")

	ctx := context.Background()

	router := client.GetAddressBook()[chainId].Hyperswap.Router

	calldata, err := erc20Abi.Pack("approve", common.HexToAddress(router), big.NewInt(1000000000000000000))
	assert.NoError(t, err, "failed to get approve calldata")

	t.Run("get valid calldata queue", func(t *testing.T) {
		_, err = nucleus.NewCalldataQueue(chainId, strategistAddress, symbol, client, rpcRegistry)
		assert.NoError(t, err, "failed to create calldata queue")
	})

	t.Run("test execute calldata queue", func(t *testing.T) {
		calldataQueue, err := nucleus.NewCalldataQueue(chainId, strategistAddress, symbol, client, rpcRegistry)
		assert.NoError(t, err, "failed to create calldata queue")

		router := client.GetAddressBook()[chainId].Hyperswap.Router

		calldata, err := erc20Abi.Pack("approve", common.HexToAddress(router), big.NewInt(1000000000000000000))
		assert.NoError(t, err, "failed to get approve calldata")

		calldataQueue.AddCall(erc20Address, calldata, zero)

		txHash, err := calldataQueue.Execute(ctx, transactors)
		assert.NoError(t, err, "failed to execute calldata queue")

		t.Logf("txHash: %s", txHash)
	})

	t.Run("test execute with invalid transactors", func(t *testing.T) {
		calldataQueue, err := nucleus.NewCalldataQueue(chainId, strategistAddress, symbol, client, rpcRegistry)
		assert.NoError(t, err, "failed to create calldata queue")

		calldataQueue.AddCall(erc20Address, calldata, zero)

		_, err = calldataQueue.Execute(ctx, invalidTransactors)
		assert.ErrorIs(t, err, nucleus.ErrInvalidSigner, "should return invalid signer error")
	})

	t.Run("test invalid nucleus api key", func(t *testing.T) {
		invalidClient, err := nucleus.NewClient("invalid_api_key", baseUrl)
		assert.NoError(t, err, "failed to create client")
		calldataQueue, err := nucleus.NewCalldataQueue(chainId, strategistAddress, symbol, invalidClient, rpcRegistry)
		assert.NoError(t, err, "failed to create calldata queue")

		calldataQueue.AddCall(erc20Address, calldata, zero)

		_, err = calldataQueue.GetCalldata(ctx)
		assert.Error(t, err, "should return err cuz invalid nucleus api key error")
	})

	t.Run("test invalid strategist", func(t *testing.T) {
		_, err = nucleus.NewCalldataQueue(chainId, invalidStrategistAddress, symbol, client, rpcRegistry)
		assert.ErrorIs(t, err, nucleus.ErrStrategiesIsInvalid, "should return invalid strategist address error")
	})

	t.Run("test invalid chain id", func(t *testing.T) {
		_, err = nucleus.NewCalldataQueue(1, invalidStrategistAddress, symbol, client, rpcRegistry)
		assert.Error(t, err, "should return invalid chain id error")
	})

	t.Run("test invalid symbol", func(t *testing.T) {
		invalidSymbol := "invalid_symbol"
		_, err = nucleus.NewCalldataQueue(1, invalidStrategistAddress, invalidSymbol, client, rpcRegistry)
		assert.Error(t, err, "should return invalid symbol error")
	})

	t.Run("test invalid target", func(t *testing.T) {
		calldataQueue, err := nucleus.NewCalldataQueue(chainId, strategistAddress, symbol, client, rpcRegistry)
		assert.NoError(t, err, "failed to create calldata queue")

		calldataQueue.AddCall(common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), calldata, zero)

		_, err = calldataQueue.GetCalldata(ctx)
		assert.Error(t, err, "should return invalid target error")
	})

	t.Run("test get calldata and target", func(t *testing.T) {
		ethClient, err := rpcRegistry.GetClient(chainId)
		assert.NoError(t, err, "failed to get eth client")

		calldataQueue, err := nucleus.NewCalldataQueue(chainId, strategistAddress, symbol, client, rpcRegistry)
		assert.NoError(t, err, "failed to create calldata queue")

		calldataQueue.AddCall(erc20Address, calldata, zero)

		calldataBytes, target, err := calldataQueue.GetCalldataBytesAndTarget(ctx)
		assert.NoError(t, err, "failed to get target and calldata")
		assert.NotEmpty(t, calldata, "targets should not be empty")
		assert.NotEmpty(t, target, "targetData should not be empty")

		strategist, err := transactors.GetStrategist(strategistAddress)
		assert.NoError(t, err, "failed to get strategist")

		nonce, err := ethClient.PendingNonceAt(ctx, strategist.From)
		assert.NoError(t, err, "failed to get nonce")

		gasPrice, err := ethClient.SuggestGasPrice(ctx)
		assert.NoError(t, err, "failed to get gas price")

		gasLimit := uint64(300_000)

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       target,
			Value:    zero,
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     calldataBytes,
		})

		signedTx, err := strategist.Signer(strategist.From, tx)
		assert.NoError(t, err, "failed to sign transaction")

		err = ethClient.SendTransaction(ctx, signedTx)
		assert.NoError(t, err, "failed to send transaction")

		t.Logf("txHash: %s", signedTx.Hash().Hex())
	})
}
