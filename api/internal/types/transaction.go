package types

import (
	"math"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	ID        string
	Block     uint
	Timestamp time.Time
	Sender    string
	Reciever  string
	Value     float64
	NumVer    int
}

type TransactionRaw struct {
	ID       string `json:"hash"`
	Block    string `json:"blockNumber"`
	Sender   string `json:"from"`
	Reciever string `json:"to"`
	Value    string `json:"value"`
	Tstamp   string `json:"timestamp"`
}

func (t TransactionRaw) RawToNormal() (*Transaction, error) {
	t.Block, t.Value, t.Tstamp = strings.Replace(t.Block, "0x", "", -1), strings.Replace(t.Value, "0x", "", -1), strings.Replace(t.Tstamp, "0x", "", -1)

	block, err := strconv.ParseInt(t.Block, 16, 64)
	if err != nil {
		return nil, err
	}

	valint, err := strconv.ParseInt(t.Value, 16, 64)
	if err != nil {
		return nil, err
	}
	value := float64(valint) / math.Pow(10, 17)

	timestamp, err := strconv.ParseInt(t.Tstamp, 16, 64)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:        t.ID,
		Block:     uint(block),
		Timestamp: time.Unix(timestamp, 0),
		Sender:    t.Sender,
		Reciever:  t.Reciever,
		Value:     value,
	}, nil
}
