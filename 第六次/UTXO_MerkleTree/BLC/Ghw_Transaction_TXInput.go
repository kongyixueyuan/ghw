package BLC

import "bytes"

type TXInput struct {
	// 1. 交易的Hash
	Ghw_TxHash      []byte
	// 2. 存储TXOutput在Vout里面的索引
	Ghw_Vout      int

	Ghw_Signature []byte // 数字签名

	Ghw_PublicKey    []byte // 公钥，钱包里面
}

// 判断当前的消费是谁的钱
func (txInput *TXInput) Ghw_UnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := Ghw_Ripemd160Hash(txInput.Ghw_PublicKey)

	return bytes.Compare(publicKey,ripemd160Hash) == 0
}