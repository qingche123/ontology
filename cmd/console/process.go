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

package console

import (
	"fmt"
	"strings"
	"reflect"
	"errors"

	"github.com/ontio/ontology-go-sdk/rpc"
	"github.com/ontio/ontology-go-sdk/wallet"
	"github.com/ontio/ontology-go-sdk/common"
)

var process * Process
var keyWords []string

func Call(statement string, cons *Console) error {
	var err error
	process.Cons = cons

	sm := strings.Split(statement, ".")
	if true {
		// todo
	} else {
		fmt.Println("Unknown statement: ", sm[0], " ", sm[1])
		err = errors.New(fmt.Sprint("Unknown statement: ", sm[0], " ", sm[1]))
	}
	return err
}

type Process struct {
	Cons *Console
}

func NewProcess(cons *Console) *Process{
	collectKeyWords(cons)
	return &Process{Cons:cons}
}

func collectKeyWords(cons *Console){

}

func CompleteKeywords(line string) []string {
	var rightKeyWords []string

	return rightKeyWords
}
