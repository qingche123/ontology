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
	"github.com/ontio/ontology/common"
)

const (
	FS_SET                     = "FsSet"
	FS_SET_INIT                = "FsSettingInit"
	FS_GETSETTING              = "FsGetSetting"
	FS_NODE_REGISTER           = "FsNodeRegister"
	FS_NODE_QUERY              = "FsNodeQuery"
	FS_NODE_UPDATE             = "FsNodeUpdate"
	FS_NODE_CANCEL             = "FsNodeCancel"
	FS_GET_NODE_LIST           = "FsGetNodeList"
	FS_STORE_FILE              = "FsStoreFile"
	FS_GET_FILE_INFO           = "FsGetFileInfo"
	FS_NODE_WITH_DRAW_PROFIT   = "FsNodeWithDrawProfit"
	FS_FILE_PROVE              = "FsFileProve"
	FS_GET_FILE_PROVE_DETAILS  = "FsGetFileProveDetails"
	FS_READ_FILE_PLEDGE        = "FsReadFilePledge"
	FS_FILE_READ_PROFIT_SETTLE = "FsFileReadProfitSettle"
)

const (
	ONTFS_SETTING          = "ontfssetting"
	ONTFS_NODE_INFO        = "ontfsnodeInfo"
	ONTFS_NODE_SET         = "ontfsnodeset"
	ONTFS_FILE_PROVE       = "ontfsfileprove"
	ONTFS_FILE_READ_PLEDGE = "ontfsfileprove"
)

func GenFsSettingKey(contract common.Address) []byte {
	return append(contract[:], ONTFS_SETTING...)
}

func GenFsNodeInfoKey(contract common.Address, walletAddr common.Address) []byte {
	key := append(contract[:], ONTFS_NODE_INFO...)
	return append(key, walletAddr[:]...)
}

func GenFsNodeSetKey(contract common.Address) []byte {
	return append(contract[:], ONTFS_NODE_SET...)
}

func GenFsProveDetailsKey(contract common.Address, fileHash []byte) []byte {
	key := append(contract[:], ONTFS_FILE_PROVE...)
	return append(key, fileHash[:]...)
}

func GenFsFileReadPledgeKey(contract common.Address, fileHash []byte) []byte {
	key := append(contract[:], ONTFS_FILE_READ_PLEDGE...)
	return append(key, fileHash[:]...)
}
