package api

import (
	"github.com/fberrez/scrooge/scrooge"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type (
	// AccountIn represents an account.
	// It is used as argument in controllers.
	AccountIn struct {
		ID string `path:"id" description:"Account ID"`
	}
)

// getBalance returns the balance of the account corresponding to the given id.
func (a *API) getBalance(c *gin.Context, in *AccountIn) (*ResponseBalanceOut, error) {
	response, err := a.scroogeClient.GetBalance(context.Background(), &scrooge.Account{
		AccountID: in.ID,
	})
	if err != nil {
		return nil, err
	}

	accountID, err := uuid.FromString(response.AccountID)
	if err != nil {
		return nil, err
	}

	return &ResponseBalanceOut{
		AccountID: accountID,
		Balance:   response.Balance,
	}, err

}
