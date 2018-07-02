package main

import (
	"flag"
	"fmt"
)

// addBlock -data "data",增加区块数据
// printchain ,打印所有的区块数据
//flag -string "早上好" -number 8 -bool
//flag -string "早上好" -number 8
func main()  {
	flagString := flag.String("string","","输出所有的区块信息")
	flagInt := flag.Int("number",6,"输出一个整数")
	flagBool := flag.Bool("bool",false,"输出一个真或者假")
	flag.Parse()
	fmt.Printf("%s\n",*flagString)
	fmt.Printf("%s\n",*flagInt)
	fmt.Printf("%s\n",*flagBool)
}
