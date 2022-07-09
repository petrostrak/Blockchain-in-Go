package main

import (
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	myBlockChanAddress := "my_blockchain_address"
	blockChain := NewBlockchain(myBlockChanAddress)
	blockChain.Print()

	blockChain.AddTransaction("A", "B", 1.0)
	blockChain.Mining()
	blockChain.Print()

	blockChain.AddTransaction("C", "D", 2.0)
	blockChain.AddTransaction("X", "Y", 3.0)
	blockChain.Mining()
	blockChain.Print()
}
