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

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

func InitFs() {
	native.Contracts[utils.OntFSContractAddress] = RegisterFsContract
}

func RegisterFsContract(native *native.NativeService) {
	//native.Register(FS_INIT, FsInit)
	native.Register(FS_GETSETTING, FsGetSetting)
	native.Register(FS_NODE_REGISTER, FsNodeRegister)
	native.Register(FS_NODE_QUERY, FsNodeQuery)
	native.Register(FS_NODE_UPDATE, FsNodeUpdate)
	native.Register(FS_NODE_CANCEL, FsNodeCancel)
	native.Register(FS_GET_NODE_LIST, FsGetNodeList)
	native.Register(FS_STORE_FILE, FsStoreFile)
	native.Register(FS_GET_FILE_INFO, FsGetFileInfo)
	native.Register(FS_NODE_WITH_DRAW_PROFIT, FsNodeWithDrawProfit)
	native.Register(FS_FILE_PROVE, FsFileProve)
	native.Register(FS_GET_FILE_PROVE_DETAILS, FsGetFileProveDetails)
	native.Register(FS_READ_FILE_PLEDGE, FsReadFilePledge)
	native.Register(FS_FILE_READ_PROFIT_SETTLE, FsFileReadProfitSettle)
	native.Register(FS_DELETE_FILE, FsDeleteFile)
}

func fsInit(native *native.NativeService) ([]byte, error) {
	var fsSetting FsSetting
	infoSource := common.NewZeroCopySource(native.Input)
	if err := fsSetting.Deserialization(infoSource); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] FsSetting deserialize error!")
	}

	setFsSetting(native, fsSetting)
	return utils.BYTE_TRUE, nil
}

func FsGetSetting(native *native.NativeService) ([]byte, error) {
	fsSetting, err := getFsSetting(native)
	if err != nil || fsSetting == nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] GetFsSetting error!")
	}
	fs := new(bytes.Buffer)
	fsSetting.Serialize(fs)
	return fs.Bytes(), nil
}

func getFsSetting(native *native.NativeService) (*FsSetting, error) {
	var fsSetting FsSetting
	contract := native.ContextRef.CurrentContext().ContractAddress
	fsSettingKey := GenFsSettingKey(contract)

	item, err := utils.GetStorageItem(native, fsSettingKey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Govern] GetFsSetting error!")
	}
	if item == nil {
		fsSetting = FsSetting{
			FsGasPrice:       2000,
			GasPerKBPerBlock: 1,
			GasPerKBForRead:  1,
			GasForChallenge:  1,
		}
		return &fsSetting, nil
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
