package nucleus

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

type AddressBook map[int64]NetworkData

type NetworkData struct {
	ID        string            `json:"name"`
	Multisig  string            `json:"multisig,omitempty"`
	Token     map[string]string `json:"token,omitempty"`
	Hyperswap HyperswapData     `json:"hyperswap,omitempty"`
	Nucleus   ChainConfig       `json:"nucleus"`
}

type HyperswapData struct {
	Router                     string `json:"router"`
	NonfungiblePositionManager string `json:"nonfungiblePositionManager"`
}

type VaultDetail struct {
	BoringVault    string `json:"boring_vault"`
	Manager        string `json:"manager"`
	Accountant     string `json:"accountant"`
	Teller         string `json:"teller"`
	RolesAuthority string `json:"roles_authority"`
}

type ChainConfig struct {
	RoosterMicroManager string                 `json:"roosterMicroManager"`
	Vaults              map[string]VaultDetail `json:"-"`
}

func (n *ChainConfig) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	n.Vaults = make(map[string]VaultDetail)
	for key, value := range raw {
		if key == "roosterMicroManager" {
			if err := json.Unmarshal(value, &n.RoosterMicroManager); err != nil {
				return err
			}
		} else {
			var vault VaultDetail
			if err := json.Unmarshal(value, &vault); err != nil {
				return err
			}
			n.Vaults[key] = vault
		}
	}
	return nil
}

type Transaction struct {
	Target common.Address `json:"target"`
	Val    *big.Int       `json:"value"`
	Data   string         `json:"calldata"`
}

type MerkleProofs struct {
	ManageProofs          [][]string `json:"proofs"`
	DecodersAndSanitizers []string   `json:"decoderAndSanitizerAddress"`
}

type Calldata struct {
	ManageProofs          [][][32]byte     `json:"proofs"`
	DecodersAndSanitizers []common.Address `json:"decoderAndSanitizerAddress"`
	Targets               []common.Address `json:"targets"`
	TargetData            [][]byte         `json:"targetData"`
	Values                []*big.Int       `json:"values"`
}

type StrategistSignerConfig struct {
	Address string `json:"address" mapstructure:"address"`
	Name    string `json:"name" mapstructure:"name"`
	ChainId int64  `json:"chainId" mapstructure:"chainId"`
}

func parseAddressBook(data []byte) (AddressBook, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	book := make(AddressBook)

	for key, val := range raw {
		id, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			continue
		}

		var net NetworkData
		if err := json.Unmarshal(val, &net); err != nil {
			return nil, err
		}
		book[id] = net
	}

	return book, nil
}
