package blc

type BlockChain struct {
	Blocks []*Block //存储区块链
}

//增加区块到区块链中
func (blc *BlockChain) AddBlockToBlockChain(data string,height int64,prehash []byte)  {
	//增加区块
	newBlock := NewBlock(data,height,prehash)
	//新区块加入链中
	blc.Blocks = append(blc.Blocks,newBlock)
}

//创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock() *BlockChain  {
	genesisBlock := CreateGenesisBlock("genesis block")
	//返回区块链对象
	return &BlockChain{[]*Block{genesisBlock}}

}


