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
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")
	genesisBlockData := addGenesisBlockCmd.String("data","Genesis block data......","创世区块交易数据......")
	switch os.Args[1] {
	case "addGenesisBlock":
		fmt.Println("增加创世区块")
		err := addGenesisBlockCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "addBlock":
		fmt.Println("增加区块")
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "printChain":
		fmt.Println("查看区块")
	err := printChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}
	if addGenesisBlockCmd.Parsed() {
		if *genesisBlockData == "" {
			fmt.Println("交易数据不能为空......")
			addGenesisBlockCmd.Usage()
			os.Exit(1)
		}
		//创建创世区块
		CreateBlockChainWithGenesisBlock(*genesisBlockData)
	}

	if printChainCmd.Parsed() {
		fmt.Println(cli)
		cli.printChain()

	}
}

func (cli *CLI) CreateBlockChainWithGenesisBlock(data string)  {

	CreateBlockChainWithGenesisBlock(data)
	fmt.Println(data)

}

//增加数据
func  (cli *CLI) addBlock(data string)  {
	if IsDBExists(dbname) == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	fmt.Println("返回区块链对象")
	fmt.Println(blockchain)
	defer blockchain.DB.Close()
	fmt.Println(blockchain)
	blockchain.AddBlockToBlockChain(data)
/*	cli.BlockChain.AddBlockToBlockChain(data)
//	cli.addBlock(data)
	fmt.Println("区块增加成功!")*/
}

func (cli *CLI) printChain() {
	fmt.Println("开始查看区块数据集")
	if IsDBExists(dbname) == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	fmt.Println("返回区块链对象")
	fmt.Println(blockchain)
	defer blockchain.DB.Close()
	fmt.Println(blockchain)
	blockchain.PrintChainsUseIterator()

	//迭代打印区块信息

	//cli.BlockChain.PrintChainsUseIterator()
	/*bci := cli.bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.IsValild()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}*/
}
func isValidArgs()  {
	fmt.Printf("参数的长度%s",len(os.Args))
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage()  {
	fmt.Println("\t usage:")
	fmt.Println("\t addGenesisBlock -data -- 增加创世区块.")
	fmt.Println("\t addBlock -data DATA -交易数据")
	fmt.Println("\t printChain -- 输出区块信息")
}
