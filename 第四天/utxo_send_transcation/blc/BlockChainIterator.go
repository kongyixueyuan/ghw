package blc


import (
	"block/bolt"
	"log"

	"fmt"
	"time"
	"math/big"
)

type BlockChainIterator struct {
	CurrentHash []byte
	DB *bolt.DB
}

func (blockChainIterator *BlockChainIterator) Next() *Block {
	var block *Block
	err := blockChainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil{
			currentBlock := b.Get(blockChainIterator.CurrentHash)
			block = DeserializeBlock(currentBlock)
			blockChainIterator.CurrentHash =block.PrevBlockHash
		}
		return nil
	})
	if  err != nil{
		log.Panic(err)
	}
	blockChainIterator.CurrentHash = block.PrevBlockHash
	return block
}
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return  &BlockChainIterator{bc.Tip,bc.DB}
}

func (blc *BlockChain) PrintChainsUseIterator(){
	fmt.Println("PrintchainPrintchainPrintchainPrintchain")
	blockChainIterator  := blc.Iterator()
	for {
		block := blockChainIterator.Next()
		fmt.Printf("区块的高度%d\n",block.Height)
		fmt.Printf("区块的hash%x\n",block.Hash)
		fmt.Printf("区块的数据%s\n",block.HashTransactions())
		fmt.Printf("区块的时间戳%d\n",time.Unix(block.Timestamp,0).Format("2006-01-02 03:04:25 PM"))
		fmt.Printf("区块的随机数%d\n",block.Nonce)
		fmt.Printf("上一个区块的高度%x\n",block.PrevBlockHash)
		fmt.Println("Txs:")
		//遍历transaction对象
		for _, tx := range block.Txs{
			fmt.Printf("%x",tx.TxHash)
			//遍历输出票
			fmt.Println("vouts")
			for _, out := range tx.Vouts  {
				fmt.Println(out.Money)
				fmt.Println(out.ScriptPublicKey)
			}
			//遍历输入票
			fmt.Println("vins")
			for _, in := range tx.Vins  {
				fmt.Printf("%s\n",in.TxHash)
				fmt.Printf("%d\n",in.Vout)
				fmt.Printf("%s\n",in.ScriptSig)
			}
		}
		
		var hashInt  big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) ==0{
			break
		}
	}
}