package blc

import (
	"flag"
	"os"
	"fmt"
	"log"
)

type CLI struct {
	BlockChain *BlockChain
}
func (cli *CLI) Run() {
	isValidArgs()//判断用户数的参数
	addGenesisBlockCmd := flag.NewFlagSet("addGenesisBlock", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)

	sendBlockFrom := sendBlockCmd.String("from", "", "源地址")
	sendBlockTo := sendBlockCmd.String("to", "", "目的地地址")
	sendBlockAmount := sendBlockCmd.String("amount", "", "转账金额")
	getbalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	flagCreateBlockchainWithAddress := addGenesisBlockCmd.String("address","","创建创世区块的地址")
	getbalanceWithAdress := getbalanceCmd.String("address","","要查询某一个账号的余额.......")

	switch os.Args[1] {
	case "addGenesisBlock":
		fmt.Println("增加创世区块")
		err := addGenesisBlockCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "send":
		fmt.Println("交易发生 增加新的区块")
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "printChain":
		fmt.Println("查看区块")
	err := printChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	//转账发生
	if sendBlockCmd.Parsed() {
		if *sendBlockFrom == "" ||*sendBlockTo == ""||*sendBlockAmount == "" {
			printUsage()
			os.Exit(1)
		}
		//转账
		cli.send(JSONToArray(*sendBlockFrom),JSONToArray(*sendBlockTo),JSONToArray(*sendBlockAmount))
	}
	if addGenesisBlockCmd.Parsed() {
		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("地址不能为空")
			addGenesisBlockCmd.Usage()
			os.Exit(1)
		}
		//创建创世区块
		cli.CreateBlockChainWithGenesisBlock( *flagCreateBlockchainWithAddress)
	}
	if printChainCmd.Parsed() {
		fmt.Println(cli)
		cli.printChain()
	}
	if getbalanceCmd.Parsed() {
		if *getbalanceWithAdress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getbalanceWithAdress)
	}
}

func (cli *CLI) send(from []string,to []string,amount []string)  {
	if IsDBExists(dbname) == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	blockChain:=BlockchainObject()
	defer  blockChain.DB.Close()
	blockChain.MineNewBlock(from,to,amount)

}

func (cli *CLI) CreateBlockChainWithGenesisBlock(address string)  {
	blockChain :=CreateBlockChainWithGenesisBlock(address)
	defer  blockChain.DB.Close()
}
func (cli *CLI) printChain() {
	fmt.Println("开始查看区块数据集")
	if IsDBExists(dbname) == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	blockchain.PrintChainsUseIterator()
}

func (cli *CLI) getBalance(address string)  {
	fmt.Println("地址：" + address)
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	amount := blockchain.GetBalance(address)
	fmt.Printf("%s一共有%d个Token\n",address,amount)
}

func isValidArgs()  {
	fmt.Printf("参数的长度%s",len(os.Args))
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
//main send -from "[\"gaohongwei\"]" -to "[\"zhangsan\"]" -amount "[\"2\"]
//main send -from "[\"ghw\",\"xiaoming\"]" -to "[\"zhangsan\",\"lisi\"]" -amount "[\"2\",\"2\"]"
func printUsage()  {
	fmt.Println("\t usage:")
	fmt.Println("\t addGenesisBlock -address -- 增加创世区块.")
	fmt.Println("\t send -from from -to to -amount amount")
	fmt.Println("\t printChain -- 输出区块信息")
	fmt.Println("\tgetbalance -address -- 余额查询.")
}
