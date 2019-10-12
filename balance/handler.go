package balance

import (
	"time"

	"github.com/fberrez/scrooge/balance/backend"
	"github.com/fberrez/scrooge/scrooge"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Server struct {
	backend *backend.Backend
}

// New initializes a new server instance.
func New() (*Server, error) {
	backend, err := backend.New()
	if err != nil {
		return nil, err
	}

	return &Server{
		backend: backend,
	}, nil
}

// MakeTransaction updates the account corresponding to the given account ID
// by adding or substracted the account's balance by the given amount.
// It returns a response containing the account id and the updated balance.
func (s *Server) MakeTransaction(ctx context.Context, in *scrooge.Transaction) (*scrooge.Response, error) {
	log.WithField("timestamp", time.Now()).Info("Request 'make transaction' received")

	// parses the received string to uuid
	accountID, err := uuid.FromString(in.AccountID)
	if err != nil {
		return nil, err
	}

	balance, err := s.backend.UpdateAccount(accountID, in.Amount)
	if err != nil {
		return nil, err
	}

	return &scrooge.Response{
		AccountID: in.AccountID,
		Balance:   int64(balance),
	}, nil
}

func (s *Server) GetBalance(ctx context.Context, in *scrooge.Account) (*scrooge.Response, error) {
	log.WithField("timestamp", time.Now()).Info("Request 'get balance' received")

	// parses the received string to uuid
	accountID, err := uuid.FromString(in.AccountID)
	if err != nil {
		return nil, err
	}

	balance, err := s.backend.GetBalance(accountID)
	if err != nil {
		return nil, err
	}

	return &scrooge.Response{
		AccountID: in.AccountID,
		Balance:   int64(balance),
	}, nil

}
