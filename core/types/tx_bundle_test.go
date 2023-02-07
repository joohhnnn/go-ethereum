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

package types_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ptr(hash common.Hash) *common.Hash {
	return &hash
}

func TestKnownAccountJSONUnmarshal(t *testing.T) {
	tests := []struct {
		input    string
		mustFail bool
		expected types.KnownAccounts
	}{
		0: {
			`{"knownAccounts":{"0x6b3A8798E5Fb9fC5603F3aB5eA2e8136694e55d0":"0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"}}`,
			false,
			types.KnownAccounts{
				Accounts: map[common.Address]types.KnownAccount{
					common.HexToAddress("0x6b3A8798E5Fb9fC5603F3aB5eA2e8136694e55d0"): types.KnownAccount{
						StorageRoot:  ptr(common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563")),
						StorageSlots: make(map[common.Hash]common.Hash),
					},
				},
			},
		},
		1: {
			`{"knownAccounts":{"0x6b3A8798E5Fb9fC5603F3aB5eA2e8136694e55d0":{"0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8":"0x0000000000000000000000000000000000000000000000000000000000000000"}}}`,
			false,
			types.KnownAccounts{
				Accounts: map[common.Address]types.KnownAccount{
					common.HexToAddress("0x6b3A8798E5Fb9fC5603F3aB5eA2e8136694e55d0"): types.KnownAccount{
						StorageRoot: nil,
						StorageSlots: map[common.Hash]common.Hash{
							common.HexToHash("0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8"): common.HexToHash("0x"),
						},
					},
				},
			},
		},
		2: {
			`{"knownAccounts":{}}`,
			false,
			types.KnownAccounts{
				Accounts: map[common.Address]types.KnownAccount{},
			},
		},
		3: {
			`{"knownAccounts":{"":""}}`,
			true,
			types.KnownAccounts{
				Accounts: map[common.Address]types.KnownAccount{},
			},
		},
	}

	for i, test := range tests {
		var ka types.KnownAccounts
		err := json.Unmarshal([]byte(test.input), &ka)
		if test.mustFail && err == nil {
			t.Errorf("Test %d should fail", i)
			continue
		}
		if !test.mustFail && err != nil {
			t.Errorf("Test %d should pass but got err: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(ka, test.expected) {
			t.Errorf("Test %d got unexpected value, want %d, got %d", i, test.expected, ka)
		}
	}
}
