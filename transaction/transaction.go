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
Package transaction represents financial transactions.
It offers:
  - parsing an instance of type Transaction from a [comma-separated values (CSV)] record;
    the parser is configured by an instance of type CSVRecordFormat
  - formatting a Transaction to either this module's CSV record (mcsv) or a [Ledger] journal entry

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
It is described by a memo and code also called the description and transaction type.
A transaction belongs to an account called this account.
Required fields must have non-zero values; optional fields may be empty string.
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
	_, a, err := parseNonZeroDecimal(t.AmountText)
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	if t.Amount != a {
		return fmt.Errorf("Validate: %w but %q and %v do not", errAmountMismatch, t.AmountText, t.Amount)
	}

	// The code is not validated because it is optional and free form.

	// The currency is optional.
	if t.Currency != "" && !IsLedgerCurrency(t.Currency) {
		return fmt.Errorf("Validate: %w not %q", errCurrency, t.Currency)
	}

	_, err = ParseModuleDate(t.Date)
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	if t.Memo == "" {
		return fmt.Errorf("Validate: %w not %q", errMemo, t.Memo)
	}

	if t.OtherAccount == "" {
		return fmt.Errorf("Validate: %w not %q", errOtherAccount, t.OtherAccount)
	}

	if t.ThisAccount == "" || t.ThisAccount == DefaultOtherAccount {
		return fmt.Errorf("Validate: %w not %q", errThisAccount, t.ThisAccount)
	}

	if t.ThisAccount == t.OtherAccount {
		return fmt.Errorf("Validate: %w not both %q", errAccounts, t.ThisAccount)
	}

	return nil
}

var (
	errAccounts     = errors.New("expect this and other account to be different")
	errMemo         = errors.New("expect memo string")
	errOtherAccount = errors.New("expect other account string")
	errThisAccount  = errors.New("expect this account string")
)
