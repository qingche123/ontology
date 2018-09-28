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
		return fmt.Errorf("[ProveData] [MultiRes:%v] serialize from error:%v", this.MultiRes, err)
	}
	if err := utils.WriteBytes(w, this.AddRes); err != nil {
		return fmt.Errorf("[ProveData] [AddRes:%v] serialize from error:%v", this.AddRes, err)
	}
	return nil
}

func (this *ProveData) Deserialize(r io.Reader) error {
	var err error
	if this.MultiRes, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveData] [MultiRes] deserialize from error:%v", err)
	}
	if this.AddRes, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveData] [AddRes] deserialize from error:%v", err)
	}
	return nil
}

type ProveParam struct {
	G        []byte
	G0       []byte
	PubKey   []byte
	FileId   []byte
	R        []byte
	Paring   []byte
}

func (this *ProveParam) Serialize(w io.Writer) error {
	if err := utils.WriteBytes(w, this.G); err != nil {
		return fmt.Errorf("[ProveParam] [G:%v] serialize from error:%v", this.G, err)
	}
	if err := utils.WriteBytes(w, this.G0); err != nil {
		return fmt.Errorf("[ProveParam] [G0:%v] serialize from error:%v", this.G0, err)
	}
	if err := utils.WriteBytes(w, this.PubKey); err != nil {
		return fmt.Errorf("[ProveParam] [PubKey:%v] serialize from error:%v", this.PubKey, err)
	}
	if err := utils.WriteBytes(w, this.FileId); err != nil {
		return fmt.Errorf("[ProveParam] [FileId:%v] serialize from error:%v", this.FileId, err)
	}
	if err := utils.WriteBytes(w, this.R); err != nil {
		return fmt.Errorf("[ProveParam] [R:%v] serialize from error:%v", this.R, err)
	}
	if err := utils.WriteBytes(w, this.Paring); err != nil {
		return fmt.Errorf("[ProveParam] [Paring:%v] serialize from error:%v", this.Paring, err)
	}
	return nil
}

func (this *ProveParam) Deserialize(r io.Reader) error {
	var err error
	if this.G, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] [G] deserialize from error:%v", err)
	}
	if this.G0, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] [G0] deserialize from error:%v", err)
	}
	if this.PubKey, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] [PubKey] deserialize from error:%v", err)
	}
	if this.FileId, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] [FileId] deserialize from error:%v", err)
	}
	if this.R, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] [R] deserialize from error:%v", err)
	}
	if this.Paring, err = utils.ReadBytes(r); err != nil {
		return fmt.Errorf("[ProveParam] [Paring] deserialize from error:%v", err)
	}
	return nil
}
