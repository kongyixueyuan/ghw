package main

import (
	"block/utxo_send_json/blc"
	"fmt"
)

func main()  {
	var str1 string
	str1 = "[\"ghw\",\"zhangsan\"]"
	str2:=blc.JSONToArray(str1)
	fmt.Println(str2)
}
