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
	FS_SET              = "FsSet"
	FS_SET_INIT         = "FsSettingInit"
	FS_GETSETTING       = "FsGetSetting"
	FS_NODE_REGISTER    = "FsNodeRegister"
	FS_NODE_QUERY       = "FsNodeQuery"
	FS_NODE_UPDATE      = "FsNodeUpdate"
	FS_NODE_CANCEL      = "FsNodeCancel"
	FS_GET_NODE_LIST    = "FsGetNodeList"
	FS_STORE_FILE       = "FsStoreFile"
)

const (
	ONTFS_SETTING     = "ontfssetting"
	ONTFS_NODE_INFO   = "ontfsnodeInfo"
	ONTFS_NODE_SET   = "ontfsnodeset"
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