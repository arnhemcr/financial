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

package transaction

import (
	"encoding/xml"
	"errors"
	"os"
	"slices"
	"strings"
	"time"
)

const CSV = "csv" // The name of this module's CSV record format.

// A CSVRecordFormat defines the format of CSV records representing financial transactions.
type CSVRecordFormat struct {
	NFields uint8 // The number of fields in a record.

	/*
		The indexes of fields in the record.
		Some fields are required but others are optional.
		The index for a required field is between 1 and NFields inclusive.
		If the record does not contain an optional field, the field's index is zero.
	*/
	// Either amount, or both credit and debit are required.
	AmountI, CreditI, DebitI uint8
	CurrencyI                uint8 // optional
	DateI                    uint8 // required
	MemoI                    uint8 // required
	OtherAccountI            uint8 // optional
	ThisAccountI             uint8 // optional

	// The layout of the date field in the records e.g. "02/01/2006".
	DateLayout string
}

/*
ReadCSVFormat returns the first CSV record format read from the named file.
If it fails to get a format, ReadCSVFormat returns the first error.
*/
func ReadCSVFormat(fileName string) (CSVRecordFormat, error) {
	var crf CSVRecordFormat

	bs, err := os.ReadFile(fileName)
	if err != nil {
		return crf, err
	}

	err = xml.Unmarshal(bs, &crf)
	if err != nil {
		return crf, err
	}

	return crf, nil
}

// GetModuleCSVFormat returns this module's CSV record format.
func GetModuleCSVFormat() CSVRecordFormat {
	return CSVRecordFormat{
		NFields:       6,
		DateI:         1,
		ThisAccountI:  2,
		OtherAccountI: 3,
		MemoI:         4,
		AmountI:       5,
		CurrencyI:     6,
		DateLayout:    "2006-01-02",
	}
}

/*
ParseCSV parses this transaction from fields according to the CSV record format.
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

	var err error

	t.Date, err = parseDate(fs[crf.DateI], crf.DateLayout)
	if err != nil {
		return err
	}

	t.Amount, err = parseValue(fs, crf)
	if err != nil {
		return err
	} else if t.Amount == 0.00 {
		return errAmount
	}

	t.Currency = fs[crf.CurrencyI]

	t.Memo = fs[crf.MemoI]
	if t.Memo == "" {
		return errMemo
	}

	t.OtherAccount = fs[crf.OtherAccountI]
	if t.OtherAccount == "" {
		t.OtherAccount = DefaultOtherAccount
	}

	switch {
	case t.ThisAccount != "":
		// This account already has a value, so ignore the this account field.
	case fs[crf.ThisAccountI] != "":
		t.ThisAccount = fs[crf.ThisAccountI]
	default:
		return errThisAccount
	}

	return nil
}

// StringCSV returns this transaction as a CSV record in this module's format.
func (t Transaction) StringCSV() string {
	a := stringAmount(t.Amount)
	fs := []string{t.Date, t.ThisAccount, t.OtherAccount, t.Memo, a, t.Currency}

	return strings.Join(fs, ",") + "\n"
}

/*
Validate returns nil if this CSV record format is valid.
If not, validate returns the first error.
*/
func (crf CSVRecordFormat) Validate() error {
	d, _ := time.Parse(crf.DateLayout, crf.DateLayout)
	if d.Format(time.DateOnly) != time.DateOnly {
		return errDateLayout
	}

	if crf.NFields < minNFields || maxNFields < crf.NFields {
		return errNFieldsRange
	}

	err := crf.areIndexesValid()
	if err != nil {
		return err
	}

	if crf.DateI == 0 {
		return errDateI
	}

	if crf.MemoI == 0 {
		return errMemoI
	}

	err = crf.areOptionsValid()
	if err != nil {
		return err
	}

	return nil
}

const (
	// The inclusive limits for the number of fields in a CSV record.
	minNFields = 3 // date, memo and amount
	maxNFields = 20
)

var (
	errAmountOpt    = errors.New("amount field index, or credit and debit indexes cannot both be zero")
	errDateI        = errors.New("date field index cannot be zero")
	errDateLayout   = errors.New("date layout in CSV record must be Go style e.g. \"02/01/2006\"")
	errIndexUnique  = errors.New("field indexes cannot share a non-zero value")
	errIndexRange   = errors.New("field index is out of range")
	errMemoI        = errors.New("memo field index cannot be zero")
	errNFields      = errors.New("unexpected number of fields in CSV record")
	errNFieldsRange = errors.New("number of fields in CSV record is out of range")
)

/*
AreIndexesValid returns nil if the field indexes in this CSV record format are valid.
It assumes the number of fields in the format is in range.
All indexes must be <= nFields, and all non-zero indexes must be unique.
If not, areIndexesValid returns the first error.
*/
func (crf CSVRecordFormat) areIndexesValid() error {
	is := [...]uint8{crf.AmountI, crf.CreditI, crf.DateI, crf.DebitI,
		crf.MemoI, crf.OtherAccountI, crf.ThisAccountI}

	var inUse [maxNFields + 1]bool

	for _, i := range is {
		switch {
		case crf.NFields < i:
			return errIndexRange
		case i == 0:
			// These CSV records do not contain this field.
		case inUse[i]:
			return errIndexUnique
		default:
			inUse[i] = true
		}
	}

	return nil
}

/*
AreOptionsValid returns nil if the combination of options is valid.
If not, areOptionsValid returns the error.
*/
func (crf CSVRecordFormat) areOptionsValid() error {
	switch {
	case crf.AmountI != 0:
		return nil
	case crf.CreditI != 0 && crf.DebitI != 0:
		return nil
	default:
		return errAmountOpt
	}
}

/*
ParseValue parses the value of a transaction from either the amount, credit or debit fields.
If it fails to parse the value, parseValue returns the error.
*/
func parseValue(fields []string, crf CSVRecordFormat) (float64, error) {
	a, c, d := fields[crf.AmountI], fields[crf.CreditI], fields[crf.DebitI]

	switch {
	case a != "":
		return parseAmount(a)
	case c != "" && d == "":
		return parsePositiveAmount(c)
	case d != "" && c == "":
		v, err := parsePositiveAmount(d)

		return v * -1.00, err
	default:
		return 0.00, errCreditDebit
	}
}
