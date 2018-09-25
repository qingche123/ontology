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
	"bytes"
	"fmt"
	"github.com/ontio/ontology/common"
	"testing"
)

func TestFileInfo_Serialize(t *testing.T) {
	fileHashStr := []byte("QmevhnWdtmz89BMXuuX5pSY2uZtqKLz7frJsrCojT5kmb6")

	var addr common.Address
	copy(addr[:], fileHashStr[0:20])
	fileInfo := FileInfo{
		FileHash:       fileHashStr,
		UserAddr:       addr,
		KeeyHours:      1,
		FileBlockNum:   2,
		FIleBlockSize:  3,
		ChallengeRate:  4,
		ChallengeTimes: 5,
		CopyNum:        6,
		Deposit:        7,
		FileProveParam: nil,
		ProveBlockNum:  8,
		BlockHeight:    9,
	}
	b := new(bytes.Buffer)
	err := fileInfo.Serialize(b)
	if err != nil {
		t.Error(err.Error())
	}
	var fileInfo2 FileInfo
	err = fileInfo2.Deserialize(b)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(string(fileInfo2.FileHash[:]))
	fmt.Println(string(fileInfo2.UserAddr[:]))
	fmt.Println(fileInfo2.KeeyHours)
	fmt.Println(fileInfo2.FileBlockNum)
	fmt.Println(fileInfo2.FIleBlockSize)
	fmt.Println(fileInfo2.ChallengeRate)
	fmt.Println(fileInfo2.ChallengeTimes)
	fmt.Println(fileInfo2.CopyNum)
	fmt.Println(fileInfo2.Deposit)
	fmt.Println(fileInfo2.FileProveParam)
	fmt.Println(fileInfo2.ProveBlockNum)
	fmt.Println(fileInfo2.BlockHeight)
}
