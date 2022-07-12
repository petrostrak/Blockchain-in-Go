package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

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

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(ws.Port())), nil))
}
