package BLC

import "fmt"

func (cli *Ghw_CLI) ghw_reindexUTXO(nodeID string)  {
	bc := Ghw_NewBlockchain(nodeID);
	defer bc.ghw_db.Close()
	utxoset := Ghw_UTXOSet{bc}
	utxoset.Ghw_Reindex()
	fmt.Println("重建成功")
}
