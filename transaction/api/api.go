package api

import (
	"net/http"
	"time"

	"github.com/fberrez/scrooge/scrooge"
	"github.com/fberrez/scrooge/transaction/backend"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/utils/jujerr"
	log "github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
	"google.golang.org/grpc"
)

// API contains each part of the API settings.
type (
	API struct {
		// fizz is the http server.
		fizz *fizz.Fizz

		// backend is the backend client used to interact
		// with the database.
		backend *backend.Backend

		// conn is the grpc client connection.
		conn *grpc.ClientConn

		// scroogeClient is the client used to interact with scrooge.
		scroogeClient scrooge.ScroogeClient
	}
)

// New parses the config file and initializes the new API.
func New() (*API, error) {
	log.Info("Initializing API")

	f := fizz.New()

	backend, err := backend.New()
	if err != nil {
		return nil, err
	}

	api := &API{
		fizz:    f,
		backend: backend,
	}

	// API informations
	infos := &openapi.Info{
		Title:       "Scrooge's Transaction API",
		Description: "Scrooge's Transaction API is used to make transactions on Scrooge accounts.",
		Version:     "1.0.0",
	}

	// Defines groups of routes
	transactionGroup := f.Group("/transaction", "Transaction", "Group of transaction routes.")
	balanceGroup := f.Group("/balance", "Balance", "Group of balance routes.")
	unsecuredGroup := f.Group("/unsecured", "Unsecured", "Group of unsecured routes.")

	// Defines Unsecured group's routes
	unsecuredGroup.GET("/openapi.json", []fizz.OperationOption{
		fizz.Summary("Generates a Swagger documentation in JSON"),
		fizz.Description("Returns a Swagger JSON containing all informations about the API."),
	}, f.OpenAPI(infos, "json"))

	unsecuredGroup.GET("/health", []fizz.OperationOption{
		fizz.Summary("Returns the health status of the API."),
	}, tonic.Handler(api.health, http.StatusOK))

	// Defines transaction group's routes
	transactionGroup.POST("/", []fizz.OperationOption{
		fizz.Summary("Makes transaction"),
		fizz.Description("It makes a new transaction on a specific account. Then, the transaction is saved."),
	}, tonic.Handler(api.makeTransaction, http.StatusOK))

	transactionGroup.PUT("/", []fizz.OperationOption{
		fizz.Summary("Updates transaction"),
		fizz.Description("It updates a transaction."),
	}, tonic.Handler(api.updateTransaction, http.StatusOK))

	// Defines balance group's routes
	balanceGroup.GET("/:id", []fizz.OperationOption{
		fizz.Summary("Gets balance"),
		fizz.Description("It returns the balance of the targeted account."),
	}, tonic.Handler(api.getBalance, http.StatusOK))

	// Sets the error hook
	tonic.SetErrorHook(jujerr.ErrHook)

	return api, nil
}

func (a *API) NewGrpcConnection(address string) error {
	// makes a new grpc connection
	var err error
	a.conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	a.scroogeClient = scrooge.NewScroogeClient(a.conn)
	return nil
}

// ServeHTTP is the implementation of http.Handler.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"timestamp":   time.Now(),
		"remote_addr": r.RemoteAddr,
		"request":     r.RequestURI,
	}).Info("Request received.")

	a.fizz.ServeHTTP(w, r)
}

// Close closes the grpc connection.
func (a *API) Close() {
	a.conn.Close()
}
