package main

import (
	"./blc"
	"fmt"
)
func main()  {
	//创建区块
	blockChain := blc.CreateBlockChainWithGenesisBlock()
	//关闭数据库
	defer blockChain.DB.Close()
	//cli:=&blc.CLI{blockChain}
	fmt.Println("开始执行")
	//cli.Run()


}
