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

type FileReadPledge struct {
	FileHash    []byte
	ReadAddr    common.Address
	FromAddr    common.Address
	Id          uint64
	TotalValue  uint64
	RestValue   uint64
}

func (this *FileReadPledge) Serialize(w io.Writer) error {
	if err := utils.WriteBytes(w, this.FileHash); err != nil {
		return fmt.Errorf("[FileReadPledge] [FileHash:%v] serialize from error:%v", this.FileHash, err)
	}
	if err := utils.WriteAddress(w, this.ReadAddr); err != nil {
		return fmt.Errorf("[FileReadPledge] [ReadAddr:%v] serialize from error:%v", this.ReadAddr, err)
	}
	if err := utils.WriteAddress(w, this.FromAddr); err != nil {
		return fmt.Errorf("[FileReadPledge] [FromAddr:%v] serialize from error:%v", this.FromAddr, err)
	}
	if err := utils.WriteVarUint(w, this.Id); err != nil {
		return fmt.Errorf("[FileReadPledge] [Id:%v] serialize from error:%v", this.Id, err)
	}
	if err := utils.WriteVarUint(w, this.TotalValue); err != nil {
		return fmt.Errorf("[FileReadPledge] [TotalValue:%v] serialize from error:%v", this.TotalValue, err)
	}
	if err := utils.WriteVarUint(w, this.RestValue); err != nil {
		return fmt.Errorf("[FileReadPledge] [RestValue:%v] serialize from error:%v", this.RestValue, err)
	}
	return nil
}

func (this *FileReadPledge) Deserialize(r io.Reader) error {
	var err error
	if this.FileHash, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [FileHash] deserialize from error:%v", err)
	}
	if this.ReadAddr, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [ReadAddr] deserialize from error:%v", err)
	}
	if this.FromAddr, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [FromAddr] deserialize from error:%v", err)
	}
	if this.Id, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [Id] deserialize from error:%v", err)
	}
	if this.TotalValue, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [TotalValue] deserialize from error:%v", err)
	}
	if this.RestValue, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [RestValue] deserialize from error:%v", err)
	}
	return nil
}

func (this *FileReadPledge) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeBytes(sink, this.FileHash[:])
	utils.EncodeAddress(sink, this.ReadAddr)
	utils.EncodeAddress(sink, this.FromAddr)
	utils.EncodeVarUint(sink, this.Id)
	utils.EncodeVarUint(sink, this.TotalValue)
	utils.EncodeVarUint(sink, this.RestValue)
}

func (this *FileReadPledge) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	this.ReadAddr, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.FromAddr, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.Id, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.TotalValue, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.RestValue, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	return nil
}
