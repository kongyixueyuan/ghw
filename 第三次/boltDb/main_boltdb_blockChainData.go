package main

import "./blc"

func main()  {
	//创建区块
	blockChain := blc.CreateBlockChainWithGenesisBlock()
	//关闭数据库
	defer blockChain.DB.Close()
	//增加新的区块内容
	blockChain.AddBlockToBlockChain("张三转账给李四100")
	blockChain.AddBlockToBlockChain("李四转账给王五1000")
	blockChain.AddBlockToBlockChain("李四请假了")
	blockChain.AddBlockToBlockChain("赵柳喝酒了")
	blockChain.AddBlockToBlockChain("赵柳喝酒了222")

	//打印区块信息
	//blockChain.PrintChains()
	//使用迭代器打印
	blockChain.PrintChainsUseIterator()

}
