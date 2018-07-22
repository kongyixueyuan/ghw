package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type Ghw_BlockchainIterator struct {
	ghw_currentHash []byte
	ghw_db          *bolt.DB
}

func (i *Ghw_BlockchainIterator) Ghw_Next() *Ghw_Block {
	var block *Ghw_Block

	err := i.ghw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.ghw_currentHash)
		block = Ghw_DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.ghw_currentHash = block.Ghw_PrevBlockHash

	return block
}