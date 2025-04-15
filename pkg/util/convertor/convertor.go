package convertor

import (
	"github.com/ethereum/go-ethereum/common"
)

func MapDecodersAndSanitizersToAddresses(proofs []string) []common.Address {
	addresses := make([]common.Address, len(proofs))
	for i, proof := range proofs {
		addresses[i] = common.HexToAddress(proof)
	}
	return addresses
}

func MapManageProofs(proofs [][]string) [][][32]byte {
	manageProofs := make([][][32]byte, len(proofs))
	for i, group := range proofs {
		manageProofs[i] = make([][32]byte, len(group))
		for j, proof := range group {
			manageProofs[i][j] = common.HexToHash(proof)
		}
	}
	return manageProofs
}
