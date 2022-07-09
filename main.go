package main

import (
	"fmt"
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

	fmt.Printf("my_blockchain_address %.1f\n", blockChain.CalculateTotalAmout("my_blockchain_address"))
	fmt.Printf("C %.1f\n", blockChain.CalculateTotalAmout("C"))
	fmt.Printf("D %.1f\n", blockChain.CalculateTotalAmout("D"))
}
