package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

const (
	TEMP_DIR = "wallet-server/templates"
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

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(ws.Port())), nil))
}
