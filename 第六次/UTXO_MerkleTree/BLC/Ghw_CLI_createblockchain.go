package BLC


// 创建创世区块
func (cli *CLI) ghw_createGenesisBlockchain(address string)  {

	blockchain := Ghw_CreateBlockchainWithGenesisBlock(address)
	//打开数据库后需要关闭链接
	defer blockchain.Ghw_DB.Close()

	utxoSet := &UTXOSet{blockchain}
	//将交易保存到文件
	utxoSet.Ghw_ResetUTXOSet()
}
