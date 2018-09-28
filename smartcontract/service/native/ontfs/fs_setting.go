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

type FsSetting struct {
	FsGasPrice       uint64
	GasPerKBPerBlock uint64 //for store file
	GasPerKBForRead  uint64 //for read file
	GasForChallenge  uint64 //for challenge
	MaxProveBlockNum uint64
}

func (this *FsSetting) Serialize(w io.Writer) error {
	if err := utils.WriteVarUint(w, this.FsGasPrice); err != nil {
		return fmt.Errorf("[FsSetting] [FsGasPrice:%v] serialize from error:%v", this.FsGasPrice, err)
	}
	if err := utils.WriteVarUint(w, this.GasPerKBPerBlock); err != nil {
		return fmt.Errorf("[FsSetting] [GasPerKBPerBlock:%v] serialize from error:%v", this.GasPerKBPerBlock, err)
	}
	if err := utils.WriteVarUint(w, this.GasPerKBForRead); err != nil {
		return fmt.Errorf("[FsSetting] [GasPerKBForRead:%v] serialize from error:%v", this.GasPerKBForRead, err)
	}
	if err := utils.WriteVarUint(w, this.GasForChallenge); err != nil {
		return fmt.Errorf("[FsSetting] [GasForChallenge:%v] serialize from error:%v", this.GasForChallenge, err)
	}
	if err := utils.WriteVarUint(w, this.MaxProveBlockNum); err != nil {
		return fmt.Errorf("[FsSetting] [MaxProveBlockNum:%v] serialize from error:%v", this.MaxProveBlockNum, err)
	}
	return nil
}

func (this *FsSetting) Deserialize(r io.Reader) error {
	var err error
	if this.FsGasPrice, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsSetting] [FsGasPrice] Deserialize from error:%v", err)
	}
	if this.GasPerKBPerBlock, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsSetting] [GasPerKBPerBlock] Deserialize from error:%v", err)
	}
	if this.GasPerKBForRead, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsSetting] [GasPerKBForRead] Deserialize from error:%v", err)
	}
	if this.GasForChallenge, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsSetting] [GasForChallenge] Deserialize from error:%v", err)
	}
	if this.MaxProveBlockNum, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsSetting] [MaxProveBlockNum] Deserialize from error:%v", err)
	}
	return nil
}

func (this *FsSetting) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarUint(sink, this.FsGasPrice)
	utils.EncodeVarUint(sink, this.GasPerKBPerBlock)
	utils.EncodeVarUint(sink, this.GasPerKBForRead)
	utils.EncodeVarUint(sink, this.GasForChallenge)
	utils.EncodeVarUint(sink, this.MaxProveBlockNum)
}

func (this *FsSetting) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FsGasPrice, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.GasPerKBPerBlock, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.GasPerKBForRead, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.GasForChallenge, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.MaxProveBlockNum, err = utils.DecodeVarUint(source)
	return err
}
