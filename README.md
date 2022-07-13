#### Blockchain in Go!

In this repo I am writting a blockchain from scratch with Go in order to ultimately:

* Understand the theory and mechanisms behind Blockchain.
* Understand the verification process with Blockchain transactions
* Understand consensus algorithm that's used for deriving nonce when mining
* Understand the theory behind sending/receiving cryptocurrency
* Learn how to develop a basic Blockchain using Go
* Understand hash's role in blockchain management
* Understand how Blockchain consensus mechanisms work

#### Workflow
To run wallet and blockchain servers run
```
go run wallet/wallet-server/*.go && go run block/blockchain-server/*.go
```

Access a wallet in
```
http://localhost:8000/
```

The wallet prepopulates the public and private keys of the user as well as blockchain 
address. To send crypto to another wallet, simply fire up another wallet tab and
use the recipients address to send the desired amount.

See pending Transactions in:
```
http://localhost:5000/transactions
```

To mine pending transactions, either mine them manually by hitting:
```
http://localhost:5000/mine
```

or simulate an automatic mining procedure by hitting:
```
http://localhost:5000/mine/start
```
This will refresh the mine functionality every 20 secs.

To see the hole blockchain with all the transaction history visit:
```
http://localhost:5000/chain
```