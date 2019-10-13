# Scrooge

<a href="https://goreportcard.com/report/github.com/fberrez/scrooge"><img src="https://goreportcard.com/badge/github.com/fberrez/scrooge"></a>

## Abstract

The main objective of this project is to test and play with the gRPC protocol.

## Project

Scrooge is a micro bank application where we can send transactions which modify the account balances in the database.

It uses two micro-services:

- balance: it manages accounts and their balance in the database. It exposes a gRPC interface and handles the received requests.
- transaction: it exposes a REST API used to create and update transactions. They are used to modify accounts balance.

## Deployment

> Note: the current configuration files have been edited for a docker environment. They must be edited before a local use.

```sh
# Local environment
## Downloads packages
$ dep ensure
## Runs services
$ export CONFIGURATION_FROM=file:transaction/config.yml; go run cmd/transaction/main.go
$ export CONFIGURATION_FROM=file:balance/config.yml; go run cmd/balance/main.go
## Builds binaries
$ make

# Builds docker images
$ docker build -t fberrez/scrooge-balance -f Dockerfile-balance .
$ docker build -t fberrez/scrooge-transaction -f Dockerfile-transaction .

# Creates directory for postgres data
$ mkdir containers/backend/data

# Runs docker images
$ docker-compose up -d
```

## API

A documentation can be find [here](https://app.swaggerhub.com/apis-docs/fberrez/Scrooge/1.0.0).

The database contains two accounts:

```sh
- account_id: 'a84f1a1b-d6eb-4819-be29-2055b8862094'
  balance: 100
- account_id: '7f87b08e-760b-43df-9af4-5354da34e7b4'
  balance: 500
```

### Example of requests:

```sh
# Makes a transaction
$ curl -L -XPOST -H 'Content-Type:application/json' --data '{
	"account_id": "7f87b08e-760b-43df-9af4-5354da34e7b4",
	"amount": 50,
	"currency": "eur"
}' localhost:3001/transaction
{"account_id":"7f87b08e-760b-43df-9af4-5354da34e7b4","balance":550,"transaction_id":"9a86ea42-6a3e-4c6a-bfec-503930bdf910","created_at":"2019-10-12T21:59:57.400141Z"}

# Get Balance
$ curl -L -XGET localhost:3001/balance/7f87b08e-760b-43df-9af4-5354da34e7b4
{"account_id":"7f87b08e-760b-43df-9af4-5354da34e7b4","balance":550}

# Update transaction
$ curl -L -XPUT -H 'Content-Type:application/json' --data '{
	"id": "9a86ea42-6a3e-4c6a-bfec-503930bdf910",
	"account_id": "7f87b08e-760b-43df-9af4-5354da34e7b4",
	"amount": -100
}' localhost:3001/transaction
{"account_id":"7f87b08e-760b-43df-9af4-5354da34e7b4","balance":400,"transaction_id":"9a86ea42-6a3e-4c6a-bfec-503930bdf910","created_at":"2019-10-12T21:59:57.400141Z"}
```

## Architecture

```sh
.
├── balance 					# balance is the microservice used to interact with account balances.
│   └── config.yml 				# balance configuration file
├── cmd 						# main folder, used to start microservices in your local environment
├── containers 					# contains containers data
│   └── backend
│       ├── data 				# contains postgres data
│       └── schema.sql 			# database schema executed on the first run
├── docker-compose.yml			
├── Dockerfile-balance			# Dockerfile used to build balance image
├── Dockerfile-transaction		# Dockerfile used to build transaction image
├── scrooge						# scrooge contains the proto file
│   ├── scrooge.pb.go			# "compiled" proto file
│   └── scrooge.proto			# proto file, used to make interact with balance from the outside
└── transaction					# transaction is the microservice used to makes transactions
    └── config.yml				# transction configuration file

```
