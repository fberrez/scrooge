package backend

import (
	"time"

	"github.com/ovh/configstore"
	uuid "github.com/satori/go.uuid"
)

type (
	// Transaction represents a transaction.
	Transaction struct {
		// ID is the transaction ID.
		ID uuid.UUID `json:"id" description:"Transaction ID"`
		// Account is the id of the account targeted by the transaction.
		AccountID uuid.UUID `json:"account_id" description:"Targeted Account ID"`
		// Amount is the transaction amount.
		Amount int64 `json:"amount" description:"Transaction amount"`
		// CreatedAt is the transaction creation date.
		CreatedAt time.Time `json:"created_at" description:"Transaction creation date"`
	}

	// UpdatedTransaction reprends a transaction update.
	UpdatedTransaction struct {
		// ID is the transaction ID.
		ID uuid.UUID `json:"id" description:"Transaction ID"`
		// Previous is the pre-update transaction.
		Previous *Transaction `json:"previous" description:"Previous transaction replaced during the update"`
		// New is the post-update transaction.
		New *Transaction `json:"new" description:"New transaction inserted during the update"`
	}
)

var (
	defaultCurrencyKey = "defaultCurrency"
)

// SaveTransaction adds a new transaction to the database.
// It returns some data about the inserted transaction.
func (b *Backend) SaveTransaction(accountID uuid.UUID, amount int64, description string, currency string) (*Transaction, error) {
	defaultCurrency, err := configstore.GetItemValue(defaultCurrencyKey)
	if err != nil {
		return nil, err
	}

	if len(currency) == 0 {
		currency = defaultCurrency
	}

	// prepares the query
	stmt, err := b.db.Prepare("INSERT INTO transaction(account_id, amount, created_at, description, currency) VALUES ($1, $2, $3, $4, $5) RETURNING uuid, created_at;")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	// executes the query
	transaction := Transaction{}
	if err = stmt.QueryRow(accountID, amount, time.Now(), description, currency).Scan(&transaction.ID, &transaction.CreatedAt); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// UpdateTransaction updates the transaction corresponding to the given UUID.
// It returns the previous transaction account id and amount.
func (b *Backend) UpdateTransaction(transactionID uuid.UUID, accountID uuid.UUID, amount int64) (*UpdatedTransaction, error) {
	// prepares the query
	stmt, err := b.db.Prepare("UPDATE transaction SET account_id = $1, amount = $2 FROM transaction old WHERE old.uuid = $3 RETURNING old.account_id, old.amount, old.created_at;")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	// executes the query and gets the old values
	previous := Transaction{}
	if err = stmt.QueryRow(accountID, amount, transactionID).Scan(&previous.AccountID, &previous.Amount, &previous.CreatedAt); err != nil {
		return nil, err
	}

	return &UpdatedTransaction{
		ID: transactionID,
		Previous: &Transaction{
			AccountID: previous.AccountID,
			Amount:    previous.Amount,
			CreatedAt: previous.CreatedAt,
		},
		New: &Transaction{
			AccountID: accountID,
			Amount:    amount,
		},
	}, nil
}
