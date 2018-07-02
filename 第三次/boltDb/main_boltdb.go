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
	//创建表
	err = db.Update(func(tx *bolt.Tx) error {
		//创建表 b表名？
		b, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return err
		}
		//向表中写入数据
		if b != nil {
			err :=b.Put([]byte("mykey"),[]byte("给张三转账了。"))
				if err !=nil {
					return fmt.Errorf("数据存储失败。。。")
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



