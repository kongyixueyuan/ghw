package BLC

import "fmt"

// 重置方法
func (cli *CLI) Ghw_TestMethod()  {

	fmt.Println("TestMethod")

	blockchain := Ghw_BlockchainObject()

	defer blockchain.Ghw_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	utxoSet.Ghw_ResetUTXOSet()

	//fmt.Println(blockchain.FindUTXOMap())
}
