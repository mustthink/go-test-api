package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mustthink/go-test-api/internal/types"
	"io"
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
		Result struct {
			Number string `json:"number"`
		} `json:"result"`
	}{}

	err = json.Unmarshal(data, &jsondata)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	lastid, err := types.ConvHexDec(jsondata.Result.Number)
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
		fmt.Println("A new block has been added ")
		time.Sleep(time.Second * time.Duration(t))
	}
}

func (app *application) InitDB(t int) {
	var (
		first  = true
		lastid int64
		err    error
	)

	for i := 0; i < 1000; i++ {
		var req *http.Request

		if !first {
			req, err = http.NewRequest("GET", app.service.GenReq(lastid), nil)

		} else if first {
			req, err = http.NewRequest("GET", app.service.GenReqLast(), nil)
			first = false
		}

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
		lastid = newid

		err = app.service.InsertTransaction(transactions, timestamp, i)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		time.Sleep(time.Second / time.Duration(t))
	}

	lastid = app.GetLastID()
	err = app.service.UpdRest(int(lastid))
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
