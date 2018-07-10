package BLC

import "fmt"

// 先用它去查询余额
func (cli *CLI) ghw_getBalance(address string)  {

	fmt.Println("地址：" + address)

	blockchain := Ghw_BlockchainObject()

	defer blockchain.Ghw_DB.Close()

	utxoSet := &UTXOSet{blockchain}

	amount := utxoSet.Ghw_GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n",address,amount)

}
