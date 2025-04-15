package nucleus

import (
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

type StrategistTransactor struct {
	strategistTransactors map[string]*bind.TransactOpts
}

func NewStrategist(configs []StrategistSignerConfig) (*StrategistTransactor, error) {
	strategistMap := make(map[string]*bind.TransactOpts)
	for _, config := range configs {
		prvKey := os.Getenv(config.Name)
		pk, err := crypto.HexToECDSA(strings.TrimPrefix(prvKey, "0x"))
		if err != nil {
			log.Panicf("Failed to parse private key: %v", err)
		}

		transactor, err := bind.NewKeyedTransactorWithChainID(pk, new(big.Int).SetInt64(config.ChainId))
		if err != nil {
			return nil, err
		}

		strategistMap[config.Address] = transactor
	}

	return &StrategistTransactor{
		strategistTransactors: strategistMap,
	}, nil
}

func (s *StrategistTransactor) GetStrategist(address string) (*bind.TransactOpts, error) {
	strategist, ok := s.strategistTransactors[address]
	if !ok {
		return nil, ErrStrategistNotFound
	}

	return strategist, nil
}
