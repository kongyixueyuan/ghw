package BLC

func (cli *Ghw_CLI) ghw_send(from []string, to []string, amount []string,nodeID string, mineNow bool) {
	bc := Ghw_NewBlockchain(nodeID)
	defer bc.ghw_db.Close()
	bc.MineNewBlock(from, to, amount,nodeID, mineNow)
}