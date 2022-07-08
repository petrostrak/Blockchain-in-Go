package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte) *Block {
	return &Block{
		timestamp:    time.Now().Unix(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: []string{},
	}
}

func (b *Block) Print() {
	fmt.Printf("timestamp	%d\n", b.timestamp)
	fmt.Printf("nonce		%d\n", b.nonce)
	fmt.Printf("previous_hash	%x\n", b.previousHash)
	fmt.Printf("transactions	%s\n", b.transactions)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int      `json:"nonce"`
		PreviousHash [32]byte `json:"previous_hash"`
		Timestamp    int64    `json:"timestamps"`
		Transactions []string `json:"transactions"`
	}{
		Nonce:        b.nonce,
		Timestamp:    b.timestamp,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

func (b *Block) Hash() [32]byte {
	m, err := b.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}

	return sha256.Sum256([]byte(m))
}
