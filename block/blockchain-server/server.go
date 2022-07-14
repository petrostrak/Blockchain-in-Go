package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/petrostrak/Blockchain-in-Go/block"
	"github.com/petrostrak/Blockchain-in-Go/utils"
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
		log.Printf("private_key %v\n", minersWallet.PrivateKeyStr())
		log.Printf("public_key %v\n", minersWallet.PublicKeyStr())
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

func (bs *BlockchainServer) Transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bs.GetBlockchain()
		transaction := bc.TransactionPool()
		m, err := json.Marshal(struct {
			Transactions []*block.Transaction `json:"transactions"`
			Length       int                  `json:"length"`
		}{
			Transactions: transaction,
			Length:       len(transaction),
		})
		if err != nil {
			log.Printf("[ERROR]: %v\n", err)
			return
		}
		io.WriteString(w, string(m))
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t *block.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("[ERROR]: %v\n", err)
			return
		}
		if !t.Validate() {
			log.Println("[ERROR]: missing field(s)")
			io.WriteString(w, string(utils.JSONStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)
		bc := bs.GetBlockchain()
		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JSONStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = utils.JSONStatus("success")
		}

		io.WriteString(w, string(m))
	case http.MethodPut:
		decoder := json.NewDecoder(r.Body)
		var t *block.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("[ERROR]: %v\n", err)
			return
		}
		if !t.Validate() {
			log.Println("[ERROR]: missing field(s)")
			io.WriteString(w, string(utils.JSONStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)
		bc := bs.GetBlockchain()
		isUpdated := bc.AddTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isUpdated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JSONStatus("fail")
		} else {
			m = utils.JSONStatus("success")
		}

		io.WriteString(w, string(m))
	case http.MethodDelete:
		bc := bs.GetBlockchain()
		bc.ClearTransactionPool()
		io.WriteString(w, string(utils.JSONStatus("success")))
	default:
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bs *BlockchainServer) Mine(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bc := bs.GetBlockchain()
		isMined := bc.Mining()

		var m []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JSONStatus("fail")
		} else {
			m = utils.JSONStatus("success")
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	default:
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bs *BlockchainServer) StartMine(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bc := bs.GetBlockchain()
		bc.StartMining()

		m := utils.JSONStatus("success")
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	default:
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bs *BlockchainServer) Amount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		blockchainAddress := r.URL.Query().Get("blockchain_address")
		amount := bs.GetBlockchain().CalculateTotalAmout(blockchainAddress)

		a := &block.AmountResponse{Amount: amount}
		m, err := a.MarshalJSON()
		if err != nil {
			log.Printf("[ERROR]: unable to marshal json: %v\n", err)
		}

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	default:
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bs *BlockchainServer) Run() {
	bs.GetBlockchain().Run()

	http.HandleFunc("/", bs.GetChain)
	http.HandleFunc("/transactions", bs.Transactions)
	http.HandleFunc("/mine", bs.Mine)
	http.HandleFunc("/mine/start", bs.StartMine)
	http.HandleFunc("/amount", bs.Amount)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(bs.Port())), nil))
}
