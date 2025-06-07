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
  - stringing of transactions into [Ledger] journal entries or this module's CSV records

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
*/
package transaction

import "errors"

/*
A transaction represents a financial transaction.
It is the transfer of an amount of currency
from one account (this account) to another (other account)
with a memo on a date.
*/
type Transaction struct {
	Amount       float64
	Currency     string // can be empty string
	Date         string
	Memo         string
	OtherAccount string // defaults to DefaultOtherAccount
	ThisAccount  string
}

const DefaultOtherAccount = "Imbalance"

/*
StringFormat returns this transaction in the named format.
If the name is not known, stringFormat returns the empty string.
*/
func (t Transaction) StringFormat(name string) string {
	switch name {
	case CSV:
		return t.StringModuleCSV()
	case Ledger:
		return t.StringLedger()
	default:
		return ""
	}
}

var (
	errAccounts     = errors.New("this account cannot be the same as other account")
	errAmount       = errors.New("amount cannot be zero")
	errCreditDebit  = errors.New("credit and debit cannot both be empty string or both non-empty string")
	errMemo         = errors.New("memo cannot be empty string")
	errOtherAccount = errors.New("other account cannot be empty string")
	errThisAccount  = errors.New("this account cannot be empty string")
)
