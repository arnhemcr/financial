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
It offers:

  - configurable parsing of [comma-separated values (CSV)] records into transactions
  - stringing of transactions into the module's CSV records or [Ledger] journal entries

[comma-separated values (CSV)]: https://www.ietf.org/rfc/rfc4180.txt
[Ledger]: https://ledger-cli.org
*/
package transaction

import (
	"errors"
	"fmt"
)

/*
A transaction represents a financial transaction.

A transaction transfers an amount of money from this account to another account
on a particular date.
*/
type Transaction struct {
	Amount       float64
	Currency     string
	Date         string
	Memo         string // This field is also known as the description.
	OtherAccount string // This field defaults to DefaultOtherAccount.
	ThisAccount  string
}

const DefaultOtherAccount = "Imbalance"

// IsValid reports whether this transaction is valid.
func (t Transaction) IsValid() bool {
	return t.Validate() == nil
}

/*
Validate returns nil if this transaction is valid.
If not, validate returns the first error.
*/
func (t Transaction) Validate() error {
	if t.Amount == 0.00 {
		return errAmount
	}

	err := validateDateOnly(t.Date)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if t.Memo == "" {
		return errMemo
	}

	if t.OtherAccount == "" {
		return errOtherAccount
	}

	if t.ThisAccount == "" {
		return errThisAccount
	}

	if t.OtherAccount == t.ThisAccount {
		return errAccounts
	}

	return nil
}

/*
StringFormat returns this transaction in the named format.
If the name is not known, stringFormat returns the empty string.
*/
func (t Transaction) StringFormat(name string) string {
	switch name {
	case CSV:
		return t.StringCSV()
	case Ledger:
		return t.StringLedger()
	default:
		return ""
	}
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
