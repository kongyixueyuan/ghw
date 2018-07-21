package BLC

func (cli *Ghw_CLI) ghw_printchain(nodeID string)  {
	Ghw_NewBlockchain(nodeID).Ghw_Printchain()
}
