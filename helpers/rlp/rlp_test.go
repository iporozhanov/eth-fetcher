package rlp_test

import (
	"errors"
	"testing"

	"eth-fetcher/helpers/rlp"
)

func TestRLPtoStrings(t *testing.T) {
	tests := []struct {
		name     string
		rlphex   string
		expected []string
		err      error
	}{
		{
			name:     "Empty input",
			rlphex:   "",
			expected: []string{},
		},
		{
			name:     "Single string",
			rlphex:   "c68568656c6c6f",
			expected: []string{"hello"},
		},
		{
			name:     "Multiple strings",
			rlphex:   "cc8568656c6c6f85776f726c64",
			expected: []string{"hello", "world"},
		},
		{
			name:     "Multiple hexes",
			rlphex:   "f90110b842307839623266366133633265316165643263636366393262613636366332326430353361643064386135646137616131666435343737646364363537376234353234b842307835613537653330353163623932653264343832353135623037653762336431383531373232613734363534363537626436346131346333396361336639636632b842307837316239653262343464343034393863303861363239383866616337373664306561633062356239363133633337663966366639613462383838613862303537b842307863356639366266316235346433333134343235643233373962643737643765643465363434663763366538343961373438333230323862333238643464373938", // ["hello", "world"]
			expected: []string{"0x9b2f6a3c2e1aed2cccf92ba666c22d053ad0d8a5da7aa1fd5477dcd6577b4524", "0x5a57e3051cb92e2d482515b07e7b3d1851722a74654657bd64a14c39ca3f9cf2", "0x71b9e2b44d40498c08a62988fac776d0eac0b5b9613c37f9f6f9a4b888a8b057", "0xc5f96bf1b54d3314425d2379bd77d7ed4e644f7c6e849a74832028b328d4d798"},
		},
		{
			name:     "Invalid hex input",
			rlphex:   "invalid",
			expected: nil,
			err:      errors.New("invalid hex input"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := rlp.RLPtoStrings(test.rlphex)
			if err != nil && test.err == nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && test.err != nil {
				t.Errorf("Expected error: %v, but got nil", test.err)
			}
			if len(result) != len(test.expected) {
				t.Errorf("Expected %d strings, but got %d", len(test.expected), len(result))
			}
			for i, str := range result {
				if str != test.expected[i] {
					t.Errorf("Expected string '%s', but got '%s'", test.expected[i], str)
				}
			}
		})
	}
}
