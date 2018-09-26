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
	"encoding/json"
	"fmt"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

func FsGetNodeList(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress

	nodeSetKey := GenFsNodeSetKey(contract)
	nodeSet, err := utils.GetStorageItem(native, nodeSetKey)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] GetStorageItem nodeSetKey error!")
	}
	if nodeSet == nil {
		return utils.BYTE_FALSE, errors.NewErr("[FS Profit] FsGetNodeList No nodeSet found!")
	}

	addrSet := utils.NewSet()
	if err = addrSet.AddrDeserialize(nodeSet.Value); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] Set deserialize error!")
	}

	items := addrSet.GetAllAddrs()
	if items == nil {
		return utils.BYTE_FALSE, nil
	}
	nodesInfo := new(bytes.Buffer)
	var fsNodeInfos FsNodeInfos
	for _, addr := range items {
		fsNodeInfo, err := getFsNodeInfo(native, addr)
		if err != nil {
			fmt.Errorf("[FS Profit] FsGetNodeList getFsNodeInfo(%v) error", addr)
			continue
		}
		fsNodeInfos.Item = append(fsNodeInfos.Item, *fsNodeInfo)
	}
	data, err := json.Marshal(fsNodeInfos)
	if _, err = nodesInfo.Write(data); err != nil {
		return utils.BYTE_FALSE, nil
	}
	return nodesInfo.Bytes(), nil
}

func FsStoreFile(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress

	var fileInfo FileInfo
	source := common.NewZeroCopySource(native.Input)
	if err := fileInfo.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsStoreFile deserialize error!")
	}

	item, err := utils.GetStorageItem(native, fileInfo.FileHash[:])
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] GetStorageItem error!")
	}
	if item != nil {
		return utils.BYTE_FALSE, errors.NewErr("[FS Profit] File have stored!")
	}

	set, err := getFsSetting(native)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsStoreFile getFsSetting error!")
	}

	fileInfo.Deposit = (fileInfo.FileBlockNum*fileInfo.FileBlockSize*set.GasPerKBPerBlock +
		fileInfo.ChallengeRate*fileInfo.ChallengeTimes*set.GasForChallenge) *
		fileInfo.CopyNum * set.FsGasPrice

	state := ont.State{From: fileInfo.UserAddr, To: contract, Value: fileInfo.Deposit}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("FS Profit] CheckWitness failed!")
	}
	err = appCallTransfer(native, utils.OngContractAddress, state.From, state.To, state.Value)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] appCallTransferOnt, ont transfer error!")
	}
	ont.AddNotifications(native, contract, &state)

	fileInfo.ProveBlockNum = 32
	fileInfo.BlockHeight = uint64(native.Height)

	bf := new(bytes.Buffer)
	if err = fileInfo.Serialize(bf); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsStoreFile fileInfo serialize error!")
	}
	utils.PutBytes(native, fileInfo.FileHash[:], bf.Bytes())

	var proveDetails FsProveDetails
	proveDetails.CopyNum = fileInfo.CopyNum
	proveDetails.ProveDetailNum = 0
	proveBuff := new(bytes.Buffer)
	if err = proveDetails.Serialize(proveBuff); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] ProveDetails serialize error!")
	}
	proveDetailsKey := GenFsProveDetailsKey(contract, fileInfo.FileHash)
	utils.PutBytes(native, proveDetailsKey, proveBuff.Bytes())

	return utils.BYTE_TRUE, nil
}

func FsGetFileInfo(native *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(native.Input)
	fileHash, err := utils.DecodeBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsGetFileInfo DecodeBytes error!")
	}
	item, err := utils.GetStorageItem(native, fileHash)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsFileInfo GetStorageItem error!")
	}
	if item == nil {
		return nil, errors.NewErr("[FS Profit] FsFileInfo not found!")
	}
	return item.Value, nil
}

func FsGetFileProveDetails(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	source := common.NewZeroCopySource(native.Input)
	fileHash, err := utils.DecodeBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsGetFileProveDetails DecodeBytes error!")
	}

	fileProveDetailKey := GenFsProveDetailsKey(contract, fileHash)
	item, err := utils.GetStorageItem(native, fileProveDetailKey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsGetFileProveDetails GetStorageItem error!")
	}
	if item == nil {
		return nil, errors.NewErr("[FS Profit] FsGetFileProveDetails not found!")
	}
	return item.Value, nil
}

func FsReadFilePledge(native *native.NativeService) ([]byte, error){
	contract := native.ContextRef.CurrentContext().ContractAddress
	var fileReadPledge FileReadPledge
	source := common.NewZeroCopySource(native.Input)
	err := fileReadPledge.Deserialization(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsReadFilePledge deserialization error!")
	}

	pledge, err := getReadFilePledge(native, fileReadPledge.FileHash)
	if err == nil && pledge != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsReadFilePledge file read pledged error!")
	}

	state := ont.State{From: fileReadPledge.FromAddr, To: contract, Value: fileReadPledge.TotalValue}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("FS Profit] CheckWitness failed!")
	}
	err = appCallTransfer(native, utils.OngContractAddress, state.From, state.To, state.Value)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] appCallTransferOnt, ont transfer error!")
	}
	ont.AddNotifications(native, contract, &state)

	fileReadPledge.Id = 0
	fileInfo, err := getFsFileInfo(native, fileReadPledge.FileHash)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsReadFilePledge getFsFileInfo error!")
	}

	fsSetting, err := getFsSetting(native)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsReadFilePledge getFsSetting error!")
	}

	readMinFee := fileInfo.FileBlockNum * fileInfo.FileBlockSize * fsSetting.FsGasPrice * fsSetting.GasPerKBForRead
	if fileReadPledge.TotalValue < readMinFee {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsReadFilePledge insufficient pay!")
	}

	key := GenFsFileReadPledgeKey(contract, fileReadPledge.FileHash)
	bf := new(bytes.Buffer)
	err = fileReadPledge.Serialize(bf)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsReadFilePledge serialize error!")
	}
	utils.PutBytes(native, key, bf.Bytes())
	return utils.BYTE_TRUE, nil
}

func FsDeleteFile(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress

	source := common.NewZeroCopySource(native.Input)
	fileHash, err := utils.DecodeBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsGetFileInfo DecodeBytes error!")
	}
	utils.DelStorageItem(native, fileHash[:])
	proveDetailsKey := GenFsProveDetailsKey(contract, fileHash)
	utils.DelStorageItem(native, proveDetailsKey)
	return utils.BYTE_TRUE, nil
}

func getReadFilePledge(native *native.NativeService, fileHash []byte) (*FileReadPledge, error){
	contract := native.ContextRef.CurrentContext().ContractAddress
	key := GenFsFileReadPledgeKey(contract, fileHash)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] GetReadFilePledge GetStorageItem error!")
	}
	if item == nil {
		return nil, errors.NewErr("[FS Profit] FileReadPledge not found!")
	}

	var fileReadPledge FileReadPledge
	fileReadPledgeSource := common.NewZeroCopySource(item.Value)
	err = fileReadPledge.Deserialization(fileReadPledgeSource)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] GetReadFilePledge deserialize error!")
	}
	return &fileReadPledge, nil
}

func getFsFileInfo(native *native.NativeService, fileHash []byte) (*FileInfo, error) {
	item, err := utils.GetStorageItem(native, fileHash)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsFileInfo GetStorageItem error!")
	}
	if item == nil {
		return nil, errors.NewErr("[FS Profit] FsFileInfo not found!")
	}

	var fsFileInfo FileInfo
	fsFileInfoSource := common.NewZeroCopySource(item.Value)
	err = fsFileInfo.Deserialization(fsFileInfoSource)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsFileInfo deserialize error!")
	}
	return &fsFileInfo, nil
}

func getFsFileProveDetail(native *native.NativeService, fileHash []byte) (*ProveDetail, error) {
	item, err := utils.GetStorageItem(native, fileHash)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] ProveDetail GetStorageItem error!")
	}
	if item == nil {
		return nil, errors.NewErr("[FS Profit] FsFileInfo not found!")
	}
	reader := bytes.NewReader(item.Value)
	var fileProveDetail ProveDetail
	err = fileProveDetail.Deserialize(reader)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] ProveDetail deserialize error!")
	}
	return &fileProveDetail, nil
}
