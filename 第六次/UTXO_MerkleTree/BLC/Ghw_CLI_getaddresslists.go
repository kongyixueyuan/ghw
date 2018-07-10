package BLC

import "fmt"

// 打印所有的钱包地址
func (cli *CLI) ghw_addressLists()  {

	fmt.Println("打印所有的钱包地址:")

	wallets,_ := Ghw_NewWallets()

	for address,_ := range wallets.Ghw_WalletsMap {

		fmt.Println(address)
	}
}