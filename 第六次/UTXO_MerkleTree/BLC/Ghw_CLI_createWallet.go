package BLC

import "fmt"

//创建钱包
func (cli *CLI) ghw_createWallet()  {

	//钱包集
	wallets,_ := Ghw_NewWallets()
	//创建钱包
	wallets.Ghw_CreateNewWallet()

	fmt.Println(len(wallets.Ghw_WalletsMap))
}
