package handlers

import (
	"encoding/json"
	"github.com/mustthink/go-test-api/internal/types"
	"io"
	"log"
	"net/http"
	"time"
)

func (app *application) BodyToData(body []byte) (int64, string, []types.TransactionRaw) {
	jsondata := struct {
		Result struct {
			Number       string                 `json:"number"`
			Timestamp    string                 `json:"timestamp"`
			Transactions []types.TransactionRaw `json:"transactions"`
		} `json:"result"`
	}{}

	err := json.Unmarshal(body, &jsondata)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	newid, err := types.ConvHexDec(jsondata.Result.Number)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	return newid, jsondata.Result.Timestamp, jsondata.Result.Transactions
}

func (app *application) GetLastID() int64 {
	req, err := http.NewRequest("GET", app.service.GetReqID(), nil)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := app.client.Do(req)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	jsondata := struct {
		Number string `json:"result"`
	}{}

	err = json.Unmarshal(data, &jsondata)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	lastid, err := types.ConvHexDec(jsondata.Number)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	return lastid
}

func (app *application) ScanBlocks(t int) {
	var lastid uint
	for {
		req, err := http.NewRequest("GET", app.service.GenReqLast(), nil)
		if err != nil {
			app.errorLog.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := app.client.Do(req)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		newid, timestamp, transactions := app.BodyToData(data)

		if uint(newid) == lastid {
			time.Sleep(time.Second)
			continue
		}

		err = app.service.UpdRest(1)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		err = app.service.InsertTransaction(transactions, timestamp, 0)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		lastid = uint(newid)
		log.Println("A new block has been added ")
		time.Sleep(time.Second * time.Duration(t))
	}
}

func (app *application) InitDB(t int) {
	log.Println("Start initiating db")
	iter := app.GetLastID()
	lastid := iter - 1000
	time.Sleep(time.Second / time.Duration(t))

	for i := 0; i < 1000; i++ {
		req, err := http.NewRequest("GET", app.service.GenReq(lastid+int64(i)), nil)

		if err != nil {
			app.errorLog.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := app.client.Do(req)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		_, timestamp, transactions := app.BodyToData(data)

		err = app.service.InsertTransaction(transactions, timestamp, int(lastid)-i)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		time.Sleep(time.Second / time.Duration(t))
	}

	lastid = app.GetLastID()
	err := app.service.UpdRest(int(lastid - iter))
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func (app *application) Validation(t types.Transaction) bool {
	req, err := http.NewRequest("GET", app.service.GetReqTransactionID(t.ID), nil)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := app.client.Do(req)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	jsondata := struct {
		Transaction types.TransactionRaw `json:"result"`
	}{}

	err = json.Unmarshal(data, &jsondata)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	valt, err := jsondata.Transaction.RawToNormal("0x0")
	if err != nil {
		app.errorLog.Fatal(err)
	}

	if valt.ID == t.ID && valt.Value == t.Value && valt.Sender == t.Sender && valt.Receiver == valt.Receiver {
		return true
	}
	return false
}
