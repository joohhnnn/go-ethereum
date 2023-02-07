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

package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type KnownAccounts struct {
	Accounts map[common.Address]KnownAccount `json:"knownAccounts"`
}

func (ka *KnownAccounts) Copy() KnownAccounts {
	cpy := KnownAccounts{
		Accounts: make(map[common.Address]KnownAccount),
	}
	for key, val := range ka.Accounts {
		cpy.Accounts[key] = val.Copy()
	}
	return cpy
}

type KnownAccount struct {
	StorageRoot  *common.Hash
	StorageSlots map[common.Hash]common.Hash
}

func (ka *KnownAccount) UnmarshalJSON(data []byte) error {
	var hash common.Hash
	if err := json.Unmarshal(data, &hash); err == nil {
		ka.StorageRoot = &hash
		ka.StorageSlots = make(map[common.Hash]common.Hash)
		return nil
	}

	var mapping map[common.Hash]common.Hash
	if err := json.Unmarshal(data, &mapping); err == nil {
		ka.StorageSlots = mapping
		return nil
	}

	return errors.New("cannot unmarshal json")
}

func (ka *KnownAccount) Copy() KnownAccount {
	cpy := KnownAccount{
		StorageRoot:  nil,
		StorageSlots: map[common.Hash]common.Hash{},
	}

	if ka.StorageRoot != nil {
		*cpy.StorageRoot = *ka.StorageRoot
	}

	for key, val := range ka.StorageSlots {
		cpy.StorageSlots[key] = val
	}

	return cpy
}

func (ka *KnownAccount) Root() (common.Hash, bool) {
	if ka.StorageRoot == nil {
		return common.Hash{}, false
	}
	return *ka.StorageRoot, true
}

func (ka *KnownAccount) Slots() (map[common.Hash]common.Hash, bool) {
	if len(ka.StorageSlots) == 0 && ka.StorageRoot != nil {
		return ka.StorageSlots, false
	}
	return ka.StorageSlots, true
}

type BundleTx struct {
	Inner         *Transaction
	KnownAccounts KnownAccounts
}

func (tx *BundleTx) copy() TxData {
	cpy := &BundleTx{
		Inner:         NewTx(tx.Inner.inner),
		KnownAccounts: tx.KnownAccounts.Copy(),
	}
	return cpy
}

// accessors for innerTx.
func (tx *BundleTx) txType() byte                 { return BundleTxType }
func (tx *BundleTx) chainID() *big.Int            { return tx.Inner.ChainId() }
func (tx *BundleTx) accessList() AccessList       { return tx.Inner.AccessList() }
func (tx *BundleTx) data() []byte                 { return tx.Inner.Data() }
func (tx *BundleTx) gas() uint64                  { return tx.Inner.Gas() }
func (tx *BundleTx) gasFeeCap() *big.Int          { return tx.Inner.GasFeeCap() }
func (tx *BundleTx) gasTipCap() *big.Int          { return tx.Inner.GasTipCap() }
func (tx *BundleTx) gasPrice() *big.Int           { return tx.Inner.GasFeeCap() }
func (tx *BundleTx) value() *big.Int              { return tx.Inner.Value() }
func (tx *BundleTx) nonce() uint64                { return tx.Inner.Nonce() }
func (tx *BundleTx) to() *common.Address          { return tx.Inner.To() }
func (tx *BundleTx) knownAccounts() KnownAccounts { return tx.KnownAccounts }

func (tx *BundleTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.Inner.RawSignatureValues()
}

func (tx *BundleTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.Inner.inner.setSignatureValues(chainID, v, r, s)
}
