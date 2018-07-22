package BLC

import "fmt"

func (cli *Ghw_CLI) ghw_createWallet(nodeID string) {
	wallets, _ := Ghw_NewWallets(nodeID)
	address := wallets.Ghw_NewWallet().Ghw_GetAddress()
	wallets.Ghw_SaveToFile(nodeID)
	fmt.Printf("钱包地址：%s\n", address)

}
