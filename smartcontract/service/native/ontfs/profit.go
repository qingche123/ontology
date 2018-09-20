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
	"fmt"
	"bytes"
	"encoding/json"

	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
)

func FsGetNodeList(native *native.NativeService) ([]byte, error) {
	fmt.Println("===FsGetNodeList===")
	contract := native.ContextRef.CurrentContext().ContractAddress

	nodeSetKey := GenFsNodeSetKey(contract)
	nodeSet, err := utils.GetStorageItem(native, nodeSetKey)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode,"[FS Profit] GetStorageItem nodeSetKey error!")
	}
	if nodeSet == nil {
		return utils.BYTE_FALSE, errors.NewErr("[FS Profit] FsGetNodeList No nodeSet found!")
	}

	addrSet := utils.NewSet()
	if err = addrSet.AddrDeserialize(nodeSet.Value); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode,"[FS Profit] Set deserialize error!")
	}

	items := addrSet.GetAllAddrs()
	if items == nil {
		return utils.BYTE_FALSE, nil
	}
	nodesInfo := new(bytes.Buffer)
	var fsNodeInfos FsNodeInfos
	for _ , addr := range items {
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
	fmt.Println("===FsStoreFile===")
	contract := native.ContextRef.CurrentContext().ContractAddress

	var fileInfo FileInfo
	source := common.NewZeroCopySource(native.Input)
	if err := fileInfo.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsStoreFile deserialize error!")
	}

	item, err := utils.GetStorageItem(native, fileInfo.FileHash[:])
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode,"[FS Profit] GetStorageItem error!")
	}
	if item != nil {
		return utils.BYTE_FALSE, errors.NewErr("[FS Profit] File have stored!")
	}

	set, err := getFsSetting(native)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsStoreFile getFsSetting error!")
	}

	fileInfo.Pay = (fileInfo.FileBlockNum * fileInfo.FIleBlockSize * set.GasPerKBPerHourPreserve  +
		fileInfo.ChallengeRate * fileInfo.ChallengeTimes * set.GasForChallenge) *
		fileInfo.CopyNum * set.FsGasPrice

	state := ont.State{From: fileInfo.UserAddr, To: contract, Value:fileInfo.Pay}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("FS Profit] CheckWitness failed!")
	}
	err = appCallTransfer(native, utils.OntContractAddress, state.From, state.To, state.Value)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] appCallTransferOnt, ont transfer error!")
	}
	ont.AddNotifications(native, contract, &state)

	bf := new(bytes.Buffer)
	if err = fileInfo.Serialize(bf); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] appCallTransferOnt, ont transfer error!")
	}
	utils.PutBytes(native, fileInfo.FileHash[:], bf.Bytes())
	return utils.BYTE_TRUE, nil
}

func FsGetFileInfo(native *native.NativeService) ([]byte, error){
	fmt.Println("===FsGetFileInfo===")

	source := common.NewZeroCopySource(native.Input)
	fileHash, err := utils.DecodeBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsGetFileInfo DecodeBytes error!")
	}
	fileInfo, err := getFsFileInfo(native, fileHash)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[FS Profit] FsGetFileInfo getFsFileInfo error!")
	}

	info := new(bytes.Buffer)
	fileInfo.Serialize(info)
	return info.Bytes(), nil
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
