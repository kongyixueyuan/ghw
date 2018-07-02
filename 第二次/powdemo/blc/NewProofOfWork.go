package blc

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

//hash 256
//0000 0000 0000 0000 ......0001 ,前面4位为0
//hash 256 的前面16个0
const targitBit  = 20
type ProofOfWork struct {
	Block *Block //当前要验证的区块
	//diff int64 //难度系数
	target *big.Int//大数据存储，防止数据溢出
}

//3.数据拼接，返回字节数组
func (proofOfWork *ProofOfWork) preparedData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.Block.PrevBlockHash,
			proofOfWork.Block.Data,
			IntToHex(proofOfWork.Block.Timestamp),
			IntToHex(int64(targitBit)),
			IntToHex(int64(nonce)),
			IntToHex(proofOfWork.Block.Height),
		},
		[]byte{},
	)
	return data
}

func (proofOfWork *ProofOfWork) IsValild() bool  {
	//1,proofOfWork.Block.Hash
	//2,proofOfWork.target
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)
	if proofOfWork.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

func (proofOfWork *ProofOfWork) Run() ([]byte,int64) {
	nonce := 0
	var hashInt big.Int //存储新的hash值
	var hash [32]byte


	//3,判断hash的有效姓
	for  {
		//1,把所有的blocks的属性拼接成字节数组
		dataBytes := proofOfWork.preparedData(nonce)

		//2,生产hash
		hash := sha256.Sum256(dataBytes)

		fmt.Printf("\r%x",hash)
		//fmt.Printf("\r%x",hash)
		//存储hash 到hashInt,判断hash 是否小于目标值（hash中target）
		hashInt.SetBytes(hash[:])
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce ++
	}
	return hash[:],int64(nonce)
}
//创建工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork  {
	//1,big.Int 对象1
	//1,创建初始值1位的target
	target := big.NewInt(1)

	//2,左移动 256 - target
	target = target.Lsh(target,256-targitBit)

	return &ProofOfWork{block,target}
}
