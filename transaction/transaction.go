/*
Copyright (C) 2025-2026 Andrew Flint.

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
  - parsing a transaction from a [CSV] record;
    an instance of type CSVRecordFormat configures the parser for the record format
  - parsing some of a transaction's fields from a [Ledger] journal entry
  - stringing a transaction to either a Ledger journal entry or
    this module's CSV record

[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
*/
package transaction

/*
A Transaction represents a financial transaction:
the transfer of an amount of currency from one account to another on a date.
It is described by a memo and code, also known as the description and transaction type respectively.
A transaction belongs to an account known as this account.
Optional fields may have the value empty string, while the rest must have non-zero values.
*/
type Transaction struct {
	Amount       float64
	Code         string // This field is optional. See also MirrorCode.
	Currency     string // This field is optional.
	Date         string
	Memo         string
	OtherAccount string // The default value of this field is DefaultOtherAccount.
	ThisAccount  string
}

const (
	DefaultOtherAccount = "Imbalance" // The default value for other account.
	MirrorCode          = "MT"        // The code value for a mirror transaction.
)

/*
StringFormat returns this transaction in the named format.
If the name is not known, StringFormat returns the empty string.
*/
func (t Transaction) StringFormat(name string) string {
	switch name {
	case Ledger, Ledger_:
		return t.StringLedger()
	case ModuleCSV, ModuleCSV_:
		return t.StringModuleCSV()
	default:
		return ""
	}
}
