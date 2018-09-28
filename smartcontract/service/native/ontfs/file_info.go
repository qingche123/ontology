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
	"io"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

type FileInfo struct {
	FileHash       []byte
	UserAddr       common.Address
	FileBlockNum   uint64
	FileBlockSize  uint64
	ChallengeRate  uint64
	ChallengeTimes uint64
	CopyNum        uint64
	Deposit        uint64
	FileProveParam []byte
	ProveBlockNum  uint64
	BlockHeight    uint64
}

func (this *FileInfo) Serialize(w io.Writer) error {
	if err := utils.WriteBytes(w, this.FileHash); err != nil {
		return fmt.Errorf("[FileInfo] [FileHash:%v] serialize from error:%v", this.FileHash, err)
	}
	if err := utils.WriteAddress(w, this.UserAddr); err != nil {
		return fmt.Errorf("[FileInfo] [UserAddr:%v] serialize from error:%v", this.UserAddr, err)
	}
	if err := utils.WriteVarUint(w, this.FileBlockNum); err != nil {
		return fmt.Errorf("[FileInfo] [FileBlockNum:%v] serialize from error:%v", this.FileBlockNum, err)
	}
	if err := utils.WriteVarUint(w, this.FileBlockSize); err != nil {
		return fmt.Errorf("[FileInfo] [FileBlockSize:%v] serialize from error:%v", this.FileBlockSize, err)
	}
	if err := utils.WriteVarUint(w, this.ChallengeRate); err != nil {
		return fmt.Errorf("[FileInfo] [ChallengeRate:%v] serialize from error:%v", this.ChallengeRate, err)
	}
	if err := utils.WriteVarUint(w, this.ChallengeTimes); err != nil {
		return fmt.Errorf("[FileInfo] [ChallengeTimes:%v] serialize from error:%v", this.ChallengeTimes, err)
	}
	if err := utils.WriteVarUint(w, this.CopyNum); err != nil {
		return fmt.Errorf("[FileInfo] [CopyNum:%v] serialize from error:%v", this.CopyNum, err)
	}
	if err := utils.WriteVarUint(w, this.Deposit); err != nil {
		return fmt.Errorf("[FileInfo] [Deposit:%v] serialize from error:%v", this.Deposit, err)
	}
	if err := utils.WriteBytes(w, this.FileProveParam); err != nil {
		return fmt.Errorf("[FileInfo] [FileProveParam:%v] serialize from error:%v", this.FileProveParam, err)
	}
	if err := utils.WriteVarUint(w, this.ProveBlockNum); err != nil {
		return fmt.Errorf("[FileInfo] [ProveBlockNum:%v] serialize from error:%v", this.ProveBlockNum, err)
	}
	if err := utils.WriteVarUint(w, this.BlockHeight); err != nil {
		return fmt.Errorf("[FileInfo] [BlockHeight:%v] serialize from error:%v", this.BlockHeight, err)
	}
	return nil
}

func (this *FileInfo) Deserialize(r io.Reader) error {
	var err error
	if this.FileHash, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileInfo] [FileHash] deserialize from error:%v", err)
	}
	if this.UserAddr, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("[FileInfo] [UserAddr] deserialize from error:%v", err)
	}
	if this.FileBlockNum, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileInfo] [FileBlockNum] deserialize from error:%v", err)
	}
	if this.FileBlockSize, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileInfo] [FileBlockSize] deserialize from error:%v", err)
	}
	if this.ChallengeRate, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileInfo] [ChallengeRate] deserialize from error:%v", err)
	}
	if this.ChallengeTimes, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileInfo] [ChallengeTimes] deserialize from error:%v", err)
	}
	if this.CopyNum, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileInfo] [CopyNum] deserialize from error:%v", err)
	}
	if this.Deposit, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileInfo] [Deposit] deserialize from error:%v", err)
	}
	if this.FileProveParam, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileInfo] [FileProveParam] deserialize from error:%v", err)
	}
	if this.ProveBlockNum, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileInfo] [ProveBlockNum] deserialize from error:%v", err)
	}
	if this.BlockHeight, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileInfo] [BlockHeight] deserialize from error:%v", err)
	}
	return nil
}

func (this *FileInfo) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeBytes(sink, this.FileHash)
	utils.EncodeAddress(sink, this.UserAddr)
	utils.EncodeVarUint(sink, this.FileBlockNum)
	utils.EncodeVarUint(sink, this.FileBlockSize)
	utils.EncodeVarUint(sink, this.ChallengeRate)
	utils.EncodeVarUint(sink, this.ChallengeTimes)
	utils.EncodeVarUint(sink, this.CopyNum)
	utils.EncodeVarUint(sink, this.Deposit)
	utils.EncodeBytes(sink, this.FileProveParam)
	utils.EncodeVarUint(sink, this.ProveBlockNum)
	utils.EncodeVarUint(sink, this.BlockHeight)
}

func (this *FileInfo) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	this.UserAddr, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.FileBlockNum, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.FileBlockSize, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.ChallengeRate, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.ChallengeTimes, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.CopyNum, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.Deposit, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.FileProveParam, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	this.ProveBlockNum, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.BlockHeight, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	return nil
}
