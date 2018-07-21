package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)

type Ghw_Block struct {
	Ghw_TimeStamp     int64
	Ghw_Transactions   []*Ghw_Transaction
	Ghw_PrevBlockHash []byte
	Ghw_Hash          []byte
	Ghw_Nonce         int
	Ghw_Height        int
}
// 生成新的区块
func Ghw_NewBlock(transactions []*Ghw_Transaction, prevBlockHash []byte, height int) *Ghw_Block {
	// 生成新的区块对象
	block := &Ghw_Block{
		time.Now().Unix(),
		transactions,
		prevBlockHash,
		[]byte{},
		0,
		height,
	}
	// 挖矿

	pow := Ghw_NewProofOfWork(block)
	nonce,hash :=pow.Ghw_Run()

	block.Ghw_Nonce = nonce
	block.Ghw_Hash = hash[:]

	return block

}

// 将交易进行hash
func (b Ghw_Block) Ghw_HashTransactions() []byte {
	var transactions [][]byte
	// 获取交易真实内容
	for _,tx := range b.Ghw_Transactions{
		transactions = append(transactions,tx.Ghw_Serialize())
	}
	//txHash := sha256.Sum256(bytes.Join(transactions,[]byte{}))
	mTree := Ghw_NewMerkelTree(transactions)
	return mTree.Ghw_RootNode.Ghw_Data
}
// 新建创世区块
func Ghw_NewGenesisBlock(coinbase *Ghw_Transaction) *Ghw_Block  {
	return Ghw_NewBlock([]*Ghw_Transaction{coinbase},[]byte{},1)
}

// 序列化
func (b *Ghw_Block) Ghw_Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func Ghw_DeserializeBlock(d []byte) *Ghw_Block {
	var block Ghw_Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
// 打印区块内容
func (block Ghw_Block) String()  {
	fmt.Println("\n==============")
	fmt.Printf("Height:\t%d\n", block.Ghw_Height)
	fmt.Printf("PrevBlockHash:\t%x\n", block.Ghw_PrevBlockHash)
	fmt.Printf("Timestamp:\t%s\n", time.Unix(block.Ghw_TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
	fmt.Printf("Hash:\t%x\n", block.Ghw_Hash)
	fmt.Printf("Nonce:\t%d\n", block.Ghw_Nonce)
	fmt.Println("Txs:")

	for _, tx := range block.Ghw_Transactions {
		tx.String()
	}
	fmt.Println("==============")
}
