package block

import "encoding/json"

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

func (a *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: a.Amount,
	})
}
