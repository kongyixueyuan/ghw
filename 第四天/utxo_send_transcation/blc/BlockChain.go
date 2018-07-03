package blc

import (
	"block/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
	"encoding/hex"
	"strconv"
)

const dbname  = "my_blockChain123.db"
const tableName  = "my_blocks"
const lashHashKey  = "lastHashKey" //最新的hash key
type BlockChain struct {
	//Blocks []*Block //存储区块链
	Tip []byte //最新区块的hash
	DB *bolt.DB
}



//打印区块信息
func (blc *BlockChain) PrintChains()  {
	var block *Block//区块对象指针
	var currentHash []byte = blc.Tip
	for   {
		err := blc.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(tableName))
			if b != nil {
				blockBytes := b.Get(currentHash)//当前区块
				block = DeserializeBlock(blockBytes)
				fmt.Printf("区块的高度%d\n",block.Height)
				fmt.Printf("区块的高度%x\n",block.Hash)
				fmt.Printf("区块的时间戳%d\n",time.Unix(block.Timestamp,0).Format("2006-01-02 03:04:25 PM"))
				fmt.Printf("区块的随机数%d\n",block.Nonce)
				fmt.Printf("上一个区块的高度%x\n",block.PrevBlockHash)


				currentHash = block.PrevBlockHash
			}
			return nil
		})
		if err != nil{
			log.Panic(err)
		}
		var hashInt big.Int
		hashInt.SetBytes([]byte(block.PrevBlockHash))
		if big.NewInt(0).Cmp(&hashInt) == 0{
			break
		}
	}

}

//增加区块到区块链中
func (blc *BlockChain) AddBlockToBlockChain(txs []*Transaction)  {
	err :=blc.DB.Update(func(tx *bolt.Tx) error {
		//1,获取表
		table := tx.Bucket([]byte(tableName))
		if table !=nil {
			//获取最新的区块
			blockBytes :=table.Get(blc.Tip)
			block := DeserializeBlock(blockBytes)

			///2,创建新的区块
			//增加区块
			newBlock := NewBlock(txs,block.Height+1,block.Hash)
			//3,新区快写入到数据库中,更新数据库中的lasthash对应的hash值
			err := table.Put(newBlock.Hash,newBlock.Serialize())
			if err != nil{
				log.Panic("入库失败")
			}
			err = table.Put([]byte(lashHashKey),newBlock.Hash)
			if err != nil{
				log.Panic("入库失败")
			}
			//4,更新blackChian tip
			blc.Tip = newBlock.Hash

		}
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
/*	//新区块加入链中
	blc.Blocks = append(blc.Blocks,newBlock)*/
}


//创建带有创世区块的区块链
func  CreateBlockChainWithGenesisBlock(address string) *BlockChain   {
	// 判断数据库是否存在
	if IsDBExists(dbname) {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}
	fmt.Println("正在创建创世区块.......")
	// 创建或者打开数据库
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var genesisHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		// 创建数据库表
		b, err := tx.CreateBucket([]byte(tableName))
		if err != nil {
			log.Panic(err)
		}
		if b != nil {
			// 创建创世区块
			txCoinBase := NewConinbaseTransaction(address)
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinBase})
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 存储最新的区块的hash
			err = b.Put([]byte(lashHashKey), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			genesisHash = genesisBlock.Hash
		}
		return nil
	})
	return &BlockChain{genesisHash,db}
}

//判断数据库是否存在
func IsDBExists(dbName string) bool {
	//if _, err := os.Stat(dbName); os.IsNotExist(err) {
	//
	//	return false
	//}
	_, err := os.Stat(dbName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// 返回Blockchain对象
func BlockchainObject() *BlockChain {
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			// 读取最新区块的Hash
			tip = b.Get([]byte(lashHashKey))

		}
		return nil
	})
	return &BlockChain{tip,db}
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *BlockChain) UnUTXOs(address string) []*TXOutput {
	var unUTXOs []*TXOutput
	spentTXOutputs := make(map[string][]int)
	//{hash:[0]}
	blockIterator := blockchain.Iterator()
	for  {
		block := blockIterator.Next()
		for _,tx := range block.Txs {
			// 花费
			if tx.IsCoinbaseTransaction() == false {
				for _,in := range tx.Vins {
					//是否能够解锁
					if in.UnLockWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key],in.Vout)
					}
				}
			}


			// Vouts
			for index,out := range tx.Vouts {
				if out.UnLockScriptPubKeyWithAddress(address) {
					fmt.Println(out)
					fmt.Println(spentTXOutputs)
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							for txHash,indexArray := range spentTXOutputs {
								for _,i := range  indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash){
										continue
									} else {
										unUTXOs = append(unUTXOs,out)
									}
								}
							}
						} else {
							unUTXOs = append(unUTXOs,out)
						}

					}
				}

			}

		}
		fmt.Println(spentTXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}
	return unUTXOs
}




// 查询余额
func (blockchain *BlockChain) GetBalance(address string) int64 {
	utxos := blockchain.UnUTXOs(address)
	var amount int64
	for _,out := range utxos {
		amount = amount + out.Money
	}
	return amount
}

func (blockChain *BlockChain) MineNewBlock(from []string,to []string,amount []string)   {
/*	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)*/
	//1,建立transaction数组
/*	var txs  []*Transaction
	var block *Block
	blockChain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b!=nil{
			hash := b.Get([]byte(lashHashKey))
			blockBytes :=b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	//2,建立新的区块
	block =NewBlock(txs,block.Height+1,block.Hash)
	// 新区块入库
	blockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		if b!=nil{
			b.Put(block.Hash,block.Serialize())
			b.Put([]byte(lashHashKey),block.Hash)
			blockChain.Tip = block.Hash
		}
		return nil
	})*/

	value ,_ := strconv.Atoi(amount[0])
	tx := NewSimpleTransaction(from[0],to[0],value)
	//1. 通过相关算法建立Transaction数组
	var txs []*Transaction
	txs = append(txs,tx)
	var block *Block
	blockChain.DB.View(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			hash := b.Get([]byte(lashHashKey))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	//2. 建立新的区块
	block = NewBlock(txs,block.Height+1,block.Hash)
	//将新区块存储到数据库
	blockChain.DB.Update(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(tableName))
		if b != nil {
			b.Put(block.Hash,block.Serialize())
			b.Put([]byte(lashHashKey),block.Hash)
			blockChain.Tip = block.Hash
		}
		return nil
	})


}


