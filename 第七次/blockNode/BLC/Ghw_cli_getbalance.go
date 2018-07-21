package BLC

import (
	"log"
	"fmt"
)

func (cli *Ghw_CLI) ghw_getBalance(address string,nodeID string) {
	if !Ghw_ValidateAddress(address) {
		log.Panic("错误：地址无效")
	}

	bc := Ghw_NewBlockchain(nodeID)
	defer bc.ghw_db.Close()
	UTXOSet := Ghw_UTXOSet{bc}

	balance := UTXOSet.Ghw_GetBalance(address)
	fmt.Printf("地址:%s的余额为：%d\n", address, balance)
}

func (cli *Ghw_CLI) ghw_getBalanceAll(nodeID string) {
	wallets,err := Ghw_NewWallets(nodeID)
	if err!=nil{
		log.Panic(err)
	}
	balances := wallets.Ghw_GetBalanceAll(nodeID)
	for address,balance := range balances{
		fmt.Printf("地址:%s的余额为：%d\n", address, balance)
	}
}