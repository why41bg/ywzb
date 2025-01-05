package cli

import (
	"flag"
	"fmt"
	"os"
	"ywzb/block"
)

type CLI struct {
	BChain *block.Chain
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// Run is the entry point for the CLI
func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		_ = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		_ = printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	err := cli.BChain.AddBlock(data)
	if err != nil {
		fmt.Println("Failed to add block")
		return
	}
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.BChain.Iterator()

	for {
		b := bci.Next()

		fmt.Printf("Timestamp: %d\n", b.Timestamp)
		fmt.Printf("Data: %s\n", b.Data)
		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Println()

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA")
	fmt.Println("  printchain")
}
