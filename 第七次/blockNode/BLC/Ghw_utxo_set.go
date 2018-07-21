package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"strings"
)

const utxoBucket = "chainstate"

type Ghw_UTXOSet struct {
	Ghw_Blockchain *Ghw_Blockchain
}

// 查询可花费的交易输出
func (u Ghw_UTXOSet) Ghw_FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Ghw_Blockchain.ghw_db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := Ghw_DeserializeOutputs(v)
			for outIdx, out := range outs.Ghw_Outputs {
				if out.Ghw_IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Ghw_Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

func (u Ghw_UTXOSet) Ghw_Reindex() {
	db := u.Ghw_Blockchain.ghw_db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		// 删除旧的bucket
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic()
		}
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Ghw_Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(key, outs.Ghw_Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

// 生成新区块的时候，更新UTXO数据库
func (u Ghw_UTXOSet) Update(block *Ghw_Block) {
	err := u.Ghw_Blockchain.ghw_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Ghw_Transactions {
			if !tx.Ghw_IsCoinbase() {
				for _, vin := range tx.Ghw_Vin {
					updatedOuts := Ghw_TXOutputs{}
					outsBytes := b.Get(vin.Ghw_Txid)
					outs := Ghw_DeserializeOutputs(outsBytes)

					// 找出Vin对应的outputs,过滤掉花费的
					for outIndex, out := range outs.Ghw_Outputs {
						if outIndex != vin.Ghw_Vout {
							updatedOuts.Ghw_Outputs = append(updatedOuts.Ghw_Outputs, out)
						}
					}
					// 未花费的交易输出TXOutput为0
					if len(updatedOuts.Ghw_Outputs) == 0 {
						err := b.Delete(vin.Ghw_Txid)
						if err != nil {
							log.Panic(err)
						}
					} else { // 未花费的交易输出TXOutput>0
						err := b.Put(vin.Ghw_Txid, updatedOuts.Ghw_Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}

			// 将所有的交易输出TXOutput存入数据库中
			newOutputs := Ghw_TXOutputs{}
			for _, out := range tx.Ghw_Vout {
				newOutputs.Ghw_Outputs = append(newOutputs.Ghw_Outputs, out)
			}
			err := b.Put(tx.Ghw_ID, newOutputs.Ghw_Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 打出某个公钥hash对应的所有未花费输出
func (u *Ghw_UTXOSet) Ghw_FindUTXO(pubKeyHash []byte) []Ghw_TXOutput {
	var UTXOs []Ghw_TXOutput

	err := u.Ghw_Blockchain.ghw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := Ghw_DeserializeOutputs(v)

			for _, out := range outs.Ghw_Outputs {
				if out.Ghw_IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// 查询某个地址的余额
func (u *Ghw_UTXOSet) Ghw_GetBalance(address string) int {
	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := u.Ghw_FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Ghw_Value
	}
	return balance
}

// 打印所有的UTXO
func (u *Ghw_UTXOSet) String() {
	//outputs := make(map[string][]Ghw_TXOutput)

	var lines []string
	lines = append(lines, "---ALL UTXO:")
	err := u.Ghw_Blockchain.ghw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := Ghw_DeserializeOutputs(v)

			lines = append(lines, fmt.Sprintf("     Key: %s", txID))
			for i, out := range outs.Ghw_Outputs {
				//outputs[txID] = append(outputs[txID], out)
				lines = append(lines, fmt.Sprintf("     Output: %d", i))
				lines = append(lines, fmt.Sprintf("         value:  %d", out.Ghw_Value))
				lines = append(lines, fmt.Sprintf("         PubKeyHash:  %x", out.Ghw_PubKeyHash))
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(strings.Join(lines, "\n"))
}
