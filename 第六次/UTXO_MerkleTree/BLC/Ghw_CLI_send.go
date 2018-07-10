package BLC

import (
	"fmt"
	"os"
)

// 转账
func (cli *CLI) ghw_send(from []string,to []string,amount []string)  {

	if Ghw_DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := Ghw_BlockchainObject()

	defer blockchain.Ghw_DB.Close()

	blockchain.Ghw_MineNewBlock(from,to,amount)

	utxoSet := &UTXOSet{blockchain}

	//转账成功以后，需要更新一下
	utxoSet.Ghw_Update()

}

