package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/petrostrak/Blockchain-in-Go/block"
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
		var t wallet.TransactionRequest
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

		transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &block.TransactionRequest{
			t.SenderBlockchainAddress,
			t.RecipientBlockchainAddress,
			t.SenderPublicKey,
			&value32,
			&signatureStr,
		}

		m, err := json.Marshal(bt)
		if err != nil {
			log.Println("[ERROR]: failed to marshal transaction request")
			io.WriteString(w, string(utils.JSONStatus("fail")))
			return
		}

		buf := bytes.NewBuffer(m)

		resp, err := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if err != nil {
			log.Println(err)
		}
		if resp.StatusCode == 201 {
			io.WriteString(w, string(utils.JSONStatus("success")))
			return
		}
		log.Println("[ERROR]: could not post to transactions")
		io.WriteString(w, string(utils.JSONStatus("fail")))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
	}
}

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		blockchainAddress := r.URL.Query().Get("blockchain_address")
		entpoint := fmt.Sprintf("%s/amount", ws.Gateway())

		client := &http.Client{}
		bsReq, _ := http.NewRequest("GET", entpoint, nil)
		q := bsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		bsReq.URL.RawQuery = q.Encode()

		bsResp, err := client.Do(bsReq)
		if err != nil {
			log.Printf("[ERROR]: unable to send request: %v\n", err)
			io.WriteString(w, string(utils.JSONStatus("fail")))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if bsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bsResp.Body)
			var bar block.AmountResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("[ERROR]: unable to decode to amount respose: %v\n", err)
				io.WriteString(w, string(utils.JSONStatus("fail")))
				return
			}

			m, err := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			if err != nil {
				log.Printf("[ERROR]: unable to marshal amount: %v\n", err)
				io.WriteString(w, string(utils.JSONStatus("fail")))
				return
			}
			io.WriteString(w, string(m))
		} else {
			io.WriteString(w, string(utils.JSONStatus("fail")))
		}

	default:
		log.Printf("[ERROR]: Invalid request method: %v\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/wallet/amount", ws.WalletAmount)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(ws.Port())), nil))
}
