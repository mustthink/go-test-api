FROM golang:1.17

WORKDIR .

COPY go.mod go.sum service/
RUN cd service && go mod download && go mod verify

COPY . service/

RUN cd service && go mod tidy && go build cmd/app.go

CMD service/app -url="172.16.1.3:8081" -connstr="mongodb://root:example@172.16.1.2:27017"
