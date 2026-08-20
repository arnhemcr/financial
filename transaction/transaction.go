/*
Copyright (C) 2026 Andrew Flint.

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
  - parsing a transaction from a [comma-separated values (CSV)] record;
    an instance of type CSVRecordFormat configures the parser for the record format
  - parsing some of a transaction's fields from a [Ledger] journal entry
  - stringing a transaction to either a Ledger journal entry or this module's CSV record (mcsv)

[comma-separated values (CSV)]: https://en.wikipedia.org/wiki/Comma-separated_values
[Ledger]: https://en.wikipedia.org/wiki/Ledger_(software)
*/
package transaction

import (
	"errors"
	"fmt"
)

/*
A Transaction represents a financial transaction:
the transfer of an amount of currency from one account to another on a date.
It is described by a memo and code, also called the description and transaction type.
A transaction belongs to an account called this account.
Optional fields may have the value empty string, while required fields must have non-zero values.
*/
type Transaction struct {
	// The amount is represented as a string and a floating-point number.
	AmountText   string
	Amount       float64
	Code         string // This field is optional.
	Currency     string // This field is optional.
	Date         string
	Memo         string
	OtherAccount string
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

/*
Validate returns nil if this transaction is valid.
If not, Validate returns the first error.
*/
func (t Transaction) Validate() error {
	_, a, err := parseDecimal(t.AmountText)
	if err != nil || t.Amount != a {
		return fmt.Errorf("Validate: %w", err)
	}

	// The code field is not validated because it is optional and free form.

	if t.Currency != "" && !IsLedgerCurrency(t.Currency) {
		return fmt.Errorf("Validate: %w", errCurrency)
	}

	_, err = ParseModuleDate(t.Date)
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	if t.Memo == "" {
		return fmt.Errorf("Validate: %w", errMemo)
	}

	if t.OtherAccount == "" {
		return fmt.Errorf("Validate: %w", errOtherAccount)
	}

	if t.ThisAccount == "" || t.ThisAccount == DefaultOtherAccount {
		return fmt.Errorf("Validate: %w", errThisAccount)
	}

	return nil
}

var (
	errMemo         = errors.New("memo cannot be empty string")
	errNFields      = errors.New("unexpected number of fields in record")
	errOtherAccount = errors.New("other account cannot be empty string")
	errThisAccount  = fmt.Errorf("this account cannot be empty string (or %q)", DefaultOtherAccount)
)
