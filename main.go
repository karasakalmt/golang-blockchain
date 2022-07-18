package main

import (
	"blockchain/blockchain"
	"fmt"
	"strconv"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("Block1")
	chain.AddBlock("Block2")
	chain.AddBlock("Block3")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x \n", block.PrevHash)
		fmt.Printf("Data in Block: %s \n", block.Data)
		fmt.Printf("Hash: %x \n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
