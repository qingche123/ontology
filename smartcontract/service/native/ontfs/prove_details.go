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

type FsProveDetails struct {
	CopyNum        uint64
	ProveDetailNum uint64
	ProveDetails   []ProveDetail
}

type ProveDetail struct {
	NodeAddr   []byte
	WalletAddr common.Address
	ProveTimes uint64
}

func (this *ProveDetail) Serialize(w io.Writer) error {
	if err := utils.WriteBytes(w, this.NodeAddr[:]); err != nil {
		return fmt.Errorf("[ProveNode] serialize from error:%v", err)
	}
	if err := utils.WriteAddress(w, this.WalletAddr); err != nil {
		return fmt.Errorf("[ProveNode] serialize from error:%v", err)
	}
	if err := utils.WriteVarUint(w, this.ProveTimes); err != nil {
		return fmt.Errorf("[ProveNode] serialize from error:%v", err)
	}
	return nil
}

func (this *ProveDetail) Deserialize(r io.Reader) error {
	var err error
	if this.NodeAddr, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveNode] deserialize from error:%v", err)
	}
	if this.WalletAddr, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("[ProveNode] deserialize from error:%v", err)
	}
	if this.ProveTimes, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[ProveNode] deserialize from error:%v", err)
	}
	return nil
}

func (this *FsProveDetails) Serialize(w io.Writer) error {
	var err error
	if err = utils.WriteVarUint(w, this.CopyNum); err != nil {
		return fmt.Errorf("[ProveDetail] serialize from error:%v", err)
	}
	if err = utils.WriteVarUint(w, this.ProveDetailNum); err != nil {
		return fmt.Errorf("[ProveDetail] serialize from error:%v", err)
	}
	for _, v := range this.ProveDetails {
		err = v.Serialize(w)
		if err != nil {
			return fmt.Errorf("[ProveDetail] serialize from error:%v", err)
		}
	}
	return nil
}

func (this *FsProveDetails) Deserialize(r io.Reader) error {
	var err error
	var tmpProveDetail ProveDetail
	if this.CopyNum, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[ProveDetail] deserialize from error:%v", err)
	}
	if this.ProveDetailNum, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("[ProveDetail] deserialize from error:%v", err)
	}
	for i := 0; uint64(i) < this.ProveDetailNum; i++ {
		if err = tmpProveDetail.Deserialize(r); err != nil {
			return fmt.Errorf("[ProveDetail] deserialize from error:%v", err)
		}
		this.ProveDetails = append(this.ProveDetails, tmpProveDetail)
	}
	return nil
}
