package blc

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

//定义区块结构
type Block struct {
	//1,区块高度,区块的编号
	Height int64
	//2,上一个区块的hash
	PrevBlockHash []byte
	//3,交易数据
	//Data []byte
	Txs []*Transaction
	//4,时间戳
	Timestamp int64
	//5, hash
	Hash []byte
	//6,随机数
	Nonce int64


}

//需要将Txs 转换为 []byte
func (block *Block) HashTransactions() []byte {
 	var txHashes [][]byte
 	var txHash [32]byte
 	for _,tx:= range block.Txs{
		txHashes = append(txHashes,tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	return txHash[:]
}


//序列号
func (block *Block) Serialize()  []byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化，区块序列化为字节数组
func DeserializeBlock(blockBytes []byte)  *Block{
	var block  Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//1,创建新的区块，
func NewBlock(txs []*Transaction,height int64,prevBlockHash []byte) *Block{
	//创建区块
	 block := &Block{Height: height, PrevBlockHash: prevBlockHash, Txs: txs,Timestamp:time.Now().Unix(), Hash: nil,Nonce:0}
	 //设置当前区块的hash == >调用工作量证明的方法，并且返回有限的Nonce  hash值
	// block.setHash()
	pow := NewProofOfWork(block)
	hash,nonce :=pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	 return block
}

// 生产创世区块
func CreateGenesisBlock(txs []*Transaction) *Block  {
	//高度，可知1，上一个区块的hash值
	return NewBlock(txs,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

func main() {
}
