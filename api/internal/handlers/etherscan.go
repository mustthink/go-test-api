package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mustthink/go-test-api/internal/types"
	"io"
	"log"
	"net/http"
	"time"
)

func genReq(url, key string) string {
	return url + "api?module=proxy&action=eth_getBlockByNumber&boolean=true&apikey=" + key
}

func (app *application) ScanBlocks(t int) {
	//var lastid int
	ticker := time.NewTicker(time.Second / time.Duration(t))
	for _ = range ticker.C {
		req, err := http.NewRequest("GET", "https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&boolean=true&apikey=GICHEEBFZVYGAXX48VVWIWCNYKYEGDMEKZ", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := app.client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		jsondata := struct {
			Result struct {
				Number       int                    `json:"number"`
				Transactions []types.TransactionRaw `json:"transactions"`
			} `json:"result"`
		}{}

		fmt.Println(resp.Body)
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return
		}

		err = json.Unmarshal(data, &jsondata)

		for _, v := range jsondata.Result.Transactions {
			fmt.Println(v)
		}

	}
}
