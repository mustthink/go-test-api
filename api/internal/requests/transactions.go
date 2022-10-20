package data

import (
	"context"
	"github.com/mustthink/go-test-api/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
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

func (s *Service) GetTransactions(sender, receiver, id, block, page string) ([]types.Transaction, error) {
	collection := s.DB.Database("test").Collection("transactions")

	filter := bson.D{}

	if sender != "" {
		filter = append(filter, primitive.E{Key: "sender", Value: sender})
	}
	if receiver != "" {
		filter = append(filter, primitive.E{Key: "receiver", Value: sender})
	}
	if id != "" {
		idint, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		filter = append(filter, primitive.E{Key: "id", Value: idint})
	}
	if block != "" {
		blockint, err := strconv.Atoi(block)
		if err != nil {
			return nil, err
		}
		filter = append(filter, primitive.E{Key: "block", Value: blockint})
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	results := []types.Transaction{}

	for cur.Next(context.TODO()) {

		var elem types.Transaction
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())

	if page != "" {
		pageint, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}
		n := len(results)
		pageint -= 1
		pageint *= 20 //20 is the number of transactions per page
		if n > pageint && pageint >= 0 {
			if n > (pageint + 20) {
				return results[pageint:(pageint + 20)], nil
			}
			return results[pageint:], nil
		}
		return nil, nil
	}
	return results, nil
}
