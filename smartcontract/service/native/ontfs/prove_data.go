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

	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

type ProveData struct {
	MultiRes []byte
	AddRes   []byte
}

func (this *ProveData) Serialize(w io.Writer) error {
	if err := utils.WriteBytes(w, this.MultiRes); err != nil {
		return fmt.Errorf("[ProveData] serialize from error:%v", err)
	}
	if err := utils.WriteBytes(w, this.AddRes); err != nil {
		return fmt.Errorf("[ProveData] serialize from error:%v", err)
	}
	return nil
}

func (this *ProveData) Deserialize(r io.Reader) error {
	var err error
	if this.MultiRes, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveData] deserialize from error:%v", err)
	}
	if this.AddRes, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveData] deserialize from error:%v", err)
	}
	return nil
}

type ProveParam struct {
	G        []byte
	G0       []byte
	PubKey   []byte
	FileId   []byte
	R        []byte
}

func (this *ProveParam) Serialize(w io.Writer) error {
	if err := utils.WriteBytes(w, this.G); err != nil {
		return fmt.Errorf("[ProveParam] serialize from error:%v", err)
	}
	if err := utils.WriteBytes(w, this.G0); err != nil {
		return fmt.Errorf("[ProveParam] serialize from error:%v", err)
	}
	if err := utils.WriteBytes(w, this.PubKey); err != nil {
		return fmt.Errorf("[ProveParam] serialize from error:%v", err)
	}
	if err := utils.WriteBytes(w, this.FileId); err != nil {
		return fmt.Errorf("[ProveParam] serialize from error:%v", err)
	}
	if err := utils.WriteBytes(w, this.R); err != nil {
		return fmt.Errorf("[ProveParam] serialize from error:%v", err)
	}
	return nil
}

func (this *ProveParam) Deserialize(r io.Reader) error {
	var err error
	if this.G, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] deserialize from error:%v", err)
	}
	if this.G0, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] deserialize from error:%v", err)
	}
	if this.PubKey, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] deserialize from error:%v", err)
	}
	if this.FileId, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] deserialize from error:%v", err)
	}
	if this.R, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] deserialize from error:%v", err)
	}
	return nil
}
