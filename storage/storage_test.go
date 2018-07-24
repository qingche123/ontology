package storage

import (
	"testing"
	"fmt"
	"bytes"
	"math/rand"
)

func randString() string {
	alpha := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	l := rand.Intn(10) + 2

	var s string
	for i := 0; i < l; i++ {
		s += string([]byte{alpha[rand.Intn(len(alpha))]})
	}
	return s
}

func TestStorage(t *testing.T) {
	password := "testkey"

	sc := NewStorageService()
	err := sc.StartStoreService()
	if err != nil {
		fmt.Printf(err.Error())
	}

	for i := 0; i < 100; i++  {
		data := randString()
		hash, err := sc.EncryptAndAdd([]byte(data), password, AES)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("File Hash: %s\n", hash)
		decData, err := sc.GetAndDecrypt(hash, password)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if bytes.Compare(decData, []byte(data)) == 0 {
			fmt.Println("GetAndDecrypt Success")
		} else {
			fmt.Println("GetAndDecrypt Failed")
		}
	}

	return
}