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
	Url string
}

func NewStorageService(url string) *StorageContext {
	return &StorageContext{Sh:nil, Url:url}
}

func (sc *StorageContext)StartStoreService() error {
	sc.Sh = storageSdk.NewShell(sc.Url)
	if !sc.Sh.IsUp() {
		return errors.New("storage service start failed, connection error")
	}
	return nil
}

func (sc *StorageContext)EncryptAndAdd(data []byte, password string, alg EncryptScheme) (string, error) {
	return sc.Sh.EncryptAndAdd(data, password, storageSdk.AES)
}

func (sc *StorageContext)GetAndDecrypt(hash string, password string) ([]byte, error) {
	return sc.Sh.GetAndDecrypt(hash, password)
}

func (sc *StorageContext)AddData(data []byte) (string, error) {
	return sc.Sh.AddData(data)
}

func (sc *StorageContext)GetData(hash string) ([]byte, error) {
	return sc.Sh.GetData(hash)
}
