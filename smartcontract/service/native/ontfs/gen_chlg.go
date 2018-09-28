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
	"github.com/daseinio/dasein-go-PoR/PoR"
	"github.com/ontio/ontology/common"
)

func GenChallenge(hash common.Uint256, fileBlockNum, proveNum uint32) []PoR.Challenge {
	var blockNumPerPart, blockNumLastPart, blockNumOfPart uint32

	if fileBlockNum <= 3 {
		blockNumPerPart = 1
		blockNumLastPart = 1
		blockNumOfPart = 1
		proveNum = fileBlockNum
	} else if fileBlockNum > 3 && fileBlockNum < proveNum {
		proveNum = (fileBlockNum + 3) / 3
		blockNumPerPart = (fileBlockNum / proveNum) + 1
		blockNumLastPart = fileBlockNum % blockNumPerPart
		blockNumOfPart = blockNumPerPart
	} else {
		blockNumPerPart = fileBlockNum / (proveNum - 1)
		blockNumLastPart = fileBlockNum % (proveNum - 1)
		blockNumOfPart = blockNumPerPart
	}

	challenge := make([]PoR.Challenge, proveNum)
	blockHash := hash.ToArray()

	var hashIndex = 0
	var i uint32

	for i = 1; i <= proveNum; i++ {
		if i == proveNum && blockNumLastPart != 0{
			blockNumOfPart = blockNumLastPart
		}

		challenge[i-1].Index = (uint32(blockHash[hashIndex]) +1)% blockNumOfPart + i*blockNumPerPart
		challenge[i-1].Rand = uint32(blockHash[hashIndex]) + 1
		hashIndex++
		hashIndex = hashIndex % common.UINT256_SIZE
	}
	return challenge
}
