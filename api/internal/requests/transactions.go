package data

import (
	"context"
	"github.com/mustthink/go-test-api/internal/types"
)

func (s *Service) InsertTransaction(t []types.TransactionRaw) error {
	collection := s.DB.Database("test").Collection("transactions")

	transactions := make([]interface{}, len(t))
	for i, v := range t {
		try, err := v.RawToNormal()
		if err != nil {
			return err
		}
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
