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

package transaction

import (
	"errors"
	"slices"
	"strings"
)

const (
	ModuleCSV = "mcsv" // The name of this module's CSV record format.
)

/*
ParseCSV parses this transaction from the CSV record fields according to the format.
It assumes the format is valid.
If ParseCSV fails to parse the transaction, it returns the first error.
*/
func (t *Transaction) ParseCSV(fields []string, crf CSVRecordFormat) error {
	if len(fields) != int(crf.NFields) {
		return errNFields
	}

	/*
		Prepend fields with an empty string,
		so a field whose index is zero has value empty string.
	*/
	fs := slices.Insert(fields, 0, "")

	err := t.parseRequired(fs, crf)
	if err != nil {
		return err
	}

	t.parseOptional(fs, crf)

	return nil
}

// StringModuleCSV returns this transaction as this module's CSV record.
func (t Transaction) StringModuleCSV() string {
	a := stringAmount(t.Amount)
	fs := []string{t.Date, t.ThisAccount, t.OtherAccount, t.Code, t.Memo, a, t.Currency}

	return strings.Join(fs, ",") + "\n"
}

var (
	errMemo        = errors.New("memo cannot be empty string")
	errNFields     = errors.New("unexpected number of fields in CSV record")
	errThisAccount = errors.New(
		"this account cannot be empty string or \"" + DefaultOtherAccount + "\"")
)

func (t *Transaction) parseRequired(fields []string, crf CSVRecordFormat) error {
	var err error

	t.Amount, err = parseAmount(fields, crf)
	if err != nil {
		return err
	}

	t.Date, err = parseDate(fields[crf.DateI], crf.DateLayout)
	if err != nil {
		return err
	}

	t.Memo = fields[crf.MemoI]
	if t.Memo == "" {
		return errMemo
	}

	t.OtherAccount = fields[crf.OtherAccountI]
	if t.OtherAccount == "" {
		t.OtherAccount = DefaultOtherAccount
	}

	a := fields[crf.ThisAccountI]

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

func (t *Transaction) parseOptional(fields []string, crf CSVRecordFormat) {
	t.Code = fields[crf.CodeI]

	// Currency may already have a value which takes precedence over its field.
	if t.Currency == "" {
		t.Currency = fields[crf.CurrencyI]
	}
}
