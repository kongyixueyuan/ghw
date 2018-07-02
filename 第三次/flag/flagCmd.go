package main

import (
	"flag"
	"os"
	"log"
	"fmt"
)

func main()  {
	addBlockCmd := flag.NewFlagSet("addBlock",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain",flag.ExitOnError)
	flagAddBlockData := addBlockCmd.String("data","默认值","交易数据")
	isValidArgs()
	switch os.Args[1] {
	case "addBlock":
		err:=addBlockCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "printChain":
		err:=printChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)//退出
	}
	if addBlockCmd.Parsed(){
		if *flagAddBlockData == "*"{
			printUsage()
			os.Exit(1)
		}
		fmt.Println(*flagAddBlockData)
	}
	if printChainCmd.Parsed(){
		fmt.Println("输出所有信息")
	}
}
func isValidArgs()  {
	if len(os.Args) <2{
		printUsage()
		os.Exit(1)
	}
}
func printUsage()  {
	fmt.Println("usage:")
	fmt.Println("\addBlock -data DATA -交易数据")
	fmt.Println("printChain -- 输出区块信息")
}
