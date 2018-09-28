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

type FileReadSettleSlice struct {
	FileHash    []byte
	PayFrom     common.Address
	PayTo       common.Address
	SlicePay    uint64
	SliceId     uint64
	Sig         []byte
	PubKey      []byte
}

func (this *FileReadSettleSlice) Serialize(w io.Writer) error {
	if err := utils.WriteBytes(w, this.FileHash); err != nil {
		return fmt.Errorf("[FileReadPledge] [FileHash:%v] serialize from error:%v", this.FileHash, err)
	}
	if err := utils.WriteAddress(w, this.PayFrom); err != nil {
		return fmt.Errorf("[FileReadPledge] [PayFrom:%v] serialize from error:%v", this.PayFrom, err)
	}
	if err := utils.WriteAddress(w, this.PayTo); err != nil {
		return fmt.Errorf("[FileReadPledge] [PayTo:%v] serialize from error:%v", this.PayTo, err)
	}
	if err := utils.WriteVarUint(w, this.SlicePay); err != nil {
		return fmt.Errorf("[FileReadPledge] [SlicePay:%v] serialize from error:%v", this.SlicePay, err)
	}
	if err := utils.WriteVarUint(w, this.SliceId); err != nil {
		return fmt.Errorf("[FileReadPledge] [SliceId:%v] serialize from error:%v", this.SliceId, err)
	}
	if err := utils.WriteBytes(w, this.Sig); err != nil {
		return fmt.Errorf("[FileReadPledge] [Sig:%v] serialize from error:%v", this.Sig, err)
	}
	if err := utils.WriteBytes(w, this.PubKey); err != nil {
		return fmt.Errorf("[FileReadPledge] [PubKey:%v] serialize from error:%v", this.PubKey, err)
	}
	return nil
}

func (this *FileReadSettleSlice) Deserialize(r io.Reader) error {
	var err error
	if this.FileHash, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [FileHash] deserialize from error:%v", err)
	}
	if this.PayFrom, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [PayFrom] deserialize from error:%v", err)
	}
	if this.PayTo, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [PayTo] deserialize from error:%v", err)
	}
	if this.SlicePay, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [SlicePay] deserialize from error:%v", err)
	}
	if this.SliceId, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [SliceId] deserialize from error:%v", err)
	}
	if this.Sig, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [Sig] deserialize from error:%v", err)
	}
	if this.PubKey, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[FileReadPledge] [PubKey] deserialize from error:%v", err)
	}
	return nil
}


func (this *FileReadSettleSlice) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeBytes(sink, this.FileHash[:])
	utils.EncodeAddress(sink, this.PayFrom)
	utils.EncodeAddress(sink, this.PayTo)
	utils.EncodeVarUint(sink, this.SlicePay)
	utils.EncodeVarUint(sink, this.SliceId)
	utils.EncodeBytes(sink, this.Sig)
	utils.EncodeBytes(sink, this.PubKey)
}

func (this *FileReadSettleSlice) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FileHash, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	this.PayFrom, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.PayTo, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}
	this.SlicePay, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.SliceId, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.Sig, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	this.PubKey, err = utils.DecodeBytes(source)
	if err != nil {
		return err
	}
	return nil
}
