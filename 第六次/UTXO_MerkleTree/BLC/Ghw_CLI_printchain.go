package BLC

import (
	"fmt"
	"os"
)

// 打印区块信息
func (cli *CLI) ghw_printchain()  {

	if Ghw_DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := Ghw_BlockchainObject()

	defer blockchain.Ghw_DB.Close()

	blockchain.Ghw_Printchain()

}