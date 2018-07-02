package main

import (
	"log"
	"github.com/boltdb/bolt"
)

func main()  {
	//打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//创建表
	err = db.Update(func(tx *bolt.Tx) error {
		//创建表 b表名？
			b := tx.Bucket([]byte("MyBucket"))
			//数据的更新
			if b != nil {
				err := b.Put([]byte("mykey2"),[]byte("转账 给李四 10000"))
				if err != nil {
					log.Panic("数组存储失败")
				}
			}
		//返回nil
		return nil
	})
	// 更新失败
	if err != nil{
		log.Panic(err)
	}


}



