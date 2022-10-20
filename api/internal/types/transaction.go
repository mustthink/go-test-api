package types

import (
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Block     uint      `json:"block"`
	Timestamp time.Time `json:"timestamp"`
	Sender    string    `json:"sender"`
	Reciever  string    `json:"reciever"`
	Value     float64   `json:"value"`
	NumVer    int       `json:"numver"`
	Fee       float64   `json:"fee"`
}

type TransactionRaw struct {
	ID       string `json:"hash"`
	Block    string `json:"blockNumber"`
	Sender   string `json:"from"`
	Reciever string `json:"to"`
	Value    string `json:"value"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
}

func (t TransactionRaw) RawToNormal(tstamp string) (*Transaction, error) {
	block, err := ConvHexDec(t.Block)
	if err != nil {
		return nil, err
	}

	value := ConvHexWeiToDecEth(t.Value)

	timestamp, err := ConvHexDec(tstamp)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:        t.ID,
		Block:     uint(block),
		Timestamp: time.Unix(int64(timestamp), 0),
		Sender:    t.Sender,
		Reciever:  t.Reciever,
		Value:     value,
	}, nil
}
