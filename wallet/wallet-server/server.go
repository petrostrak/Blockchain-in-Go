package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/petrostrak/Blockchain-in-Go/utils"
	"github.com/petrostrak/Blockchain-in-Go/wallet"
)

const (
	TEMP_DIR = "templates"
)

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port: port, gateway: gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(TEMP_DIR, "index.html"))
		if err != nil {
			log.Println(err)
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Println(err)
		}
	default:
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, err := myWallet.MarshalJSON()
		if err != nil {
			log.Println(err)
		}

		io.WriteString(w, string(m))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t wallet.TreansactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("[ERROR]: could not decode: %v\n", err)
			io.WriteString(w, string(utils.JSONStatus("fail")))
			return
		}

		if !t.Validate() {
			log.Println("[ERROR]: missing fields")
			io.WriteString(w, string(utils.JSONStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Println("[ERROR]: parse failed")
			io.WriteString(w, string(utils.JSONStatus("fail")))
			return
		}

		value32 := float32(value)

		w.Header().Add("Content-Type", "application/json")
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(ws.Port())), nil))
}
