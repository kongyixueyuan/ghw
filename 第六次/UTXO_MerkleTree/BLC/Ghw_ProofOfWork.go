package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)


//0000 0000 0000 0000 1001 0001 0000 .... 0001

// 256位Hash里面前面至少要有16个零
const ghw_targetBit  = 20

type ProofOfWork struct {
	Ghw_Block *Block // 当前要验证的区块
	ghw_target *big.Int // 大数据存储
}

// 数据拼接，返回字节数组
func (pow *ProofOfWork) ghw_prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Ghw_Block.Ghw_PrevBlockHash,
			pow.Ghw_Block.Ghw_HashTransactions(),
			Ghw_IntToHex(pow.Ghw_Block.Ghw_Timestamp),
			Ghw_IntToHex(int64(ghw_targetBit)),
			Ghw_IntToHex(int64(nonce)),
			Ghw_IntToHex(int64(pow.Ghw_Block.Ghw_Height)),
		},
		[]byte{},
	)

	return data
}

// 开始挖矿
func (proofOfWork *ProofOfWork) Ghw_Run() ([]byte,int64) {

	//1. 将Block的属性拼接成字节数组

	//2. 生成hash
	
	//3. 判断hash有效性，如果满足条件，跳出循环

	nonce := 0

	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte

	for {
		//准备数据
		dataBytes := proofOfWork.ghw_prepareData(nonce)

		// 生成hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x",hash)

		// 将hash存储到hashInt
		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于Block里面的target
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if proofOfWork.ghw_target.Cmp(&hashInt) == 1 {
			break
		}

		nonce = nonce + 1
	}

	return hash[:],int64(nonce)
}


// 创建新的工作量证明对象
func Ghw_NewProofOfWork(block *Block) *ProofOfWork  {

	//1.big.Int对象 1
	// 2
	//0000 0001
	// 8 - 2 = 6
	// 0100 0000  64
	// 0010 0000
	// 0000 0000 0000 0001 0000 0000 0000 0000 0000 0000 .... 0000

	//1. 创建一个初始值为1的target

	target := big.NewInt(1)

	//2. 左移256 - targetBit

	target = target.Lsh(target,256 - ghw_targetBit)

	return &ProofOfWork{block,target}
}






