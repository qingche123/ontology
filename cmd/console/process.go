/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package console

import (
	"fmt"
	"strings"
	"reflect"
	"errors"

	"github.com/ontio/ontology-go-sdk/rpc"
	"github.com/ontio/ontology-go-sdk/wallet"
	"github.com/ontio/ontology-go-sdk/common"
)

var process * Process
var keyWords []string

func Call(statement string, cons *Console) error {
	var err error
	process.Cons = cons

	sm := strings.Split(statement, ".")
	if true {
		// todo
	} else {
		fmt.Println("Unknown statement: ", sm[0], " ", sm[1])
		err = errors.New(fmt.Sprint("Unknown statement: ", sm[0], " ", sm[1]))
	}
	return err
}

type Process struct {
	Cons *Console
}

func NewProcess(cons *Console) *Process{
	collectKeyWords(cons)
	return &Process{Cons:cons}
}

func collectKeyWords(cons *Console){
	sdkType := reflect.TypeOf(cons.OntSDK)
	sdkMethodNum := sdkType.NumMethod()
	for i := 0; i < sdkMethodNum ; i++ {
		keyWords = append(keyWords, "wallet." + sdkType.Method(i).Name)
	}

	rpcClt :=  rpc.NewRpcClient(common.CRYPTO_SCHEME_DEFAULT)
	rpcType := reflect.TypeOf(rpcClt)
	rpcMethodNum := rpcType.NumMethod()
	for i := 0; i < rpcMethodNum ; i++ {
		keyWords = append(keyWords, "rpc." + rpcType.Method(i).Name)
	}

	ontWallet := &wallet.OntWallet{}
	ontWalletType := reflect.TypeOf(ontWallet)
	ontWalletMethodNum := ontWalletType.NumMethod()
	for i := 0; i < ontWalletMethodNum ; i++ {
		keyWords = append(keyWords, "wallet." + ontWalletType.Method(i).Name)
	}

	keyWords = append(keyWords, "wallet")
	keyWords = append(keyWords, "rpc")
	keyWords = append(keyWords, "contract")
}

func CompleteKeywords(line string) []string {
	var rightKeyWords []string
	for i := 0; i < len(keyWords); i++  {
		b := strings.HasPrefix(keyWords[i], line)
		if b {
			rightKeyWords = append(rightKeyWords, keyWords[i])
		}
	}
	return rightKeyWords
}