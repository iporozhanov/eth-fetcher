package rlp

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/rlp"
)

func RLPtoStrings(rlphex string) ([]string, error) {
	if rlphex == "" {
		return []string{}, nil
	}

	bytes, err := hex.DecodeString(rlphex)
	if err != nil {
		return nil, err
	}

	var decoded []string

	err = rlp.DecodeBytes(bytes, &decoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
