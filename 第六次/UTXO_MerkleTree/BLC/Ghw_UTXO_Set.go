package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"bytes"
)

const ghw_utxoTableName  = "utxoTableName"

type UTXOSet struct {
	Ghw_Blockchain *Blockchain
}

// 重置数据库表
func (utxoSet *UTXOSet) Ghw_ResetUTXOSet()  {

	err := utxoSet.Ghw_Blockchain.Ghw_DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(ghw_utxoTableName))

		if b != nil {

			err := tx.DeleteBucket([]byte(ghw_utxoTableName))

			if err!= nil {
				log.Panic(err)
			}

		}

		b ,_ = tx.CreateBucket([]byte(ghw_utxoTableName))
		if b != nil {

			//[string]*TXOutputs
			txOutputsMap := utxoSet.Ghw_Blockchain.Ghw_FindUTXOMap()

			for keyHash,outs := range txOutputsMap {

				txHash,_ := hex.DecodeString(keyHash)

				b.Put(txHash,outs.Ghw_Serialize())

			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

// 查谒地址对应的花费记录
func (utxoSet *UTXOSet) ghw_findUTXOForAddress(address string) []*UTXO{

	var utxos []*UTXO

	utxoSet.Ghw_Blockchain.Ghw_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(ghw_utxoTableName))

		// 游标
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOutputs := Ghw_DeserializeTXOutputs(v)

			for _,utxo := range txOutputs.Ghw_UTXOS  {

				if utxo.Ghw_Output.Ghw_UnLockScriptPubKeyWithAddress(address) {
					utxos = append(utxos,utxo)
				}
			}
		}

		return nil
	})

	return utxos
}

// 查询余额
func (utxoSet *UTXOSet) Ghw_GetBalance(address string) int64 {

	UTXOS := utxoSet.ghw_findUTXOForAddress(address)

	var amount int64

	for _,utxo := range UTXOS  {
		amount += utxo.Ghw_Output.Ghw_Value
	}

	return amount
}

// 返回要凑多少钱，对应TXOutput的TX的Hash和index
func (utxoSet *UTXOSet) Ghw_FindUnPackageSpendableUTXOS(from string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	for _,tx := range txs {

		if tx.Ghw_IsCoinbaseTransaction() == false {

			for _, in := range tx.Ghw_Vins {
				//是否能够解锁
				publicKeyHash := Ghw_Base58Decode([]byte(from))

				ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]

				if in.Ghw_UnLockRipemd160Hash(ripemd160Hash) {

					key := hex.EncodeToString(in.Ghw_TxHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Ghw_Vout)
				}

			}
		}
	}

	for _,tx := range txs {

	Work1:
		for index,out := range tx.Ghw_Vouts {

			if out.Ghw_UnLockScriptPubKeyWithAddress(from) {

				fmt.Println(from)

				fmt.Println(spentTXOutputs)

				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.Ghw_TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.Ghw_TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.Ghw_TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.Ghw_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}

	return unUTXOs

}

// 查询花费支出
func (utxoSet *UTXOSet) Ghw_FindSpendableUTXOS(from string,amount int64,txs []*Transaction) (int64,map[string][]int)  {

	unPackageUTXOS := utxoSet.Ghw_FindUnPackageSpendableUTXOS(from,txs)

	spentableUTXO := make(map[string][]int)

	var money int64 = 0

	for _,UTXO := range unPackageUTXOS {

		money += UTXO.Ghw_Output.Ghw_Value;
		txHash := hex.EncodeToString(UTXO.Ghw_TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash],UTXO.Ghw_Index)
		if money >= amount{
			return  money,spentableUTXO
		}
	}

	// 钱还不够
	utxoSet.Ghw_Blockchain.Ghw_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(ghw_utxoTableName))

		if b != nil {

			c := b.Cursor()
			UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := Ghw_DeserializeTXOutputs(v)

				for _,utxo := range txOutputs.Ghw_UTXOS {

					money += utxo.Ghw_Output.Ghw_Value
					txHash := hex.EncodeToString(utxo.Ghw_TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash],utxo.Ghw_Index)

					if money >= amount {
						 break UTXOBREAK;
					}
				}
			}

		}

		return nil
	})

	if money < amount{
		log.Panic("余额不足......")
	}

	return  money,spentableUTXO
}


// 更新
func (utxoSet *UTXOSet) Ghw_Update()  {

	// 最新的Block
	block := utxoSet.Ghw_Blockchain.Ghw_Iterator().Ghw_Next()

	ins := []*TXInput{}

	outsMap := make(map[string]*TXOutputs)

	// 找到所有我要删除的数据
	for _,tx := range block.Ghw_Txs {

		for _,in := range tx.Ghw_Vins {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.Ghw_Txs  {

		utxos := []*UTXO{}

		for index,out := range tx.Ghw_Vouts  {

			isSpent := false

			for _,in := range ins  {

				if in.Ghw_Vout == index && bytes.Compare(tx.Ghw_TxHash ,in.Ghw_TxHash) == 0 && bytes.Compare(out.Ghw_Ripemd160Hash,Ghw_Ripemd160Hash(in.Ghw_PublicKey)) == 0 {

					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &UTXO{tx.Ghw_TxHash,index,out}
				utxos = append(utxos,utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.Ghw_TxHash)
			outsMap[txHash] = &TXOutputs{utxos}
		}

	}

	err := utxoSet.Ghw_Blockchain.Ghw_DB.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(ghw_utxoTableName))

		if b != nil {
			// 删除
			for _,in := range ins {

				txOutputsBytes := b.Get(in.Ghw_TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				fmt.Println(txOutputsBytes)

				txOutputs := Ghw_DeserializeTXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _,utxo := range txOutputs.Ghw_UTXOS  {

					if in.Ghw_Vout == utxo.Ghw_Index && bytes.Compare(utxo.Ghw_Output.Ghw_Ripemd160Hash,Ghw_Ripemd160Hash(in.Ghw_PublicKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}

				if isNeedDelete {
					b.Delete(in.Ghw_TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.Ghw_TxHash)]

						preTXOutputs.Ghw_UTXOS = append(preTXOutputs.Ghw_UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.Ghw_TxHash)] = preTXOutputs

					}
				}

			}

			// 新增
			for keyHash,outPuts := range outsMap  {
				keyHashBytes,_ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes,outPuts.Ghw_Serialize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}




