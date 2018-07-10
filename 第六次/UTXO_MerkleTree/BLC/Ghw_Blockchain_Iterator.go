package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	Ghw_CurrentHash []byte
	Ghw_DB  *bolt.DB
}

// 遍历一条记录
func (blockchainIterator *BlockchainIterator) Ghw_Next() *Block {

	var block *Block

	err := blockchainIterator.Ghw_DB.View(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(ghw_blockTableName))

		if b != nil {
			currentBloclBytes := b.Get(blockchainIterator.Ghw_CurrentHash)
			//  获取到当前迭代器里面的currentHash所对应的区块
			block = Ghw_DeserializeBlock(currentBloclBytes)

			// 更新迭代器里面CurrentHash
			blockchainIterator.Ghw_CurrentHash = block.Ghw_PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block
}