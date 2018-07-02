package blc

import (
	"time"

)

//定义区块结构
type Block struct {
	//1,区块高度,区块的编号
	Height int64
	//2,上一个区块的hash
	PrevBlockHash []byte
	//3,交易数据
	Data []byte
	//4,时间戳
	Timestamp int64
	//5, hash
	Hash []byte
	//6,随机数
	Nonce int64
}

/*func (block *Block) setHash()  {
	//1,Height ==> 字节数组
	heightBytes := IntToHex(block.Height)
	//2,Timestamp == >字节数组
	timeString := strconv.FormatInt(block.Timestamp,2)
	timestamp := []byte(timeString)
	//3,属性数据拼接
	blockBytes := bytes.Join([][]byte{heightBytes,block.PrevBlockHash,block.Data,timestamp,block.Hash},[]byte{})
	hash := sha256.Sum256(blockBytes)

	block.Hash = hash[:]
}*/

//1,创建新的区块
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {

	//创建区块
	 block := &Block{Height: height, PrevBlockHash: prevBlockHash, Data: []byte(data),Timestamp:time.Now().Unix(), Hash: nil,Nonce:0}

	 //设置当前区块的hash == >调用工作量证明的方法，并且返回有限的Nonce  hash值
	// block.setHash()
	pow := NewProofOfWork(block)

	hash,nonce :=pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce;

	 return block
}

// 生产创世区块
func CreateGenesisBlock(data string) *Block  {
	//高度，可知1，上一个区块的hash值
	return NewBlock(data,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}
func main() {
}
