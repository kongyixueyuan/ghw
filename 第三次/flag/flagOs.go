package main

import (
	"os"
	"fmt"
)

func main()  {

	args := os.Args
	//fmt.Println(args)
	fmt.Println(args[0])
	fmt.Println(args[1])
	fmt.Println(args[2])
}
