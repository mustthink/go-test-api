# go-test-api

> Uses: etherscan API, MongoDB, Docker

### Handles requests to the address (default URL)
1) `localhost:8080/transactions` with tags 
    1) `from`, `to`, `id`, `block` - for filter transactions
    2) `page` - for pagination (20 transactions per page)
    3) `v`- for validation (comparing transactions from the database with the blockchain)
### How API works
> #### Flags for configuration settings
> 1) `url` - URL server address, `default: localhost:8080`
> 2) `connstr` - MongoDB connection string, `default: mongodb://root:example@127.0.0.1:27017`
> 3) `apikey` - your APIKEY for etherscan, `default: GICHEEBFZVYGAXX48VVWIWCNYKYEGDMEKZ`
> 4) `ethurl` - link to the API etherscan, `default: https://api.etherscan.io`
> 5) `reqps` - maximum requests per second to etherscan API `default: 5`
> 6) `btime` - average block time `default: 12`

> #### Launch and initiation 
> On startup the server reads all flags and then tries to connect to MongoDB.
> If successful, the server logs will indicate a successful connection.
> 
> Next, the database is checked for entries. 
> If there are no records, the database initiation is called
> At the etherscan API a request is called to specify the block number and then 1000 blocks are written.
> > With 5 requests per second, this process should take about 3 minutes and 20 seconds. However, it took me about 10 minutes.
> 
> Then the gorutin runs which cyclically makes a request to the etherscan API to write a new block to MongoDB. 
> If there is no new block after waiting for block time, then the request is made every second until it gets a new block.
> 
> At the same time, the server processes `/transactions` requests that should return a list of transactions.
> It reads filters and returns `502 status` if they are missing. 
> Otherwise, it forms query to MongoDB, checks pagination (if the `page` tag is present, it returns 20 transactions of the specified page). 
> If validation is also specified in the query, all received transactions are checked against the original ones using `eth_getTransactionByHash` query. 
> The result of the validation is displayed in the server console.