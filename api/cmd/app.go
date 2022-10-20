package main

import (
	"flag"
	"github.com/mustthink/go-test-api/internal/db"
	"github.com/mustthink/go-test-api/internal/handlers"
	"log"
	"net/http"
	"os"
	"strconv"
)

func readFlags() (string, string, string, string, string, string) {
	var a = flag.String("url", "localhost:8080", "Server url")
	var b = flag.String("connstr", "mongodb://root:example@127.0.0.1:27017", "MongoDB url")
	var c = flag.String("apikey", "GICHEEBFZVYGAXX48VVWIWCNYKYEGDMEKZ", "API key for etherscan")
	var d = flag.String("ethapi", "https://api.etherscan.io", "Link to the etherscan api")
	var e = flag.String("reqps", "5", "Requests per second")
	var f = flag.String("btime", "12", "Approximate block time")
	flag.Parse()
	return *a, *b, *c, *d, *e, *f
}

func main() {
	addr, mongourl, apikey, ethurl, req, bt := readFlags()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	t, err := strconv.Atoi(req)
	if err != nil {
		errorLog.Fatal(err)
	}
	print(t)

	btime, err := strconv.Atoi(bt)
	if err != nil {
		errorLog.Fatal(err)
	}

	client, err := db.ConnClient(mongourl)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := handlers.NewApplication(errorLog, addr, ethurl, apikey, client)

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	go app.ScanBlocks(btime)

	log.Println("Starting Hosting on ", addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
