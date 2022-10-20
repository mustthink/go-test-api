package data

import (
	"github.com/mustthink/go-test-api/internal/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Ethurl string
	Apikey string
	DB     *mongo.Client
}

func (s *Service) GenReq(id int64) string {
	return s.Ethurl + "/api?module=proxy&action=eth_getBlockByNumber&tag=" + types.ConvDecHex(id) + "&boolean=true&apikey=" + s.Apikey
}

func (s *Service) GenReqLast() string {
	return s.Ethurl + "/api?module=proxy&action=eth_getBlockByNumber&boolean=true&apikey=" + s.Apikey
}

func (s *Service) GetReqID() string {
	return s.Ethurl + "/api?module=proxy&action=eth_blockNumber&apikey=" + s.Apikey
}

func (s *Service) GetReqTransactionID(hash string) string {
	return s.Ethurl + "/api?module=proxy&action=eth_getTransactionByHash&txhash=" + hash + "&apikey=" + s.Apikey
}
