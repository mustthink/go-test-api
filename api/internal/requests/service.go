package data

import "go.mongodb.org/mongo-driver/mongo"

type Service struct {
	Ethurl string
	Apikey string
	DB     *mongo.Client
}
