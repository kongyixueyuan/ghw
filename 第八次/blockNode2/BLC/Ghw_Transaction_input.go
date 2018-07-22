package BLC

import "bytes"

type Ghw_TXInput struct {
	Ghw_Txid      []byte
	Ghw_Vout      int      // Vout的index
	Ghw_Signature []byte   // 签名
	Ghw_PubKey    []byte   // 公钥
}

func (in Ghw_TXInput) UsesKey(pubKeyHash []byte) bool  {
	lockingHash := Ghw_HashPubKey(in.Ghw_PubKey)

	return bytes.Compare(lockingHash,pubKeyHash) == 0
}
