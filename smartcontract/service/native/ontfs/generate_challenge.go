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
	if fileBlockNum <= 3 {
		proveNum = 3
	} else if fileBlockNum > 3 && fileBlockNum < proveNum {
		proveNum = fileBlockNum / 3
	}

	challenge := make([]PoR.Challenge, proveNum)
	blockHash := hash.ToArray()

	blockNumPerPart := fileBlockNum / (proveNum - 1)
	blockNumLastPart := fileBlockNum % (proveNum - 1)
	blockNumOfPart := blockNumPerPart

	var hashIndex = 0
	var i uint32
	for i = 0; i < proveNum; i++ {
		if i == proveNum-1 {
			blockNumOfPart = blockNumLastPart
		}
		challenge[i].Index = uint32(blockHash[hashIndex]) % blockNumOfPart + i*blockNumPerPart
		challenge[i].Rand = uint32(blockHash[hashIndex] & byte(i))

		hashIndex++
		hashIndex = hashIndex % common.UINT256_SIZE
	}
	return challenge
}
