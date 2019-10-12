package backend

import (
	"github.com/satori/go.uuid"
)

var (
	defaultCurrencyKey = "defaultCurrency"
)

// UpdateAccount updates the balance account by adding (or subtracted) the given amount.
// It returns the updated balance.
func (b *Backend) UpdateAccount(accountID uuid.UUID, amount int64) (int, error) {
	// prepares the query
	stmt, err := b.db.Prepare("UPDATE account SET balance = balance + $1 WHERE uuid = $2 RETURNING balance;")
	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	// executes the query and gets the updated balance
	var balance int
	if err = stmt.QueryRow(amount, accountID).Scan(&balance); err != nil {
		return -1, err
	}

	return balance, nil
}

// GetBalance returns the balance of the account corresponding to the given
// account id.
func (b *Backend) GetBalance(accountID uuid.UUID) (int, error) {
	// prepares the query
	rows, err := b.db.Query("SELECT balance FROM account WHERE uuid = $1 LIMIT 1;", accountID)
	if err != nil {
		return -1, err
	}

	// executes the query and gets the returned balance
	var balance int
	for rows.Next() {
		if err := rows.Scan(&balance); err != nil {
			return -1, err
		}
	}

	return balance, nil
}
