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
	if strings.Contains(sm[1], "CreateWallet"){
		err = process.CreateWallet()
	} else if strings.Contains(sm[1], "OpenWallet") {
		err = process.OpenWallet()
	} else if strings.Contains(sm[1], "OpenOrCreateWallet") {
		err = process.OpenOrCreateWallet()
	} else if strings.Contains(sm[1], "GetCryptScheme") {
		err = process.GetCryptScheme()
	} else if strings.Contains(sm[1], "SetCryptScheme") {
		err = process.SetCryptScheme()
	} else if strings.Contains(sm[1], "GetDefaultAccount") {
		err = process.GetDefaultAccount()
	} else if strings.Contains(sm[1], "ChangePassword") {
		err = process.ChangePassword()
	} else if strings.Contains(sm[1], "CreateAccount") {
		err = process.CreateAccount()
	} else if strings.Contains(sm[1], "GetBalance") {
		err = process.GetBalance()
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

func (s *Process) CreateWallet() error {
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		return err
	}
	password, err := s.Cons.prompter.PromptPassword("Password: ")
	if err != nil {
		return err
	}
	confirm, err := s.Cons.prompter.PromptPassword("Repeat Password: ")
	if err != nil {
		return err
	}
	if password != confirm {
		err = errors.New("Password don't match! ")
		return err
	}
	_, err = s.Cons.OntSDK.CreateWallet(walletName, password)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (s *Process) OpenWallet() error {
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		return err
	}
	password, err := s.Cons.prompter.PromptPassword("Password: ")
	if err != nil {
		return err
	}
	_, err = s.Cons.OntSDK.OpenWallet(walletName, password)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (s *Process) OpenOrCreateWallet() error {
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		return err
	}
	password, err := s.Cons.prompter.PromptPassword("Password: ")
	if err != nil {
		return err
	}
	_, err = s.Cons.OntSDK.OpenOrCreateWallet(walletName, password)
	if err != nil {
		fmt.Fprintln(s.Cons.printer, err.Error())
	}
	return err
}

func (s *Process)  GetCryptScheme() error {
	cryptScheme := s.Cons.OntSDK.GetCryptScheme()
	fmt.Fprintln(s.Cons.printer, "\n\tCryptScheme: " + cryptScheme + "\n")
	return nil
}

func (s *Process) SetCryptScheme() error {
	s.Cons.OntSDK.SetCryptScheme(common.CRYPTO_SCHEME_DEFAULT)
	cryptScheme := s.Cons.OntSDK.GetCryptScheme()
	if 0 == strings.Compare(cryptScheme, common.CRYPTO_SCHEME_DEFAULT) {
		fmt.Fprintln(s.Cons.printer, "\tSetCryptScheme success!")
	}
	return nil
}

func (s *Process) GetDefaultAccount() error {
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		return err
	}
	password, err := s.Cons.prompter.PromptPassword("Password: ")
	if err != nil {
		return err
	}
	ontWallet, err := s.Cons.OntSDK.OpenWallet(walletName, password)
	if err != nil {
		fmt.Println(err.Error())
	}
	defAccount, err := ontWallet.GetDefaultAccount()
	if err != nil {
		return err
	}

	fmt.Fprintln(s.Cons.printer, "Address: ", defAccount.Address)
	fmt.Fprintln(s.Cons.printer, "PublicKey: ", defAccount.PublicKey)
	fmt.Fprintln(s.Cons.printer, "SigScheme: ",  defAccount.SigScheme)
	fmt.Fprintln(s.Cons.printer, "PrivateKey: ",  defAccount.PrivateKey)

	return err
}

func (s *Process) ChangePassword() error {
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		return err
	}
	password, err := s.Cons.prompter.PromptPassword("Password: ")
	if err != nil {
		return err
	}
	ontWallet, err := s.Cons.OntSDK.OpenWallet(walletName, password)
	if err != nil {
		fmt.Println(err.Error())
	}
	newPassword, err := s.Cons.prompter.PromptPassword("New Password: ")
	if err != nil {
		return err
	}
	confirm, err := s.Cons.prompter.PromptPassword("Repeat New Password: ")
	if err != nil {
		return err
	}
	if password != confirm {
		err = errors.New("Password don't match! ")
		return err
	}
	err = ontWallet.ChangePassword(password, newPassword)
	return err
}

func (s *Process) CreateAccount() error {
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		return err
	}
	password, err := s.Cons.prompter.PromptPassword("Password: ")
	if err != nil {
		return err
	}
	ontWallet, err := s.Cons.OntSDK.OpenWallet(walletName, password)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	account, err := ontWallet.CreateAccount()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Fprintln(s.Cons.printer, "Address: ", account.Address)
	fmt.Fprintln(s.Cons.printer, "PublicKey: ", account.PublicKey)
	fmt.Fprintln(s.Cons.printer, "SigScheme: ",  account.SigScheme)
	fmt.Fprintln(s.Cons.printer, "PrivateKey: ",  account.PrivateKey)

	return err
}

func (s *Process) GetBalance() error {
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if walletName == "" {
		walletName = "wallet.dat"
	}
	password, err := s.Cons.prompter.PromptPassword("Password: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	ontWallet, err := s.Cons.OntSDK.OpenWallet(walletName, password)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defAccount, err := ontWallet.GetDefaultAccount()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	balance, err := s.Cons.OntSDK.Rpc.GetBalance(defAccount.Address)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("ONT: %d; ONG: %d; ONGAppove: %d\n", balance.Ont.Int64(), balance.Ong.Int64(), balance.OngAppove.Int64())
	return nil
}
