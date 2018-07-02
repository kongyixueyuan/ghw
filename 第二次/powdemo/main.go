package main

import (
	"block/powdemo/blc"
)

func main()  {
/*	genesisBlock := blc.CreateGenesisBlock("Geneis block")
	fmt.Println()
	fmt.Println("创世区块")
	fmt.Println(genesisBlock)
	block:=blc.NewBlock("Geneis block",1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	fmt.Print(block)*/

	//创建带有创世区块的区块链
/*	genesisBlockChain :=blc.CreateBlockChainWithGenesisBlock()
	fmt.Println(genesisBlockChain)
	fmt.Println(genesisBlockChain.Blocks)*/

	//创建新的区块
	blockChain :=blc.CreateBlockChainWithGenesisBlock()
	//fmt.Println(blockChain)
	//fmt.Println(blockChain.Blocks)

	//添加新的区块1
	blockChain.AddBlockToBlockChain("转账100 到张三",
		blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,
		blockChain.Blocks[len(blockChain.Blocks)-1].Hash,
			)
	//添加新的区块2
	blockChain.AddBlockToBlockChain("转账200 到李四",
		blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,
		blockChain.Blocks[len(blockChain.Blocks)-1].Hash,
	)
	//添加新的区块3
	blockChain.AddBlockToBlockChain("转账200 到王五",
		blockChain.Blocks[len(blockChain.Blocks)-1].Height+1,
		blockChain.Blocks[len(blockChain.Blocks)-1].Hash,
	)
/*	fmt.Println("区块。。。们")
	fmt.Println(blockChain.Blocks)
	fmt.Println("第三个区块")
	fmt.Println(blockChain.Blocks[2])*/






}



