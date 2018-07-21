package BLC

import (
	"log"
	"fmt"
)

func (cli *Ghw_CLI) ghw_listAddrsss(nodeID string)  {
	wallets,err := Ghw_NewWallets(nodeID)

	if err!=nil{
		log.Panic(err)
	}
	addresses := wallets.Ghw_GetAddresses()

	for _,address := range addresses{
		fmt.Printf("%s\n",address)
	}
}
