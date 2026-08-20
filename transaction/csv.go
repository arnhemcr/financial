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
func (t *Transaction) ParseCSV(fields []string, f CSVRecordFormat) error {
	if len(fields) != int(f.NFields) {
		return fmt.Errorf("ParseCSV: %w", errNFields)
	}

	// Prepend fields with an empty string, so a field whose index is zero has value empty string.
	fields = slices.Insert(fields, 0, "")

	t.AmountText, t.Amount, _ = parseAmount(fields, f)

	t.Code = fields[f.CodeI]

	// If currency has a value, it takes precedence over its field.
	if t.Currency == "" {
		t.Currency = fields[f.CurrencyI]
	}

	t.Date, _ = ParseDate(fields[f.DateI], f.DateLayout)

	t.Memo = fields[f.MemoI]

	t.OtherAccount = fields[f.OtherAccountI]
	if t.OtherAccount == "" {
		t.OtherAccount = DefaultOtherAccount
	}

	// If this account has a value, it takes precedence over its field.
	if t.ThisAccount == "" {
		t.ThisAccount = fields[f.ThisAccountI]
	}

	err := t.Validate()
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	return nil
}

// StringModuleCSV returns this transaction as a CSV record in this module's format (mcsv).
func (t Transaction) StringModuleCSV() string {
	fields := []string{t.Date, t.ThisAccount, t.OtherAccount, t.Code, t.Memo, t.AmountText, t.Currency}

	return strings.Join(fields, ",") + "\n"
}
