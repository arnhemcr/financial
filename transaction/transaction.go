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
Package transaction represents financial transactions as instances of type Transaction.
It offers:

  - parsing a transaction from a CSV record;
    an instance of type CSVRecordFormat configures the parser for the record format
  - stringing a transaction to either a Ledger journal entry or this module's CSV record

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
*/
package transaction

/*
A transaction represents a financial transaction.
It is the transfer of an amount of currency from one account to another on a date.
It is described by a memo and code, also known as the description and type respectively.
A transaction belong to an account known as this account.
*/
type Transaction struct {
	Amount       float64
	Code         string // or transaction type, can be empty string
	Currency     string // can be empty string
	Date         string
	Memo         string // or description
	OtherAccount string // defaults to "Imbalance"
	ThisAccount  string
}

const DefaultOtherAccount = "Imbalance" // The default value for other account.

/*
StringFormat returns this transaction in the named format.
If the name is not known, StringFormat returns the empty string.
*/
func (t Transaction) StringFormat(name string) string {
	switch name {
	case Ledger:
		return t.StringLedger()
	case ModuleCSV:
		return t.StringModuleCSV()
	default:
		return ""
	}
}
