package main

import (
	"fmt"
	"log"
	"time"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []string
}

func (b *Block) Print() {
	fmt.Printf("timestamp	%d\n", b.timestamp)
	fmt.Printf("nonce		%d\n", b.nonce)
	fmt.Printf("previous_hash	%s\n", b.previousHash)
	fmt.Printf("transactions	%s\n", b.transactions)
}

func NewBlock(nonce int, previousHash string) *Block {
	return &Block{
		timestamp:    time.Now().Unix(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: []string{},
	}
}

func main() {
	b := NewBlock(0, "init hash")
	b.Print()
}
