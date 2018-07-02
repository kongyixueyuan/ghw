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
	fmt.Println("Run方法")
	isValidArgs()
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
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

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
//增加数据
func  (cli *CLI) addBlock(data string)  {
	cli.addBlock(data)
	fmt.Println("区块增加成功!")
}

func (cli *CLI) printChain() {
	//迭代打印区块信息
	cli.BlockChain.PrintChainsUseIterator()
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
