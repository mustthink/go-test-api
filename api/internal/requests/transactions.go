package data

import (
	"context"
	"github.com/mustthink/go-test-api/internal/types"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) InsertTransaction(t []types.TransactionRaw, tstamp string, n int) error {
	collection := s.DB.Database("test").Collection("transactions")

	transactions := make([]interface{}, len(t))
	for i, v := range t {
		try, err := v.RawToNormal(tstamp)
		if err != nil {
			return err
		}
		try.NumVer = n
		transactions[i] = *try
	}

	//итак: я не понимаю почему структура которая имплементирует интерфейс - может быть параметром в виде интерфейса
	//но при этом слайс структур не может :/
	_, err := collection.InsertMany(context.TODO(), transactions)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdRest(n int) error {
	collection := s.DB.Database("test").Collection("transactions")

	filter := bson.D{}

	update := bson.D{
		{"$inc", bson.D{
			{"numver", n},
		}},
	}

	_, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
