package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 5000, "TCP port number of the blockchain server")
	flag.Parse()

	app := NewBlockchainServer(uint16(*port))
	app.Run()
}
