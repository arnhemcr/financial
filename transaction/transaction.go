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

[Ledger 3 Manual]: https://ledger-cli.org/doc/ledger3.html
*/
package transaction

/*
A transaction represents a financial transaction.
It is the transfer of an amount of currency from one account to another.
The transfer takes place on a date.
It is described by a memo and code,
also known as the description and transaction type respectively.
A transaction belongs to an account known as this account.
*/
type Transaction struct {
	Amount       float64
	Code         string // can be empty string
	Currency     string // can be empty string
	Date         string
	Memo         string // or description
	OtherAccount string // defaults to "Imbalance"
	ThisAccount  string
}

const DefaultOtherAccount = "Imbalance" // The default value for other account.

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
