package handlers

import (
	"context"
	"fmt"
	"github.com/mustthink/go-test-api/internal/requests"
	"github.com/mustthink/go-test-api/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"runtime/debug"
)

type application struct {
	errorLog *log.Logger
	url      string
	service  *data.Service
	client   *http.Client
}

func NewApplication(errorLog *log.Logger, url, eth, key string, client *mongo.Client) *application {
	return &application{
		errorLog: errorLog,
		url:      url,
		service: &data.Service{
			Ethurl: eth,
			Apikey: key,
			DB:     client,
		},
		client: &http.Client{},
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) CheckDB() (bool, error) {
	collection := app.service.DB.Database("test").Collection("transactions")

	var result types.Transaction
	filter := bson.D{}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != mongo.ErrNoDocuments {
		return true, err
	} else if err == nil {
		return true, nil
	}

	return false, nil
}
