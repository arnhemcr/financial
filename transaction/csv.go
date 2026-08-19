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
	"errors"
	"fmt"
	"slices"
	"strings"
)

const ModuleCSV = "mcsv" // The name of this module's CSV record format.

/*
ParseCSV parses this transaction from the fields of a CSV record according to the format.
It assumes the format is valid, which can be checked with [CSVRecordFormat.Validate].
If ParseCSV fails to parse the transaction, it returns the first error.
*/
func (t *Transaction) ParseCSV(fields []string, f CSVRecordFormat) error {
	if len(fields) != int(f.NFields) {
		return fmt.Errorf("ParseCSV: %w", errNFields)
	}

	// Prepend fields with an empty string, so a field whose index is zero has value empty string.
	fields = slices.Insert(fields, 0, "")

	err := t.parseRequired(fields, f)
	if err != nil {
		return fmt.Errorf("ParseCSV: %w", err)
	}

	err = t.parseOptional(fields, f)
	if err != nil {
		return fmt.Errorf("ParseCSV: %w", err)
	}

	return nil
}

// StringModuleCSV returns this transaction as a CSV record in this module's format (mcsv).
func (t Transaction) StringModuleCSV() string {
	fields := []string{t.Date, t.ThisAccount, t.OtherAccount, t.Code, t.Memo, t.AmountText, t.Currency}

	return strings.Join(fields, ",") + "\n"
}

var (
	errMemo        = errors.New("memo cannot be empty string")
	errNFields     = errors.New("unexpected number of fields in record")
	errThisAccount = fmt.Errorf("this account cannot be empty string (or %q)", DefaultOtherAccount)
)

func (t *Transaction) parseRequired(fields []string, f CSVRecordFormat) (err error) {
	t.AmountText, t.Amount, err = parseAmount(fields, f)
	if err != nil {
		return err
	}

	t.Date, err = ParseDate(fields[f.DateI], f.DateLayout)
	if err != nil {
		return err
	}

	t.Memo = fields[f.MemoI]
	if t.Memo == "" {
		return errMemo
	}

	t.OtherAccount = fields[f.OtherAccountI]
	if t.OtherAccount == "" {
		t.OtherAccount = DefaultOtherAccount
	}

	a := fields[f.ThisAccountI]

	switch {
	case t.ThisAccount == DefaultOtherAccount || a == DefaultOtherAccount:
		return errThisAccount
	case t.ThisAccount != "":
		// This account already has a value which takes precedence over its field.
	case a != "":
		t.ThisAccount = a
	default:
		return errThisAccount
	}

	return nil
}

func (t *Transaction) parseOptional(fields []string, f CSVRecordFormat) error {
	t.Code = fields[f.CodeI]

	if t.Currency != "" {
		// The existing currency value takes precedence over its field.
		return nil
	}

	c := fields[f.CurrencyI]
	if c == "" {
		return nil
	}

	if !IsLedgerCurrency(c) {
		return fmt.Errorf("%q %w", c, errCurrency)
	}

	t.Currency = c

	return nil
}
