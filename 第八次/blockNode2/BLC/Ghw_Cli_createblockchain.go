package BLC

import "log"

func (cli *Ghw_CLI) ghw_createblockchain(address string,nodeID string)  {
	//验证地址是否有效
	if !Ghw_ValidateAddress(address){
		log.Panic("地址无效")
	}
	bc := Ghw_CreateBlockchain(address,nodeID)
	defer bc.ghw_db.Close()

	// 生成UTXOSet数据库
	UTXOSet := Ghw_UTXOSet{bc}
	UTXOSet.Ghw_Reindex()
}
