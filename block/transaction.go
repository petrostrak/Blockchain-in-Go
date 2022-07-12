package block

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	sendersBlockchainAddress   string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(sender, recipient string, value float32) *Transaction {
	return &Transaction{
		sender,
		recipient,
		value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address	%s\n", t.sendersBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address	%s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value				%.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.sendersBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blochchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blochchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	return tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil ||
		tr.Signature == nil
}
