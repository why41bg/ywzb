package main

import (
	"github.com/boltdb/bolt"
	"ywzb/block"
	"ywzb/cli"
)

func main() {
	bc, err := block.InitBlockChain()
	if err != nil {
		panic(err)
	}
	defer func(Db *bolt.DB) {
		err := Db.Close()
		if err != nil {
			panic(err)
		}
	}(bc.Db)

	c := cli.CLI{BChain: bc}
	c.Run()
}
