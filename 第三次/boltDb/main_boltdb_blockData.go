package main


import (
	"fmt"
	"./blc"
	"log"
	"github.com/boltdb/bolt"
)

func main()  {
	//创建区块数据
	block:=blc.NewBlock("创世区块",1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
	fmt.Println("原始区块")
	fmt.Printf("nonce值%d\n",block.Nonce)
	fmt.Printf("hash值 %x\n",block.Hash)

	//打开数据库
	db, err := bolt.Open("boltdb_blocks.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//创建表
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("boltdb_blocks"))
		if b != nil {
			//区块数据序列化
			blockStr := block.Serialize()
			err = b.Put([]byte("firstBlocks"), []byte(blockStr))
			if err != nil {
				return fmt.Errorf("数据存储失败")
			}
		}
		if err != nil {
			return err
		}
		return nil
	})
	// 更新失败
	if err != nil{
		log.Panic(err)
	}

	//数据查询
	fmt.Println("从数据库中读取区块")
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("boltdb_blocks"))
		if b !=nil {
			blockData := b.Get([]byte("firstBlocks"))
			bd := blc.DeserializeBlock(blockData)
			fmt.Printf("The answer is: %s\n", bd)
			fmt.Printf("nonce值%d\n",bd.Nonce)
			fmt.Printf("hash值 %x\n",bd.Hash)

		}
		return nil
	})
}
