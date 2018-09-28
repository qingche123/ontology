package ontfs

import (
	"testing"
	"math/rand"
	"fmt"
	"github.com/ontio/ontology/common"

)

func TestGenChallenge(t *testing.T) {
	bt := make([]byte, 32)
	rand.Seed(10)
	rand.Read(bt)

	var hash common.Uint256
	copy(hash[:], bt)

	for fileBlockNum := 1; fileBlockNum < 128; fileBlockNum++  {
		challenge := GenChallenge(hash, uint32(fileBlockNum), 32)
		fmt.Println(challenge)
	}
}