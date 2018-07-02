package main

import (
	"log"
	"github.com/boltdb/bolt"
	"fmt"
)

func main()  {
	//打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//获取表
	err = db.View(func(tx *bolt.Tx) error {
		//创建表 b表名？
			b := tx.Bucket([]byte("MyBucket"))
			//数据的更新
			if b != nil {
				data := b.Get([]byte("mykey"))
				fmt.Printf("数据库中的数据%s\n",data)
				data = b.Get([]byte("mykey2"))
				fmt.Printf("数据库中的数据%s\n",data)
			}
		//返回nil
		return nil
	})
	// 更新失败
	if err != nil{
		log.Panic(err)
	}


}



