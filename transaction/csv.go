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

package transaction

import (
	"fmt"
	"slices"
	"strings"
)

const ModuleCSV = "mcsv" // The name of this module's CSV record format.

/*
ParseCSV parses this transaction from the fields of a CSV record according to the format.
It assumes the format is valid, which can be checked with [CSVRecordFormat.Validate].
If ParseCSV fails to parse a valid transaction, it returns the first error.
*/
func (t *Transaction) ParseCSV(fields []string, f CSVRecordFormat) (err error) {
	n, m := int(f.NFields), len(fields)
	if n != m {
		return fmt.Errorf("ParseCSV: expect %v fields not %v", n, m)
	}

	// Prepend fields with an empty string so fields with index zero have value empty string.
	fields = slices.Insert(fields, 0, "")

	t.AmountText, t.Amount, err = parseAmount(fields, f)
	if err != nil {
		return fmt.Errorf("ParseCSV: %w", err)
	}

	t.Code = fields[f.CodeI]

	// If currency already has a value, it takes precedence over its field.
	if t.Currency == "" {
		t.Currency = fields[f.CurrencyI]
	}

	if t.Currency != "" && !IsLedgerCurrency(t.Currency) {
		return fmt.Errorf("ParseCSV: %w not %q", errCurrency, t.Currency)
	}

	t.Date, err = ParseDate2(fields[f.DateI], f.DateLayout)
	if err != nil {
		return fmt.Errorf("ParseCSV: %w", err)
	}

	t.Memo = fields[f.MemoI]
	if t.Memo == "" {
		return fmt.Errorf("ParseCSV: %w not %q", errMemo, t.Memo)
	}

	t.OtherAccount = fields[f.OtherAccountI]
	if t.OtherAccount == "" {
		t.OtherAccount = DefaultOtherAccount
	}

	// If this account already has a value, it takes precedence over its field.
	if t.ThisAccount == "" {
		t.ThisAccount = fields[f.ThisAccountI]
	}

	if t.ThisAccount == "" || t.ThisAccount == DefaultOtherAccount {
		return fmt.Errorf("ParseCSV: %w not %q", errThisAccount, t.ThisAccount)
	}

	if t.ThisAccount == t.OtherAccount {
		return fmt.Errorf("ParseCSV: %w not both %q", errAccounts, t.ThisAccount)
	}

	return nil
}

// StringModuleCSV returns this transaction formatted as this module's CSV record (mcsv).
func (t Transaction) StringModuleCSV() string {
	vs := []string{t.Date, t.ThisAccount, t.OtherAccount, t.Code, t.Memo, t.AmountText, t.Currency}

	return strings.Join(vs, ",") + "\n"
}
