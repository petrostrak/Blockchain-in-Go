package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/petrostrak/Blockchain-in-Go/block"
	"github.com/petrostrak/Blockchain-in-Go/wallet"
)

var (
	cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)
)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port: port}
}

func (bs *BlockchainServer) Port() uint16 {
	return bs.port
}

func (bs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, exists := cache["blockchain"]
	if !exists {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minersWallet.BlockchainAddress(), bs.Port())
		cache["blockchain"] = bc
		log.Printf("private_key %v\n", minersWallet.PrivateKey())
		log.Printf("public_key %v\n", minersWallet.PublicKey())
		log.Printf("blockchain_address %v\n", minersWallet.BlockchainAddress())
	}

	return bc
}

func (bs *BlockchainServer) GetChain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bs.GetBlockchain()
		m, err := bc.MarshalJSON()
		if err != nil {
			log.Println(err)
		}

		io.WriteString(w, string(m[:]))
	default:
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
	}
}

func (bs *BlockchainServer) Run() {
	http.HandleFunc("/", bs.GetChain)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(bs.Port())), nil))
}
