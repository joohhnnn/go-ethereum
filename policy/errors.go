// Copyright 2022 The go-ethereum Authors
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

package policy

import (
	"github.com/ethereum/go-ethereum/rpc"
)

// TxOptionsError is a standardized error message for eip-4337 UserOperations
// also containing any custom error message Geth might include.
type TxOptionsError struct {
	code int
	msg  string
	err  error
}

func (e *TxOptionsError) ErrorCode() int { return e.code }
func (e *TxOptionsError) Error() string  { return e.msg }
func (e *TxOptionsError) ErrorData() interface{} {
	if e.err == nil {
		return nil
	}
	return struct {
		Error string `json:"err"`
	}{e.err.Error()}
}

// With returns a copy of the error with a new embedded custom data field.
func (e *TxOptionsError) With(err error) *TxOptionsError {
	return &TxOptionsError{
		code: e.code,
		msg:  e.msg,
		err:  err,
	}
}

var (
	_ rpc.Error     = new(TxOptionsError)
	_ rpc.DataError = new(TxOptionsError)
)

var (
	//TODO: confirm whether the code should be the same.
	OutOfTimestampRange   = &TxOptionsError{code: -32503, msg: "Out of timestamp range"}
	OutOfBlockNumberRange = &TxOptionsError{code: -32503, msg: "Out of blockNumber range"}
	KnownAccountsNotMatch = &TxOptionsError{code: -32503, msg: "knownAccounts is not match"}
)
