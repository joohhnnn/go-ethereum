// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package policy_test

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/policy"
)

func ptr(hash common.Hash) *common.Hash {
	return &hash
}

func u64Ptr(v uint64) *uint64 {
	return &v
}

// It is known that marshaling is broken
// https://github.com/golang/go/issues/55890
func TestTxPolicyJSONUnMarshalTrip(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		mustFail bool
		expected policy.TxOptions
	}{
		{
			"StateRoot",
			`{"knownAccounts":{"0x6b3A8798E5Fb9fC5603F3aB5eA2e8136694e55d0":"0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"}}`,
			false,
			policy.TxOptions{
				KnownAccounts: map[common.Address]policy.KnownAccount{
					common.HexToAddress("0x6b3A8798E5Fb9fC5603F3aB5eA2e8136694e55d0"): policy.KnownAccount{
						StorageRoot:  ptr(common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563")),
						StorageSlots: make(map[common.Hash]common.Hash),
					},
				},
			},
		},
		{
			"StorageSlots",
			`{"knownAccounts":{"0x6b3A8798E5Fb9fC5603F3aB5eA2e8136694e55d0":{"0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8":"0x0000000000000000000000000000000000000000000000000000000000000000"}}}`,
			false,
			policy.TxOptions{
				KnownAccounts: map[common.Address]policy.KnownAccount{
					common.HexToAddress("0x6b3A8798E5Fb9fC5603F3aB5eA2e8136694e55d0"): policy.KnownAccount{
						StorageRoot: nil,
						StorageSlots: map[common.Hash]common.Hash{
							common.HexToHash("0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8"): common.HexToHash("0x"),
						},
					},
				},
			},
		},
		{
			"EmptyObject",
			`{"knownAccounts":{}}`,
			false,
			policy.TxOptions{
				KnownAccounts: make(map[common.Address]policy.KnownAccount),
			},
		},
		{
			"EmptyStrings",
			`{"knownAccounts":{"":""}}`,
			true,
			policy.TxOptions{
				KnownAccounts: nil,
			},
		},
		{
			"BlockNumberMin",
			`{"blockNumberMin":"0x1"}`,
			false,
			policy.TxOptions{
				BlockNumberMin: big.NewInt(1),
			},
		},
		{
			"BlockNumberMax",
			`{"blockNumberMin":"0x1", "blockNumberMax":"0x2"}`,
			false,
			policy.TxOptions{
				BlockNumberMin: big.NewInt(1),
				BlockNumberMax: big.NewInt(2),
			},
		},
		{
			"TimestampMin",
			`{"timestampMin":"0xffff"}`,
			false,
			policy.TxOptions{
				TimestampMin: u64Ptr(uint64(0xffff)),
			},
		},
		{
			"TimestampMax",
			`{"timestampMax":"0xffffff"}`,
			false,
			policy.TxOptions{
				TimestampMax: u64Ptr(uint64(0xffffff)),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var opts policy.TxOptions
			err := json.Unmarshal([]byte(test.input), &opts)
			if test.mustFail && err == nil {
				t.Errorf("Test %s should fail", test.name)
			}
			if !test.mustFail && err != nil {
				t.Errorf("Test %s should pass but got err: %v", test.name, err)
			}

			if !reflect.DeepEqual(opts, test.expected) {
				t.Errorf("Test %s got unexpected value, want %#v, got %#v", test.name, test.expected, opts)
			}
		})
	}
}
