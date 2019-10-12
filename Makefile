.PHONY: transaction balance
	
BUILD_ARGS := -ldflags "-X github/fberrez/build.Version=0.0.3"
BINARY := ./output
TRANSACTION_MAIN_LOCATION := cmd/transaction
TRANSACTION_BINARY := $(BINARY)/transaction
BALANCE_MAIN_LOCATION := cmd/balance
BALANCE_BINARY := $(BINARY)/balance

all:
	go build -o $(TRANSACTION_BINARY) $(BUILD_ARGS) $(TRANSACTION_MAIN_LOCATION)/*.go
	go build -o $(BALANCE_BINARY) $(BUILD_ARGS) $(BALANCE_MAIN_LOCATION)/*.go

transaction:
	go build -o $(TRANSACTION_BINARY) $(BUILD_ARGS) $(TRANSACTION_MAIN_LOCATION)/*.go

balance:
	go build -o $(BALANCE_BINARY) $(BUILD_ARGS) $(BALANCE_MAIN_LOCATION)/*.go

clean:
	rm -rf $(BINARY)

test:
	go test -v -race -cover -bench=. -coverprofile=cover.profile ./...

fmt:
	for filename in $$(find . -path ./vendor -prune -o -name '*.go' -print); do \
		gofmt -w -l -s $$filename ;\
	done

