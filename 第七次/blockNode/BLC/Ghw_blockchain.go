package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"encoding/hex"
	"strconv"
	"crypto/ecdsa"
	"bytes"
	"github.com/pkg/errors"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "genesis data"

type Ghw_Blockchain struct {
	ghw_tip []byte
	ghw_db  *bolt.DB
}

// 打印区块链内容
func (bc *Ghw_Blockchain) Ghw_Printchain() {
	bci := bc.Ghw_Iterator()

	for {
		block := bci.Ghw_Next()
		block.String()
		if len(block.Ghw_PrevBlockHash) == 0 {
			break
		}
	}

}

// 通过交易hash,查找交易
func (bc *Ghw_Blockchain) Ghw_FindTransaction(ID []byte) (Ghw_Transaction, error) {
	bci := bc.Ghw_Iterator()
	for {
		block := bci.Ghw_Next()
		for _, tx := range block.Ghw_Transactions {
			if bytes.Compare(tx.Ghw_ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.Ghw_PrevBlockHash) == 0 {
			break
		}
	}
	fmt.Printf("查找%x的交易失败",ID)
	return Ghw_Transaction{}, errors.New("未找到交易")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Ghw_Blockchain) FindUTXO() map[string]Ghw_TXOutputs {
	// 未花费的交易输出
	// key:交易hash   txID
	UTXO := make(map[string]Ghw_TXOutputs)
	// 已经花费的交易txID : TXOutputs.index
	spentTXOs := make(map[string][]int)
	bci := bc.Ghw_Iterator()

	for {
		block := bci.Ghw_Next()

		// 循环区块中的交易
		for _, tx := range block.Ghw_Transactions {
			// 将区块中的交易hash，转为字符串
			txID := hex.EncodeToString(tx.Ghw_ID)

		Outputs:
			for outIdx, out := range tx.Ghw_Vout { // 循环交易中的 TXOutputs
				// Was the output spent?
				// 如果已经花费的交易输出中，有此输出，证明已经花费
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx { // 如果花费的正好是此笔输出
							continue Outputs // 继续下一次循环
						}
					}
				}

				outs := UTXO[txID] // 获取UTXO指定txID对应的TXOutputs
				outs.Ghw_Outputs = append(outs.Ghw_Outputs, out)
				UTXO[txID] = outs
			}

			if tx.Ghw_IsCoinbase() == false { // 非创世区块
				for _, in := range tx.Ghw_Vin {
					inTxID := hex.EncodeToString(in.Ghw_Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Ghw_Vout)
				}
			}
		}
		// 如果上一区块的hash为0，代表已经到创世区块，循环结束
		if len(block.Ghw_PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// 获取迭代器
func (bc *Ghw_Blockchain) Ghw_Iterator() *Ghw_BlockchainIterator {
	return &Ghw_BlockchainIterator{bc.ghw_tip, bc.ghw_db}
}

// 新建区块链(包含创世区块)
func Ghw_CreateBlockchain(address string,nodeID string) *Ghw_Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if ghw_dbExists(dbFile) {
		fmt.Println("blockchain数据库已经存在.")
		os.Exit(1)
	}

	var tip []byte
	cbtx := Ghw_NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := Ghw_NewGenesisBlock(cbtx)

	//genesis.String()

	// 打开数据库，如果不存在自动创建
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		// 新区块存入数据库
		err = b.Put(genesis.Ghw_Hash, genesis.Ghw_Serialize())
		if err != nil {
			log.Panic(err)
		}
		// 将创世区块的hash存入数据库
		err = b.Put([]byte("l"), genesis.Ghw_Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Ghw_Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Ghw_Blockchain{tip, db}
}

// 获取blockchain对象
func Ghw_NewBlockchain(nodeID string) *Ghw_Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if !ghw_dbExists(dbFile) {
		log.Panic("区块链还未创建")
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Ghw_Blockchain{tip, db}
}

// 生成新的区块(挖矿)
func (bc *Ghw_Blockchain) MineNewBlock(from []string, to []string, amount []string,nodeID string , mineNow bool) *Ghw_Block {
	UTXOSet := Ghw_UTXOSet{bc}

	wallets, err := Ghw_NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	var txs []*Ghw_Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		if value<=0 {
			log.Panic("错误：转账金额需要大于0")
		}
		wallet := wallets.Ghw_GetWallet(address)
		tx := Ghw_NewUTXOTransaction(&wallet, to[index], value, &UTXOSet, txs)
		txs = append(txs, tx)
	}

	if mineNow {
		// 挖矿奖励
		tx := Ghw_NewCoinbaseTX(from[0], "")
		txs = append(txs, tx)

		//=====================================
		newBlock := bc.Ghw_MineBlock(txs)
		UTXOSet.Update(newBlock)
		return newBlock
	}else{
		// 如果不立即挖矿，将交易写到内存中
		//var txs_all []Ghw_Transaction
		//for _,value := range txs{
		//	txs_all= append(txs_all, *value)
		//}
		ghw_sendTxs(knownNodes[0],txs)
		return nil
	}


}

// 挖矿
func (bc *Ghw_Blockchain) Ghw_MineBlock(txs []*Ghw_Transaction) *Ghw_Block  {
	var lashHash []byte
	var lastHeight int

	// 检查交易是否有效，验证签名
	for _, tx := range txs {
		if !bc.Ghw_VerifyTransaction(tx, txs) {
			log.Panic("错误：无效的交易")
		}
	}
	// 获取最后一个区块的hash,然后获取最后一个区块的信息，进而获得height
	err := bc.ghw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lashHash = b.Get([]byte("l"))
		blockData := b.Get(lashHash)
		block := Ghw_DeserializeBlock(blockData)
		lastHeight = block.Ghw_Height
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	// 生成新的区块
	newBlock := Ghw_NewBlock(txs, lashHash, lastHeight+1)

	// 将新区块的内容更新到数据库中
	err = bc.ghw_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Ghw_Hash, newBlock.Ghw_Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.Ghw_Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.ghw_tip = newBlock.Ghw_Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return newBlock
}

// 签名
func (bc *Ghw_Blockchain) Ghw_SignTransaction(tx *Ghw_Transaction, privKey ecdsa.PrivateKey,txs []*Ghw_Transaction) {
	prevTXs := make(map[string]Ghw_Transaction)

	// 找到交易输入中，之前的交易
	Vin:
	for _, vin := range tx.Ghw_Vin {
		for _, tx := range txs {
			if bytes.Compare(tx.Ghw_ID, vin.Ghw_Txid) == 0 {
				prevTX := *tx
				prevTXs[hex.EncodeToString(prevTX.Ghw_ID)] = prevTX
				continue Vin
			}
		}

		prevTX, err := bc.Ghw_FindTransaction(vin.Ghw_Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Ghw_ID)] = prevTX

	}

	tx.Ghw_Sign(privKey, prevTXs)
}

// 验证签名
func (bc *Ghw_Blockchain) Ghw_VerifyTransaction(tx *Ghw_Transaction,txs []*Ghw_Transaction) bool {
	if tx.Ghw_IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Ghw_Transaction)
	Vin:
	for _, vin := range tx.Ghw_Vin {
		for _, tx := range txs {
			if bytes.Compare(tx.Ghw_ID, vin.Ghw_Txid) == 0 {
				prevTX := *tx
				prevTXs[hex.EncodeToString(prevTX.Ghw_ID)] = prevTX
				continue Vin
			}
		}
		prevTX, err := bc.Ghw_FindTransaction(vin.Ghw_Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Ghw_ID)] = prevTX
	}

	return tx.Ghw_Verify(prevTXs)
}

// 判断数据库是否存在
func ghw_dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// 获取BestHeight
func (bc *Ghw_Blockchain) Ghw_GetBestHeight() int {
	var lastBlock Ghw_Block

	err := bc.ghw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *Ghw_DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Ghw_Height
}

// 获取所有区块的hash
func (bc *Ghw_Blockchain) Ghw_GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Ghw_Iterator()

	for {
		block := bci.Ghw_Next()

		blocks = append(blocks, block.Ghw_Hash)

		if len(block.Ghw_PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

// 根据hash获取某个区块的内容
func (bc *Ghw_Blockchain) Ghw_GetBlock(blockHash []byte) (Ghw_Block, error) {
	var block Ghw_Block

	err := bc.ghw_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("未找到区块")
		}

		block = *Ghw_DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// 将区块添加到链中
func (bc *Ghw_Blockchain) Ghw_AddBlock(block *Ghw_Block) {
	err := bc.ghw_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Ghw_Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Ghw_Serialize()
		err := b.Put(block.Ghw_Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := Ghw_DeserializeBlock(lastBlockData)

		if block.Ghw_Height > lastBlock.Ghw_Height {
			err = b.Put([]byte("l"), block.Ghw_Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.ghw_tip = block.Ghw_Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}