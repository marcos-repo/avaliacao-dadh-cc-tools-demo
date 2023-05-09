package main

import (
	"github.com/goledgerdev/cc-tools-demo/chaincode/txdefs"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

var txList = []tx.Transaction{
	txdefs.CriarProprietario,
	txdefs.CriarToken,
	txdefs.TransferirToken,
	txdefs.ContabilizarTokens,
}
