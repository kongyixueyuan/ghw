package BLC

func (cli *Ghw_CLI) ghw_printutxo(nodeID string) {
	bc := Ghw_NewBlockchain(nodeID)
	UTXOSet := Ghw_UTXOSet{bc}
	defer bc.ghw_db.Close()
	UTXOSet.String()
}
