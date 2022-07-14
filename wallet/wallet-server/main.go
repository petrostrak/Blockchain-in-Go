package main

import (
	"flag"
	"log"

	"github.com/petrostrak/Blockchain-in-Go/utils"
)

func init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	port := flag.Uint("port", 8000, "TCP port number of the wallet server")
	gateway := flag.String("gateway", utils.GetHost()+":5000", "Blockchain Gateway")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *gateway)
	app.Run()
}
