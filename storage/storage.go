package storage

import (
	"errors"

	storageSdk "github.com/qingche123/go-ipfs-api"
)

type EncryptScheme byte

const (
	AES EncryptScheme = iota
)

type StorageContext struct {
	Sh *storageSdk.Shell
}

func NewStorageService() *StorageContext {
	return &StorageContext{Sh:nil}
}

func (sc *StorageContext)StartStoreService() error {
	sc.Sh = storageSdk.NewShell("localhost:5001")
	if sc.Sh == nil {
		return errors.New("storage service start failed")
	}
	return nil
}

func (sc *StorageContext)EncryptAndAdd(data []byte, password string, alg EncryptScheme) (string, error) {
	return sc.Sh.EncryptAndAdd(data, password, storageSdk.AES)
}

func (sc *StorageContext)GetAndDecrypt(hash string, password string) ([]byte, error) {
	return sc.Sh.GetAndDecrypt(hash, password)
}


func (sc *StorageContext)Add(data []byte) (string, error) {
	return sc.Sh.Add(data)
}

func (sc *StorageContext)Get(hash string) ([]byte, error) {
	return sc.Sh.Get(hash)
}
