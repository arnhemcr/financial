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
It offers methods to:

  - parse a transaction from a [comma-separated values] (CSV) record in an arbitrary format
  - string a transaction to the package's own CSV record format or [Ledger] journal format

[comma-separated values]: https://www.ietf.org/rfc/rfc4180.txt
[Ledger]: https://ledger-cli.org
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

	if trn.OtherAccount == trn.ThisAccount {
		return errAccounts
	}

	return nil
}

var (
	errAccounts     = errors.New("this account cannot be the same as that account")
	errAmount       = errors.New("amount cannot be zero")
	errCreditDebit  = errors.New("credit and debit cannot both be empty string or both non-empty string")
	errMemo         = errors.New("memo cannot be empty string")
	errNFields      = errors.New("wrong number of fields")
	errOtherAccount = errors.New("other account cannot be empty string")
	errThisAccount  = errors.New("this account cannot be empty string")
)

/*
ParseDate returns the date parsed from str, according to the date format,
as a string and nil.
If it fails to parse the date, parseDate returns the error.
*/
func parseDate(str, format string) (string, error) {
	val, err := time.Parse(format, str)
	if err != nil {
		return "", fmt.Errorf("parsedate: %w", err)
	}

	return val.Format(time.DateOnly), nil
}
