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
	"os"
	"fmt"
	"bytes"
	"time"
	"strings"
	"reflect"
	"errors"
	"strconv"
	"math/big"
	"encoding/hex"
	"encoding/json"

	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology-go-sdk/rpc"
	"github.com/ontio/ontology-go-sdk/wallet"
	"github.com/ontio/ontology-go-sdk/common"
	sctypes "github.com/ontio/ontology/smartcontract/types"
	comm "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
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
	} else if strings.Contains(sm[1], "DeploySmartContract") {
		err = process.DeploySmartContract()
	} else if strings.Contains(sm[1], "GetBalanceWithBase58") {
		err = process.GetBalanceWithBase58()
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
	if walletName == "" {
		walletName = "wallet.dat"
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
	if walletName == "" {
		walletName = "wallet.dat"
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
	if walletName == "" {
		walletName = "wallet.dat"
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
	pro := s.Cons.prompter
	var names []string = []string{
		"SHA224withECDSA",
		"SHA256withECDSA",
		"SHA384withECDSA",
		"SHA512withECDSA",
		"SHA3-224withECDSA",
		"SHA3-256withECDSA",
		"SHA3-384withECDSA",
		"SHA3-512withECDSA",
		"RIPEMD160withECDSA",
		"SM3withSM2",
		"SHA512withEdDSA",
	}

CRYPTSCHEME:
	fmt.Fprintln(s.Cons.printer, "Crypt Scheme: ")
	for n := range names{
		fmt.Fprintln(s.Cons.printer, n, ". ", names[n])
	}

	choiceStr, err := pro.PromptInput("Choose the crypt scheme: ")
	if err != nil {
		fmt.Println("addr:", err.Error())
		return nil
	}
	choice, _ := strconv.Atoi(choiceStr)
	if choice < 0 && choice > 11{
		goto CRYPTSCHEME
	}
	s.Cons.OntSDK.SetCryptScheme(names[choice])
	return nil
}

func (s *Process) GetDefaultAccount() error {
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		return err
	}
	if walletName == "" {
		walletName = "wallet.dat"
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
	if walletName == "" {
		walletName = "wallet.dat"
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
	if walletName == "" {
		walletName = "wallet.dat"
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

func (s *Process) DeploySmartContract() error {
	pro := s.Cons.prompter

	walletName, err := pro.PromptInput("Wallet Name:")
	if err != nil {
		fmt.Println("WalletName:", err.Error())
		return err
	}
	if walletName == "" {
		walletName = "wallet.dat"
	}
	password, err := pro.PromptPassword("Password:")
	if err != nil {
		fmt.Println("Password:", err.Error())
		return err
	}
	ontWallet, err := s.Cons.OntSDK.OpenWallet(walletName, password)
	if err != nil {
		fmt.Println("OpenWallet:", err.Error())
		return err
	}
	defAccount, err := ontWallet.GetDefaultAccount()
	if err != nil {
		fmt.Println("GetDefaultAccount:", err.Error())
		return err
	}

	var vmType sctypes.VmType
	var tmp int

VMTYPE:
	fmt.Println("-----------------------------------------------")
	fmt.Println("| VM Type:    1.Native    2.NEOVM    3.WASMVM |")
	fmt.Println("-----------------------------------------------")
	vmTypeStr, err := pro.PromptInput("Choose the vm type:")
	if err != nil {
		fmt.Println("VM Type:", err.Error())
		return err
	}

	tmp, _ = strconv.Atoi(vmTypeStr)
	switch tmp {
	case 1:
		vmType = sctypes.Native
	case 2:
		vmType = sctypes.NEOVM
	case 3:
		vmType = sctypes.WASMVM
	default:
		goto VMTYPE
	}

	needStorage, err := pro.PromptConfirm("Does it need to be store: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	contractName, err := pro.PromptInput("Input the contract name: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	contractCode, err := pro.PromptInput("Input the contract code: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	contractVersion, err := pro.PromptInput("Input the code version: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	authorName, err := pro.PromptInput("Input the author name: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	email, err := pro.PromptInput("Input your email: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	desc, err := pro.PromptInput("Input the contract description: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	hash, err := s.Cons.OntSDK.Rpc.DeploySmartContract(defAccount, vmType, needStorage, contractCode,
		contractName, contractVersion, authorName, email, desc)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Deploy success, Contract Hash: ", hash)
	return nil
}

func (s *Process)GetBalanceWithBase58() error {
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

	balance, err := s.Cons.OntSDK.Rpc.GetBalanceWithBase58(defAccount.Address.ToBase58())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("ONT: %d; ONG: %d; ONGAppove: %d\n", balance.Ont.Int64(), balance.Ong.Int64(), balance.OngAppove.Int64())
	return nil
}

func (s *Process)GetBlockByHash() error{
	blockHash, err := s.Cons.prompter.PromptInput("Block Hash[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	hashBytes, err := hex.DecodeString(blockHash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	hash, err := comm.Uint256ParseFromBytes(hashBytes)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	block, err := s.Cons.OntSDK.Rpc.GetBlockByHash(hash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoInvGracefully(block)
	return nil
}

func (s *Process)GetBlockByHashWithHexString () error{
	blockHash, err := s.Cons.prompter.PromptInput("Block Hash[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	block, err := s.Cons.OntSDK.Rpc.GetBlockByHashWithHexString(blockHash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoInvGracefully(block)
	return nil
}

func (s *Process)GetBlockByHeight () error{
	blockHeightStr, err := s.Cons.prompter.PromptInput("Block Height: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	blockHeight, err := strconv.Atoi(blockHeightStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	block, err := s.Cons.OntSDK.Rpc.GetBlockByHeight(uint32(blockHeight))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoInvGracefully(block)
	return nil
}

func (s *Process)GetBlockCount () error{
	blockCount, err := s.Cons.OntSDK.Rpc.GetBlockCount()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("BlockCount: ", blockCount)
	return nil
}

func (s *Process)GetBlockHash () error{
	blockHeightStr, err := s.Cons.prompter.PromptInput("Block Height: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	blockHeight, err := strconv.Atoi(blockHeightStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	hash, err := s.Cons.OntSDK.Rpc.GetBlockHash(uint32(blockHeight))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("BlockHash: ", hex.EncodeToString(hash.ToArray()))
	return nil
}

func (s *Process)GetCurrentBlockHash () error{
	hash, err := s.Cons.OntSDK.Rpc.GetCurrentBlockHash()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("CurrentBlockHash: ", hex.EncodeToString(hash.ToArray()))
	return nil
}

func (s *Process)GetGenerateBlockTime () error{
	time, err := s.Cons.OntSDK.Rpc.GetGenerateBlockTime()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("GenerateBlockTime: ", time)
	return nil
}

func (s *Process)GetMerkleProof () error{
	txHash, err := s.Cons.prompter.PromptInput("Transaction Hash[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	hashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	hash, err := comm.Uint256ParseFromBytes(hashBytes)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	merkleProof, err := s.Cons.OntSDK.Rpc.GetMerkleProof(hash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoMerkleGracefully(merkleProof)
	return nil
}

func (s *Process)GetMerkleProofWithHexString () error{
	txHash, err := s.Cons.prompter.PromptInput("Transaction Hash[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	merkleProof, err := s.Cons.OntSDK.Rpc.GetMerkleProofWithHexString(txHash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoMerkleGracefully(merkleProof)
	return nil
}

func (s *Process)GetRawTransaction () error{
	txHash, err := s.Cons.prompter.PromptInput("Transaction Hash[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	hashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	hash, err := comm.Uint256ParseFromBytes(hashBytes)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	rawTx, err := s.Cons.OntSDK.Rpc.GetRawTransaction(hash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoInvGracefully(rawTx)
	return nil
}

func (s *Process) GetRawTransactionWithHexString () error{
	txHash, err := s.Cons.prompter.PromptInput("Transaction Hash[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	rawTx, err := s.Cons.OntSDK.Rpc.GetRawTransactionWithHexString(txHash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoInvGracefully(rawTx)
	return nil
}

func (s *Process) GetSmartContract () error{
	scAddrStr, err := s.Cons.prompter.PromptInput("SmartContract Address[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	scAddrBytes, _ := hex.DecodeString(scAddrStr)
	scAddr, _ := comm.AddressParseFromBytes(scAddrBytes)
	deployCode, err := s.Cons.OntSDK.Rpc.GetSmartContract(scAddr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoSmartContractGracefully(deployCode)
	return nil
}

func (s *Process) GetSmartContractEvent () error{
	txHash, err := s.Cons.prompter.PromptInput("Transaction Hash[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	hashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	hash, err := comm.Uint256ParseFromBytes(hashBytes)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	smartContractEvent, err := s.Cons.OntSDK.Rpc.GetSmartContractEvent(hash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoSmartContractEventGracefully(smartContractEvent)
	return nil
}

func (s *Process) GetSmartContractEventWithHexString () error{
	txHash, err := s.Cons.prompter.PromptInput("Transaction Hash[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	smartContractEvent, err := s.Cons.OntSDK.Rpc.GetSmartContractEventWithHexString(txHash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	echoSmartContractEventGracefully(smartContractEvent)
	return nil
}

func (s *Process) GetStorage() error{
	scAddrStr, err := s.Cons.prompter.PromptInput("SmartContract Address[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	scAddrBytes, _ := hex.DecodeString(scAddrStr)
	scAddr, _ := comm.AddressParseFromBytes(scAddrBytes)

	scKeyStr, err := s.Cons.prompter.PromptInput("SmartContract Key[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	scKeyBytes, _ := hex.DecodeString(scKeyStr)

	storage, err := s.Cons.OntSDK.Rpc.GetStorage(scAddr, scKeyBytes)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(hex.EncodeToString(storage))
	return nil
}

func (s *Process) GetVersion() error{
	version, err := s.Cons.OntSDK.Rpc.GetVersion()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Ontology Version", version)
	return nil
}

func (s *Process) InvokeNeoVMSmartContract() error{
	pro := s.Cons.prompter
	walletName, err := pro.PromptInput("Wallet Name:")
	if err != nil {
		fmt.Println("WalletName:", err.Error())
		return err
	}
	if walletName == "" {
		walletName = "wallet.dat"
	}
	password, err := pro.PromptPassword("Password:")
	if err != nil {
		fmt.Println("Password:", err.Error())
		return err
	}
	ontWallet, err := s.Cons.OntSDK.OpenWallet(walletName, password)
	if err != nil {
		fmt.Println("OpenWallet:", err.Error())
		return err
	}
	defAccount, err := ontWallet.GetDefaultAccount()
	if err != nil {
		fmt.Println("GetDefaultAccount:", err.Error())
		return err
	}
	gasLimitStr, err := pro.PromptInput("GasLimit:")
	if err != nil {
		fmt.Println("GasLimit:", err.Error())
		return err
	}
	gasLimit, _ := strconv.Atoi(gasLimitStr)
	bGsaLimit := big.NewInt(int64(gasLimit))

	scAddrStr, err := s.Cons.prompter.PromptInput("SmartContract Address[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	scAddrBytes, _ := hex.DecodeString(scAddrStr)
	scAddr, _ := comm.AddressParseFromBytes(scAddrBytes)

	hash, err := s.Cons.OntSDK.Rpc.InvokeNeoVMSmartContract(defAccount, bGsaLimit, scAddr, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Invoke success, Contract Hash: ", hash)
	return nil
}

func (s *Process) InvokeWasmVMSmartContract() error{
	pro := s.Cons.prompter
	walletName, err := pro.PromptInput("Wallet Name:")
	if err != nil {
		fmt.Println("WalletName:", err.Error())
		return err
	}
	if walletName == "" {
		walletName = "wallet.dat"
	}
	password, err := pro.PromptPassword("Password:")
	if err != nil {
		fmt.Println("Password:", err.Error())
		return err
	}
	ontWallet, err := s.Cons.OntSDK.OpenWallet(walletName, password)
	if err != nil {
		fmt.Println("OpenWallet:", err.Error())
		return err
	}
	defAccount, err := ontWallet.GetDefaultAccount()
	if err != nil {
		fmt.Println("GetDefaultAccount:", err.Error())
		return err
	}
	gasLimitStr, err := pro.PromptInput("GasLimit:")
	if err != nil {
		fmt.Println("GasLimit:", err.Error())
		return err
	}
	gasLimit, _ := strconv.Atoi(gasLimitStr)
	bGsaLimit := big.NewInt(int64(gasLimit))

	scAddrStr, err := s.Cons.prompter.PromptInput("SmartContract Address[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	scAddrBytes, _ := hex.DecodeString(scAddrStr)
	scAddr, _ := comm.AddressParseFromBytes(scAddrBytes)
	methodName, err := s.Cons.prompter.PromptInput("Input the method name: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

PARAMTYPE:
	var paramType wasmvm.ParamType
	fmt.Println("-----------------------------------------------")
	fmt.Println("| Param Type:    1.Json          2.Raw        |")
	fmt.Println("-----------------------------------------------")
	paramTypeStr, err := s.Cons.prompter.PromptInput("Input the paramType: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	tmp, _ := strconv.Atoi(paramTypeStr)
	switch tmp {
	case 1:
		paramType = wasmvm.Json
	case 2:
		paramType = wasmvm.Raw
	default:
		goto PARAMTYPE
	}

	versionStr, err := s.Cons.prompter.PromptInput("Input the version: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	version, _ := strconv.Atoi(versionStr)
	hash, err := s.Cons.OntSDK.Rpc.InvokeWasmVMSmartContract(defAccount, bGsaLimit, scAddr, methodName, paramType, byte(version), nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Invoke success, Contract Hash: ", hash)
	return nil
}

func (s *Process) PrepareInvokeNeoVMSmartContract() error {
	pro := s.Cons.prompter
	gasLimitStr, err := pro.PromptInput("GasLimit:")
	if err != nil {
		fmt.Println("GasLimit:", err.Error())
		return err
	}
	gasLimit, _ := strconv.Atoi(gasLimitStr)
	bGsaLimit := big.NewInt(int64(gasLimit))

	scAddrStr, err := s.Cons.prompter.PromptInput("SmartContract Address[HexString]: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	scAddrBytes, _ := hex.DecodeString(scAddrStr)
	scAddr, _ := comm.AddressParseFromBytes(scAddrBytes)

	NVTYPE:
	var nvType common.NeoVMReturnType
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("| NeoVMReturn Type:                                            |")
	fmt.Println("| 1.NEOVM_TYPE_BOOL            2.NEOVM_TYPE_INTEGER            |")
	fmt.Println("| 3.NEOVM_TYPE_BYTE_ARRAY      4.NEOVM_TYPE_STRING             |")
	fmt.Println("| 4.NEOVM_TYPE_ARRAY                                           |")
	fmt.Println("----------------------------------------------------------------")

	nvTypeStr, err := s.Cons.prompter.PromptInput("Input the paramType: ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	tmp, _ := strconv.Atoi(nvTypeStr)
	switch tmp {
	case 1:
		nvType = common.NEOVM_TYPE_BOOL
	case 2:
		nvType = common.NEOVM_TYPE_INTEGER
	case 3:
		nvType = common.NEOVM_TYPE_BYTE_ARRAY
	case 4:
		nvType = common.NEOVM_TYPE_STRING
	case 5:
		nvType = common.NEOVM_TYPE_ARRAY
	default:
		goto NVTYPE
	}

	hash, err := s.Cons.OntSDK.Rpc.PrepareInvokeNeoVMSmartContract(bGsaLimit, scAddr, nil, nvType)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Invoke success, Contract Hash: ", hash)
	return nil
}

func (s *Process)SendRawTransaction() error {
	var tx *types.Transaction
	//todo:
	hash, err := s.Cons.OntSDK.Rpc.SendRawTransaction(tx)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Invoke success, Contract Hash: ", hash)
	return nil
}

func (s *Process)SetAddress() *rpc.RpcClient {
	pro := s.Cons.prompter
	addr, err := pro.PromptInput("Input the address: ")
	if err != nil {
		fmt.Println("addr:", err.Error())
		return nil
	}

	rpcClt := s.Cons.OntSDK.Rpc.SetAddress(addr)
	if rpcClt == nil {
		fmt.Println("SetAddress the return value is nil")
		return nil
	}
	fmt.Println("SetAddress success")
	return rpcClt
}

func (s *Process)RPCSetCryptScheme() error {
	pro := s.Cons.prompter
	var names []string = []string{
		"SHA224withECDSA",
		"SHA256withECDSA",
		"SHA384withECDSA",
		"SHA512withECDSA",
		"SHA3-224withECDSA",
		"SHA3-256withECDSA",
		"SHA3-384withECDSA",
		"SHA3-512withECDSA",
		"RIPEMD160withECDSA",
		"SM3withSM2",
		"SHA512withEdDSA",
	}

	CRYPTSCHEME:
	fmt.Fprintln(s.Cons.printer, "Crypt Scheme: ")
	for n := range names{
		fmt.Fprintln(s.Cons.printer, n, ". \t", names[n])
	}

	choiceStr, err := pro.PromptInput("Choose the crypt scheme: ")
	if err != nil {
		fmt.Println("addr:", err.Error())
		return nil
	}
	choice, _ := strconv.Atoi(choiceStr)
	if choice < 0 && choice > 11{
		goto CRYPTSCHEME
	}
	s.Cons.OntSDK.Rpc.SetCryptScheme(names[choice])
	return nil
}

func (s *Process)SetHttpClient() *rpc.RpcClient {
	//todo:
/*	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(25 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
*/
	rpcClt := s.Cons.OntSDK.Rpc.SetHttpClient(nil)
	if rpcClt == nil {
		fmt.Println("SetHttpClient the return value is nil")
		return nil
	}
	fmt.Println("SetHttpClient success")
	return rpcClt
}

func (s *Process)Transfer() error {
	pro := s.Cons.prompter
	walletName, err := s.Cons.prompter.PromptInput("Wallet Name: ")
	if err != nil {
		return err
	}
	if walletName == "" {
		walletName = "wallet.dat"
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
	token, err := pro.PromptInput("Input the token name: ")
	if err != nil {
		fmt.Println("addr:", err.Error())
		return nil
	}
	toAddr, err := pro.PromptInput("Input the address which you transfer to : ")
	if err != nil {
		fmt.Println("addr:", err.Error())
		return nil
	}
	base58ToAddr, err := comm.AddressFromBase58(toAddr)
	toAccount := &account.Account{Address:base58ToAddr}

	countStr, err := pro.PromptInput("Input the amount you will transfer : ")
	if err != nil {
		fmt.Println("addr:", err.Error())
		return nil
	}
	count, _ := strconv.Atoi(countStr)
	bnCount := big.NewInt(int64(count))

	hash, err := s.Cons.OntSDK.Rpc.Transfer(token, defAccount, toAccount, bnCount)
	if err == nil {
		fmt.Println("Transfer the return value is nil")
		return err
	}
	fmt.Println("Send Transaction success, Tx Hash: ", hash)
	return nil
}

func (s *Process)WaitForGenerateBlock() error {
	pro := s.Cons.prompter
	heightStr, err := pro.PromptInput("Input the block height: ")
	if err != nil {
		fmt.Println("addr:", err.Error())
		return nil
	}
	height, _ := strconv.Atoi(heightStr)
	timeStr, err := pro.PromptInput("Input the block height: ")
	if err != nil {
		fmt.Println("addr:", err.Error())
		return nil
	}
	timeDur, _ := strconv.Atoi(timeStr)
	blockGenerated, err := s.Cons.OntSDK.Rpc.WaitForGenerateBlock(time.Duration(timeDur), uint32(height))
	if err != nil {
		fmt.Println("WaitForGenerateBlock the return value is nil")
		return err
	}
	if blockGenerated {
		fmt.Println("Block has been generated")
	} else {
		fmt.Println("WaitForGenerateBlock failed")
	}

	return nil
}

func echoSmartContractEventGracefully(smartContractEvent []*common.SmartContactEvent) {
	for e := range smartContractEvent {
		fmt.Println("[", e, "]TransactionHash:", smartContractEvent[e].TxHash)
		fmt.Println("[", e, "]ContractAddress:", smartContractEvent[e].ContractAddress)
	}
}

func echoMerkleGracefully(merkleProof *common.MerkleProof) {
	fmt.Println("BlockHeight:", merkleProof.BlockHeight)
	fmt.Println("CurBlockHeight:", merkleProof.CurBlockHeight)
	fmt.Println("CurBlockRoot:", merkleProof.CurBlockRoot)
	fmt.Println("TargetHashes:", merkleProof.TargetHashes)
	fmt.Println("TransactionsRoot:", merkleProof.TransactionsRoot)
	fmt.Println("Type:", merkleProof.Type)
}

func echoSmartContractGracefully(deployCode *payload.DeployCode) {
	fmt.Println("Version: ", deployCode.Version)
	fmt.Println("Name: ", deployCode.Name)
	fmt.Println("Author: ", deployCode.Author)
	fmt.Println("Code: ", deployCode.Code)
	fmt.Println("Description: ", deployCode.Description)
	fmt.Println("Email: ", deployCode.Email)
	fmt.Println("NeedStorage: ", deployCode.NeedStorage)
}

func echoInvGracefully(inv interface{}) {
	jsons, errs := json.Marshal(inv)
	if errs != nil {
		fmt.Println("Marshal json err:%s", errs.Error())
	}

	var out bytes.Buffer
	err := json.Indent(&out, jsons, "", "\t")
	if err != nil {
		fmt.Println("Gracefully format json err: %s", err.Error())
	}
	out.WriteTo(os.Stdout)
}
