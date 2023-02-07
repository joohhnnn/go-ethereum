package state

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func IsValidBundleTransaction(tx *types.Transaction, state *StateDB) (bool, error) {
	if tx.Type() != types.BundleTxType {
		return false, nil
	}
	knownAccounts := tx.KnownAccounts()
	for addr, acct := range knownAccounts.Accounts {
		root, isRoot := acct.Root()
		if isRoot {
			storageTrie, err := state.StorageTrie(addr)
			if err != nil {
				return false, err
			}
			if storageTrie.Hash() != root {
				return false, nil
			}
		}
		slots, isSlots := acct.Slots()
		if isSlots {
			for key, value := range slots {
				if value != state.GetState(addr, key) {
					return false, nil
				}
			}
		}
	}
	return true, nil
}
