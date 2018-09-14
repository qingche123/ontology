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
	"github.com/ontio/ontology/common/serialization"
)

// Transfers
type FsNodeInfo struct {
	Pledge      uint64
	Volume      uint64
	ServiceTime uint64
	WalletAddr  common.Address
	NodeAddr    []byte
}

type FsSetting struct {
	FsGasPrice       uint64
	GasPerKBForStore uint64
	GasPerKBForRead  uint64
}

func (this *FsNodeInfo) Serialize(w io.Writer) error {
	if err := utils.WriteVarUint(w, this.Pledge); err != nil {
		return fmt.Errorf("[FsNodeInfo] serialize from error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.Volume); err != nil {
		return fmt.Errorf("[FsNodeInfo] serialize from error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.ServiceTime); err != nil {
		return fmt.Errorf("[FsNodeInfo] serialize from error:%v", err)
	}
	if err := utils.WriteAddress(w, this.WalletAddr); err != nil {
		return fmt.Errorf("[FsNodeInfo] serialize from error:%v", err)
	}
	if err := serialization.WriteVarBytes(w, this.NodeAddr); err != nil {
		return fmt.Errorf("[FsNodeInfo] serialize from error:%v", err)
	}
	return nil
}

func (this *FsNodeInfo) Deserialize(r io.Reader) error {
	var err error
	if this.Pledge, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsNodeInfo] Deserialize from error:%v", err)
	}
	if this.Volume, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsNodeInfo] Deserialize from error:%v", err)
	}
	if this.ServiceTime, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsNodeInfo] Deserialize from error:%v", err)
	}
	if this.WalletAddr, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("[FsNodeInfo] Deserialize from error:%v", err)
	}
	if this.NodeAddr, err = serialization.ReadVarBytes(r); err != nil {
		return fmt.Errorf("[FsNodeInfo] Deserialize from error:%v", err)
	}
	return nil
}

func (this *FsNodeInfo) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarUint(sink, this.Pledge)
	utils.EncodeVarUint(sink, this.Volume)
	utils.EncodeVarUint(sink, this.ServiceTime)
	utils.EncodeAddress(sink, this.WalletAddr)
	sink.WriteVarBytes(this.NodeAddr)
}

func (this *FsNodeInfo) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.Pledge, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.Volume, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.ServiceTime, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.WalletAddr, err = utils.DecodeAddress(source)
	if err != nil {
		return err
	}

	from, _, irregular, eof := source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}
	this.NodeAddr = from
	return nil
}

func (this *FsSetting) Serialize(w io.Writer) error {
	if err := utils.WriteVarUint(w, this.FsGasPrice); err != nil {
		return fmt.Errorf("[FsSetting] serialize from error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.GasPerKBForStore); err != nil {
		return fmt.Errorf("[FsSetting] serialize from error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.GasPerKBForRead); err != nil {
		return fmt.Errorf("[FsSetting] serialize from error:%v", err)
	}
	return nil
}

func (this *FsSetting) Deserialize(r io.Reader) error {
	var err error
	if this.FsGasPrice, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsSetting] Deserialize from error:%v", err)
	}
	if this.GasPerKBForStore, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsSetting] Deserialize from error:%v", err)
	}
	if this.GasPerKBForRead, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[FsSetting] Deserialize from error:%v", err)
	}
	return nil
}

func (this *FsSetting) Serialization(sink *common.ZeroCopySink) {
	utils.EncodeVarUint(sink, this.FsGasPrice)
	utils.EncodeVarUint(sink, this.GasPerKBForStore)
	utils.EncodeVarUint(sink, this.GasPerKBForRead)
}

func (this *FsSetting) Deserialization(source *common.ZeroCopySource) error {
	var err error
	this.FsGasPrice, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.GasPerKBForStore, err = utils.DecodeVarUint(source)
	if err != nil {
		return err
	}
	this.GasPerKBForRead, err = utils.DecodeVarUint(source)
	return err
}
