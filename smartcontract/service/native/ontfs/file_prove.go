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

type FileProve struct {
	FileHash    []byte
	MultiRes    []byte
	AddRes      []byte
	BlockHeight uint64
	WalletAddr  common.Address
	Profit      uint64
}

func (this *FileProve) Serialize(w io.Writer) error {
	if err := utils.WriteBytes(w, this.FileHash[:]); err != nil {
		return fmt.Errorf("[FileProve] serialize from error:%v", err)
	}
	if err := utils.WriteBytes(w, this.MultiRes); err != nil {
		return fmt.Errorf("[FileProve] serialize from error:%v", err)
	}
	if err := utils.WriteBytes(w, this.AddRes); err != nil {
		return fmt.Errorf("[FileProve] serialize from error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.BlockHeight); err != nil {
		return fmt.Errorf("[FileProve] serialize from error:%v", err)
	}
	if err := utils.WriteAddress(w, this.WalletAddr); err != nil {
		return fmt.Errorf("[FileProve] serialize from error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.Profit); err != nil {
		return fmt.Errorf("[FileProve] serialize from error:%v", err)
	}
	return nil
}

func (this *FileProve) Deserialize(r io.Reader) error {
	var err error
	if this.FileHash, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileProve] deserialize from error:%v", err)
	}
	if this.MultiRes, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileProve] deserialize from error:%v", err)
	}
	if this.AddRes, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileProve] deserialize from error:%v", err)
	}
	if this.BlockHeight, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileProve] deserialize from error:%v", err)
	}
	if this.WalletAddr, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("[FileProve] deserialize from error:%v", err)
	}
	if this.Profit, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileProve] deserialize from error:%v", err)
	}
	return nil
}

func (this *FileProve) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeBytes(sink, this.FileHash)
	utils.EncodeBytes(sink, this.MultiRes)
	utils.EncodeBytes(sink, this.AddRes)
	utils.EncodeVarUint(sink, this.BlockHeight)
	utils.EncodeAddress(sink, this.WalletAddr)
	utils.EncodeVarUint(sink, this.Profit)
}

func (this *FileProve) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	this.MultiRes, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	this.AddRes, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	this.BlockHeight, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.WalletAddr, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.Profit, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	return nil
}
