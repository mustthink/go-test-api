package handlers

import (
	"encoding/json"
	"github.com/mustthink/go-test-api/internal/types"
	"io"
	"net/http"
	"time"
)

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

		jsondata := struct {
			Result struct {
				Number       string                 `json:"number"`
				Timestamp    string                 `json:"timestamp"`
				Transactions []types.TransactionRaw `json:"transactions"`
			} `json:"result"`
		}{}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		err = json.Unmarshal(data, &jsondata)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		newid, err := types.ConvHexDec(jsondata.Result.Number)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		if uint(newid) == lastid {
			time.Sleep(time.Second)
			continue
		}

		err = app.service.UpdRest()
		if err != nil {
			app.errorLog.Fatal(err)
		}

		err = app.service.InsertTransaction(jsondata.Result.Transactions, jsondata.Result.Timestamp)
		if err != nil {
			app.errorLog.Fatal(err)
		}

		time.Sleep(time.Second * time.Duration(t))
	}
}
