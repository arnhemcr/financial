/*
Copyright (C) 2025 Andrew Flint.

This file is part of arnhemcr/financial.

Arnhemcr/financial is free software:
you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Arnhemcr/financial is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with arnhemcr/financial.
If not, see <https://www.gnu.org/licenses/>.
*/

/*
Package transaction represents financial transactions.
It provides methods to translate transactions to and from formats
used in files including account statements and financial journals.
The formats include comma-separated values (CSV) and [Ledger].

[Ledger]: https://ledger-cli.org "Ledger command-line accounting"

[Ledger 3 Manual]: https://ledger-cli.org/docs.html
*/
package transaction

import (
	"errors"
	"fmt"
	"time"
)

/*
A transaction represents a financial transaction.

A transaction transfers an amount of money from this account to another account
on a particular date.
*/
type Transaction struct {
	Amount       float64
	Date         string
	Memo         string // This field is also known as the description.
	OtherAccount string // This field defaults to DefaultOtherAccount.
	ThisAccount  string
}

const DefaultOtherAccount = "Imbalance"

// IsValid reports whether this transaction is valid.
func (trn Transaction) IsValid() bool {
	return trn.Validate() == nil
}

/*
Validate returns nil if this transaction is valid.
If not, validate returns the first error.
*/
func (trn Transaction) Validate() error {
	if trn.Amount == 0.00 {
		return errAmount
	}

	_, err := time.Parse(time.DateOnly, trn.Date)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if trn.Memo == "" {
		return errMemo
	}

	if trn.OtherAccount == "" {
		return errOtherAccount
	}

	if trn.ThisAccount == "" {
		return errThisAccount
	}

	return nil
}

var (
	errAmount       = errors.New("amount cannot be zero")
	errCreditDebit  = errors.New("credit and debit cannot both be empty string or non-empty string")
	errMemo         = errors.New("memo cannot be empty string")
	errNFields      = errors.New("wrong number of fields")
	errOtherAccount = errors.New("this account cannot be empty string")
	errThisAccount  = errors.New("this account cannot be empty string")
)

/*
ParseDate returns the date parsed from str, as a string, and nil.
If it fails to parse a date, parseDate returns an error.
*/
func parseDate(str, format string) (string, error) {
	val, err := time.Parse(format, str)
	if err != nil {
		return "", fmt.Errorf("parsedate: %w", err)
	}

	return val.Format(time.DateOnly), nil
}
