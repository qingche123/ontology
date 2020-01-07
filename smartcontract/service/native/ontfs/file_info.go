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
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/native"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

const (
	FileStorageTypeUseSpace = 0
	FileStorageTypeUseFile  = 1
)

type FileHash struct {
	FHash []byte
}

type FileInfo struct {
	FileHash       []byte
	FileOwner      common.Address
	FileDesc       []byte
	FileBlockCount uint64
	RealFileSize   uint64
	CopyNumber     uint64
	PayAmount      uint64
	RestAmount     uint64
	FileCost       uint64
	PdpInterval    uint64
	TimeStart      uint64
	TimeExpired    uint64
	PdpParam       []byte
	ValidFlag      bool
	StorageType    uint64
}

type FileInfoList struct {
	FilesI []FileInfo
}

type FileHashList struct {
	FilesH []FileHash
}

func (this *FileInfo) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.FileHash)
	utils.EncodeAddress(sink, this.FileOwner)
	sink.WriteVarBytes(this.FileDesc)
	utils.EncodeVarUint(sink, this.FileBlockCount)
	utils.EncodeVarUint(sink, this.RealFileSize)
	utils.EncodeVarUint(sink, this.CopyNumber)
	utils.EncodeVarUint(sink, this.PayAmount)
	utils.EncodeVarUint(sink, this.RestAmount)
	utils.EncodeVarUint(sink, this.FileCost)
	utils.EncodeVarUint(sink, this.PdpInterval)
	utils.EncodeVarUint(sink, this.TimeStart)
	utils.EncodeVarUint(sink, this.TimeExpired)
	sink.WriteVarBytes(this.PdpParam)
	sink.WriteBool(this.ValidFlag)
	utils.EncodeVarUint(sink, this.StorageType)
}

func (this *FileInfo) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.FileOwner, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.FileDesc, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.FileBlockCount, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.RealFileSize, err = utils.DecodeVarUint(source)
	if err != nil {
		return nil
	}
	this.CopyNumber, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.PayAmount, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.RestAmount, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.FileCost, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.PdpInterval, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.TimeStart, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.TimeExpired, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.PdpParam, err = DecodeVarBytes(source)
	if err != nil {
		return err
	}
	this.ValidFlag, err = DecodeBool(source)
	if err != nil {
		return err
	}
	this.StorageType, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}

	return nil
}

func (this *FileInfoList) Serialization(sink *common.ZeroCopySink) {
	fileCount := uint64(len(this.FilesI))
	utils.EncodeVarUint(sink, fileCount)

	for _, fileInfo := range this.FilesI {
		sinkTmp := common.NewZeroCopySink(nil)
		fileInfo.Serialization(sinkTmp)
		sink.WriteVarBytes(sinkTmp.Bytes())
	}
}

func (this *FileInfoList) Deserialization(source *common.ZeroCopySource) error {
	fileCount, err := utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	if 0 == fileCount {
		return nil
	}

	for i := uint64(0); i < fileCount; i++ {
		var fileInfo FileInfo
		fileInfoTmp, err := DecodeVarBytes(source)
		if err != nil {
			return err
		}
		src := common.NewZeroCopySource(fileInfoTmp)
		if err = fileInfo.Deserialization(src); err != nil {
			return err
		}
		this.FilesI = append(this.FilesI, fileInfo)
	}
	return nil
}

func (this *FileHashList) Serialization(sink *common.ZeroCopySink) {
	fileCount := uint64(len(this.FilesH))
	utils.EncodeVarUint(sink, fileCount)

	if 0 != fileCount {
		for _, fileHash := range this.FilesH {
			sink.WriteVarBytes(fileHash.FHash)
		}
	}
}

func (this *FileHashList) Deserialization(source *common.ZeroCopySource) error {
	fileCount, err := utils.DecodeVarUint(source)
	if err != nil {
		return err
	}

	for i := uint64(0); i < fileCount; i++ {
		fileHashSrc, err := DecodeVarBytes(source)
		if err != nil {
			return err
		}
		fileHashTmp := make([]byte, len(fileHashSrc))
		copy(fileHashTmp, fileHashSrc)
		fileHash := FileHash{FHash: fileHashTmp}
		this.FilesH = append(this.FilesH, fileHash)
	}
	return nil
}

func addFileInfo(native *native.NativeService, fileInfo *FileInfo) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	fileInfoKey := GenFsFileInfoKey(contract, fileInfo.FileOwner, fileInfo.FileHash)

	sink := common.NewZeroCopySink(nil)
	fileInfo.Serialization(sink)

	utils.PutBytes(native, fileInfoKey, sink.Bytes())
}

func delFileInfo(native *native.NativeService, fileOwner common.Address, fileHash []byte) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	fileInfoKey := GenFsFileInfoKey(contract, fileOwner, fileHash)
	native.CacheDB.Delete(fileInfoKey)
}

func fileInfoExist(native *native.NativeService, fileOwner common.Address, fileHash []byte) bool {
	contract := native.ContextRef.CurrentContext().ContractAddress
	fileInfoKey := GenFsFileInfoKey(contract, fileOwner, fileHash)

	item, err := utils.GetStorageItem(native, fileInfoKey)
	if err != nil || item == nil || item.Value == nil {
		return false
	}
	return true
}

func getFileInfoFromDb(native *native.NativeService, fileOwner common.Address, fileHash []byte) *FileInfo {
	contract := native.ContextRef.CurrentContext().ContractAddress
	fileInfoKey := GenFsFileInfoKey(contract, fileOwner, fileHash)

	item, err := utils.GetStorageItem(native, fileInfoKey)
	if err != nil || item == nil || item.Value == nil {
		return nil
	}

	var fileInfo FileInfo
	source := common.NewZeroCopySource(item.Value)
	if err := fileInfo.Deserialization(source); err != nil {
		return nil
	}

	if fileInfo.StorageType == FileStorageTypeUseSpace {
		space := getSpaceInfoFromDb(native, fileOwner)
		if space == nil {
			fileInfo.ValidFlag = false
		} else {
			fileInfo.TimeExpired = space.TimeExpired
			if fileInfo.TimeExpired < uint64(native.Time) {
				fileInfo.ValidFlag = false
			}
		}
	}

	return &fileInfo
}

func getFileRawRealInfo(native *native.NativeService, fileOwner common.Address, fileHash []byte) []byte {
	fileInfo := getFileInfoFromDb(native, fileOwner, fileHash)
	if fileInfo == nil {
		return nil
	}

	if uint64(native.Time) > fileInfo.TimeExpired {
		fileInfo.ValidFlag = false
	}

	sink := common.NewZeroCopySink(nil)
	fileInfo.Serialization(sink)
	return sink.Bytes()
}

func getAndUpdateFileInfo(native *native.NativeService, fileOwner common.Address, fileHash []byte) *FileInfo {
	fileInfo := getFileInfoFromDb(native, fileOwner, fileHash)
	if fileInfo == nil {
		return nil
	}

	if uint64(native.Time) > fileInfo.TimeExpired {
		fileInfo.ValidFlag = false
		addFileInfo(native, fileInfo)
	}
	return fileInfo
}

func getFileInfoByHash(native *native.NativeService, fileHash []byte) *FileInfo {
	fileOwner, err := getFileOwner(native, fileHash)
	if err != nil {
		return nil
	}
	fileInfo := getAndUpdateFileInfo(native, fileOwner, fileHash)
	if fileInfo == nil {
		return nil
	}
	return fileInfo
}

func getFileHashList(native *native.NativeService, fileOwner common.Address) *FileHashList {
	contract := native.ContextRef.CurrentContext().ContractAddress

	fileInfoPrefix := GenFsFileInfoPrefix(contract, fileOwner)
	fileInfoPrefixLen := len(fileInfoPrefix)

	var fileHashList FileHashList

	iter := native.CacheDB.NewIterator(fileInfoPrefix[:])
	for has := iter.First(); has; has = iter.Next() {
		key := iter.Key()
		fileHash := FileHash{
			FHash: make([]byte, len(key[fileInfoPrefixLen:])),
		}
		copy(fileHash.FHash, key[fileInfoPrefixLen:])
		fileHashList.FilesH = append(fileHashList.FilesH, fileHash)
	}
	iter.Release()

	return &fileHashList
}

//func getFileInfoMap(native *native.NativeService, fileOwner common.Address) map[FileHash]*FileInfo {
//	contract := native.ContextRef.CurrentContext().ContractAddress
//
//	fileInfoPrefix := GenFsFileInfoPrefix(contract, fileOwner)
//	fileInfoPrefixLen := len(fileInfoPrefix)
//
//	fsFileInfoList := make(map[FileHash]*FileInfo)
//	iter := native.CacheDB.NewIterator(fileInfoPrefix[:])
//	for has := iter.First(); has; has = iter.Next() {
//		key := iter.Key()
//		item, err := utils.GetStorageItem(native, iter.Key())
//		if err != nil || item == nil || item.Value == nil {
//			log.Error("[File Info] GetPdpRecordMap GetStorageItem ", err)
//			continue
//		}
//
//		fileHash := FileHash{
//			FHash: make([]byte, len(key[fileInfoPrefixLen:])),
//		}
//		copy(fileHash.FHash, key[fileInfoPrefixLen:])
//
//		var fsFileInfo FileInfo
//		source := common.NewZeroCopySource(item.Value)
//		if err := fsFileInfo.Deserialization(source); err != nil {
//			log.Errorf("[File Info] getFileInfoMap error: ", err.Error())
//			continue
//		}
//
//		fsFileInfoList[fileHash] = &fsFileInfo
//	}
//	iter.Release()
//
//	return fsFileInfoList
//}
