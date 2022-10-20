package types

import (
	"math"
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Block     uint      `json:"block"`
	Timestamp time.Time `json:"timestamp"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Value     float64   `json:"value"`
	NumVer    int       `json:"numver"`
	Fee       float64   `json:"fee"`
}

type TransactionRaw struct {
	ID       string `json:"hash"`
	Block    string `json:"blockNumber"`
	Sender   string `json:"from"`
	Receiver string `json:"to"`
	Value    string `json:"value"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
}

func (t TransactionRaw) RawToNormal(tstamp string) (*Transaction, error) {
	value := ConvHexWeiToDecEth(t.Value)

	gas, err := ConvHexDec(t.Gas)
	if err != nil {
		return nil, err
	}
	gasprice, err := ConvHexDec(t.GasPrice)
	if err != nil {
		return nil, err
	}
	fee := float64(gas) * float64(gasprice) / math.Pow(10, 18)

	timestamp, err := ConvHexDec(tstamp)
	if err != nil {
		return nil, err
	}

	block, err := ConvHexDec(t.Block)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:        t.ID,
		Block:     uint(block),
		Timestamp: time.Unix(timestamp, 0),
		Sender:    t.Sender,
		Receiver:  t.Receiver,
		Value:     value,
		Fee:       fee,
	}, nil
}
