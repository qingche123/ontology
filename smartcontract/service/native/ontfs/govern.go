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

package ontfs

import (
	"bytes"
	"fmt"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
)

func InitFs() {
	native.Contracts[utils.OntFSContractAddress] = RegisterFsContract
}

func RegisterFsContract(native *native.NativeService) {
	native.Register(FS_SET, FsSet)
	native.Register(FS_SET_INIT, FsSetInit)
	native.Register(FS_GETSETTING, FsGetSetting)
	native.Register(FS_NODE_REGISTER, FsNodeRegister)
	native.Register(FS_NODE_QUERY, FsNodeQuery)
	native.Register(FS_NODE_UPDATE, FsNodeUpdate)
	native.Register(FS_NODE_CANCEL, FsNodeCancel)
}

func FsSetInit(native *native.NativeService) ([]byte, error) {
	var fsSetting FsSetting

	fsSetting.FsGasPrice = 1
	fsSetting.GasPerKBForStore = 1
	fsSetting.GasPerKBForRead = 1

	setFsSetting(native, fsSetting)
	return utils.BYTE_TRUE, nil
}

func FsSet(native *native.NativeService) ([]byte, error) {
	fmt.Println("===FsSet===")
	var fsSetting FsSetting
	infoSource := common.NewZeroCopySource(native.Input)
	if err := fsSetting.Deserialization(infoSource); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsSetting deserialize error!")
	}

	setFsSetting(native, fsSetting)
	return utils.BYTE_TRUE, nil
}

func FsGetSetting(native *native.NativeService) ([]byte, error) {
	fmt.Println("===FsGetSetting===")
	fsSetting, err := getFsSetting(native)
	if err != nil || fsSetting == nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] GetFsSetting error!")
	}
	fs := new(bytes.Buffer)
	fsSetting.Serialize(fs)
	return fs.Bytes(), nil
}

func FsNodeRegister(native *native.NativeService) ([]byte, error) {
	fmt.Println("===FsNodeRegister===")
	contract := native.ContextRef.CurrentContext().ContractAddress

	fsSetting, err := getFsSetting(native)
	if err != nil || fsSetting == nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] GetFsSetting error!")
	}

	var fsNodeInfo FsNodeInfo
	infoSource := common.NewZeroCopySource(native.Input)
	if err := fsNodeInfo.Deserialization(infoSource); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OntFS Govern] FsNodeInfo deserialize error!")
	}

	fsNodeInfoKey := GenFsNodeInfoKey(contract, fsNodeInfo.WalletAddr)
	item, err := utils.GetStorageItem(native, fsNodeInfoKey)
	if err == nil && item != nil{
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] Node have registered!")
	}

	pledge := fsSetting.FsGasPrice * fsSetting.GasPerKBForStore * fsNodeInfo.Volume
	//===========================================================================
	state := ont.State{From: fsNodeInfo.WalletAddr, To: contract, Value:pledge}

	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("FS Govern] CheckWitness failed!")
	}
	err = appCallTransfer(native, utils.OntContractAddress, state.From, state.To, state.Value)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] appCallTransferOnt, ont transfer error!")
	}
	ont.AddNotifications(native, contract, &state)
	//===========================================================================

	fsNodeInfo.Pledge = pledge
	info := new(bytes.Buffer)
	fsNodeInfo.Serialize(info)
	utils.PutBytes(native, fsNodeInfoKey, info.Bytes())

	return utils.BYTE_TRUE, nil
}

func FsNodeQuery(native *native.NativeService) ([]byte, error) {
	fmt.Println("===FsNodeQuery===")
	contract := native.ContextRef.CurrentContext().ContractAddress

	source := common.NewZeroCopySource(native.Input)
	walletAddr, err := utils.DecodeAddress(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] DecodeAddress error!")
	}

	fsNodeInfoKey := GenFsNodeInfoKey(contract, walletAddr)
	item, err := utils.GetStorageItem(native, fsNodeInfoKey)
	if err != nil || item == nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsNodeInfo GetStorageItem error!")
	}

	var fsNodeInfo FsNodeInfo
	fsNodeInfoSource := common.NewZeroCopySource(item.Value)
	err = fsNodeInfo.Deserialization(fsNodeInfoSource)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsNodeInfo deserialize error!")
	}

	info := new(bytes.Buffer)
	fsNodeInfo.Serialize(info)
	return info.Bytes(), nil
}

func FsNodeUpdate(native *native.NativeService) ([]byte, error) {
	fmt.Println("===FsNodeUpdate===")
	contract := native.ContextRef.CurrentContext().ContractAddress

	fsSetting, err := getFsSetting(native)
	if err != nil || fsSetting == nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] GetFsSetting error!")
	}

	var newFsNodeInfo, oldFsNodeInfo FsNodeInfo
	newInfoSource := common.NewZeroCopySource(native.Input)
	if err := newFsNodeInfo.Deserialization(newInfoSource); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsNodeInfo deserialize error!")
	}

	fsNodeInfoKey := GenFsNodeInfoKey(contract, newFsNodeInfo.WalletAddr)
	item, err := utils.GetStorageItem(native, fsNodeInfoKey)
	if err != nil || item == nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsNodeInfo GetStorageItem error!")
	}

	fsNodeInfoSource := common.NewZeroCopySource(item.Value)
	err = oldFsNodeInfo.Deserialization(fsNodeInfoSource)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsNodeInfo deserialize error!")
	}
	newPledge := fsSetting.FsGasPrice * fsSetting.GasPerKBForStore * newFsNodeInfo.Volume

	if newFsNodeInfo.WalletAddr != oldFsNodeInfo.WalletAddr {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsNodeInfo walletAddr changed!")
	}

	var state ont.State
	if newPledge < oldFsNodeInfo.Pledge {
		state = ont.State{From:contract, To:newFsNodeInfo.WalletAddr, Value:oldFsNodeInfo.Pledge - newPledge}
	} else if newPledge > oldFsNodeInfo.Pledge {
		state = ont.State{From:newFsNodeInfo.WalletAddr, To:contract, Value:newPledge - oldFsNodeInfo.Pledge}
	}
	if newPledge != oldFsNodeInfo.Pledge {
		if native.ContextRef.CheckWitness(newFsNodeInfo.WalletAddr) == false {
			return utils.BYTE_FALSE, errors.NewErr("FS Govern] CheckWitness failed!")
		}

		err = appCallTransfer(native, utils.OntContractAddress, state.From, state.To, state.Value)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] appCallTransferOnt, ont transfer error!")
		}
		ont.AddNotifications(native, contract, &state)
	}

	newFsNodeInfo.Pledge = newPledge
	info := new(bytes.Buffer)
	newFsNodeInfo.Serialize(info)
	utils.PutBytes(native, fsNodeInfoKey, info.Bytes())
	return utils.BYTE_TRUE, nil
}

func FsNodeCancel(native *native.NativeService) ([]byte, error) {
	fmt.Println("===FsNodeCancel===")
	contract := native.ContextRef.CurrentContext().ContractAddress

	source := common.NewZeroCopySource(native.Input)
	addr, err := utils.DecodeAddress(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] DecodeAddress error!")
	}

	fsNodeInfoKey := GenFsNodeInfoKey(contract, addr)

	var fsNodeInfo FsNodeInfo
	item, err := utils.GetStorageItem(native, fsNodeInfoKey)
	if err != nil || item == nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsNodeInfo GetStorageItem error!")
	}
	fsNodeInfoSource := common.NewZeroCopySource(item.Value)
	err = fsNodeInfo.Deserialization(fsNodeInfoSource)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsNodeInfo deserialize error!")
	}
	//===========================================================================
	var state ont.State
	if fsNodeInfo.Pledge > 0 {
		state = ont.State{From:contract, To:fsNodeInfo.WalletAddr, Value:fsNodeInfo.Pledge}
		if native.ContextRef.CheckWitness(state.To) == false {
			return utils.BYTE_FALSE, errors.NewErr("FS Govern] CheckWitness failed!")
		}
		err = appCallTransfer(native, utils.OntContractAddress, contract, state.To, state.Value)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] appCallTransferOnt, ont transfer error!")
		}
		ont.AddNotifications(native, contract, &state)
	}
	//===========================================================================
	utils.DelStorageItem(native, fsNodeInfoKey)
	return utils.BYTE_TRUE, nil
}

func getFsSetting(native *native.NativeService) (*FsSetting, error){
	contract := native.ContextRef.CurrentContext().ContractAddress

	var fsSetting FsSetting
	fsSettingKey := GenFsSettingKey(contract)

	item, err := utils.GetStorageItem(native, fsSettingKey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] GetFsSetting error!")
	}
	if item == nil {
		return nil, fmt.Errorf("[FS Govern] Not found fsSetting")
	}
	settingSource := common.NewZeroCopySource(item.Value)
	if err := fsSetting.Deserialization(settingSource); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsSetting Deserialization error!")
	}
	return &fsSetting, nil
}

func setFsSetting(native *native.NativeService, fsSetting FsSetting) {
	contract := native.ContextRef.CurrentContext().ContractAddress

	info := new(bytes.Buffer)
	fsSetting.Serialize(info)

	fsSettingKey := GenFsSettingKey(contract)
	utils.PutBytes(native, fsSettingKey, info.Bytes())
}

func appCallTransfer(native *native.NativeService, contract common.Address, from common.Address, to common.Address, amount uint64) error {
	var sts []ont.State
	sts = append(sts, ont.State{
		From:  from,
		To:    to,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	sink := common.NewZeroCopySink(nil)
	transfers.Serialization(sink)

	if _, err := native.NativeCall(contract, "transfer", sink.Bytes()); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "appCallTransfer, appCall error!")
	}
	return nil
}
