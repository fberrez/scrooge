package api

import (
	"time"

	"github.com/fberrez/scrooge/scrooge"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

type (
	// TransactionIn represents a transaction.
	// It is used as argument in controllers.
	TransactionIn struct {
		// AccountID is the id of the account targeted by the transaction.
		AccountID uuid.UUID `json:"account_id" description:"Targeted Account ID" validate:"required"`
		// Amount is the transaction amount.
		Amount int64 `json:"amount" description:"Transaction amount" validate:"required"`
		// Description is the transaction description.
		Description string `json:"description" description:"Transaction description"`
		// Currency is the transaction currency.
		Currency string `json:"currency" description:"Transaction currency"`
	}

	// UpdateTransactionIn represents a transaction.
	// It is used as argument in the update controller.
	UpdateTransactionIn struct {
		// ID is the transaction ID.
		ID uuid.UUID `json:"id" description:"Transaction ID" validate="required"`
		// AccountID is the id of the account targeted by the transaction.
		AccountID uuid.UUID `json:"account_id" description:"Targeted Account ID" validate:"required"`
		// Amount is the transaction amount.
		Amount int64 `json:"amount" description:"Transaction amount" validate:"required"`
	}

	// ResponseTransactionOut represents a response after a transaction is added or updated.
	ResponseTransactionOut struct {
		// AccountID is the id of the account targeted by the transaction.
		AccountID uuid.UUID `json:"account_id" description:"Targeted Account ID"`
		// Balance is the targeted account balance.
		Balance int64 `json:"balance" description:"Targeted Account balance (updated after the transaction)"`
		// TransactionID is the transaction id.
		TransactionID uuid.UUID `json:"transaction_id" description:"Transaction ID"`
		// CreatedAt is the transaction creation date.
		CreatedAt time.Time `json:"created_at" description:"Transaction creation date"`
	}

	// ResponseBalanceOut is used to give a reponse on a balance route.
	ResponseBalanceOut struct {
		// AccountID is the targeted account id.
		AccountID uuid.UUID `json:"account_id" description:"Targeted Account ID"`
		// Balance is the targeted account balance.
		Balance int64 `json:"balance" description:"Targeted Account balance (updated after the transaction)"`
	}
)

// makeTransaction is the controller of POST /transaction.
// It makes a transaction on the targeted account and returns
// its new balance value.
func (a *API) makeTransaction(c *gin.Context, transaction *TransactionIn) (*ResponseTransactionOut, error) {
	response, err := a.scroogeClient.MakeTransaction(context.Background(), &scrooge.Transaction{
		AccountID: transaction.AccountID.String(),
		Amount:    transaction.Amount,
	})
	if err != nil {
		return nil, err
	}

	newTransaction, err := a.backend.SaveTransaction(transaction.AccountID, transaction.Amount, transaction.Description, transaction.Currency)
	if err != nil {
		return nil, err
	}

	accountID, err := uuid.FromString(response.AccountID)
	if err != nil {
		return nil, err
	}

	return &ResponseTransactionOut{
		AccountID:     accountID,
		Balance:       response.Balance,
		TransactionID: newTransaction.ID,
		CreatedAt:     newTransaction.CreatedAt,
	}, nil
}

// updateTransaction updates the transaction corresponding to the given ID.
func (a *API) updateTransaction(c *gin.Context, transaction *UpdateTransactionIn) (*ResponseTransactionOut, error) {
	updatedTransaction, err := a.backend.UpdateTransaction(transaction.ID, transaction.AccountID, transaction.Amount)
	if err != nil {
		return nil, err
	}

	// cancel the previous transaction
	_, err = a.scroogeClient.MakeTransaction(context.Background(), &scrooge.Transaction{
		AccountID: updatedTransaction.Previous.AccountID.String(),
		Amount:    updatedTransaction.Previous.Amount * -1,
	})
	if err != nil {
		return nil, err
	}

	// apply the updated transaction
	response, err := a.scroogeClient.MakeTransaction(context.Background(), &scrooge.Transaction{
		AccountID: transaction.AccountID.String(),
		Amount:    transaction.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &ResponseTransactionOut{
		AccountID:     updatedTransaction.New.AccountID,
		Balance:       response.Balance,
		TransactionID: updatedTransaction.ID,
		CreatedAt:     updatedTransaction.Previous.CreatedAt,
	}, nil
}
