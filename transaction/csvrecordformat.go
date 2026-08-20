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
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"time"
)

// A CSVRecordFormat defines the format of CSV records representing financial transactions.
type CSVRecordFormat struct {
	NFields uint8 // The number of fields in a record.

	/*
		The indexes of fields in the record.
		Some fields are required, while the rest are optional.
		The index for a required field is between one and NFields inclusive.
		If an optional field is not contained in a record, its index is zero.
	*/
	// Either amount, or credit and debit are required.
	AmountI         uint8
	CreditI, DebitI uint8
	CurrencyI       uint8
	CodeI           uint8
	DateI           uint8 // This field is required.
	MemoI           uint8 // This field is required.
	OtherAccountI   uint8
	ThisAccountI    uint8

	// The Go-style date layout in the record e.g. "01/02/2006".
	DateLayout string
}

/*
NewCSVRecordFormat returns the CSV record format read from the named XML file.
Fields in the format default to zero except DateLayout which defaults to "2006-02-01".
If it fails to read a valid format, NewCSVRecordFormat returns the first error.
*/
func NewCSVRecordFormat(fileName string) (f CSVRecordFormat, err error) {
	bs, err := os.ReadFile(fileName)
	if err != nil {
		return f, fmt.Errorf("NewCSVRecordFormat: %w", err)
	}

	err = xml.Unmarshal(bs, &f)
	if err != nil {
		return f, fmt.Errorf("NewCSVRecordFormat: %w", err)
	}

	if f.DateLayout == "" {
		f.DateLayout = time.DateOnly
	}

	err = f.Validate()
	if err != nil {
		return f, fmt.Errorf("NewCSVRecordFormat: %w", err)
	}

	return f, nil
}

// NewModuleCSVRecordFormat returns this module's CSV record format (mcsv).
func NewModuleCSVRecordFormat() CSVRecordFormat {
	return CSVRecordFormat{
		NFields: 7,

		DateI:         1,
		ThisAccountI:  2,
		OtherAccountI: 3,
		CodeI:         4,
		MemoI:         5,
		AmountI:       6,
		CurrencyI:     7,

		DateLayout: time.DateOnly,
	}
}

/*
Validate returns nil if this CSV record format is valid.
If not, Validate returns the first error.
*/
func (f CSVRecordFormat) Validate() error {
	n := f.NFields
	if n < minNFields || maxNFields < n {
		return fmt.Errorf("Validate: %w", errNFieldsRange)
	}

	err := f.validateIndexes()
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	err = f.validateOptions()
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	if !IsDateLayout(f.DateLayout) {
		return fmt.Errorf("Validate: %q %w", f.DateLayout, ErrDateLayout)
	}

	return nil
}

const (
	// The inclusive limits for the number of fields in a CSV record.
	minNFields = 3 // date, memo and amount
	maxNFields = 20
)

var (
	errAmountOption = errors.New("amount field index, or credit and debit indexes cannot be zero")
	errDateI        = errors.New("date field index cannot be zero")
	errIndexUnique  = errors.New("field indexes cannot share a non-zero value")
	errIndexRange   = errors.New("field index is out of range")
	errMemoI        = errors.New("memo field index cannot be zero")
	errNFieldsRange = errors.New("number of fields is out of range")
)

/*
ValidateIndexes returns nil if the field indexes in this CSV record format are valid.
Indexes must be between zero and NFields inclusive.
Required indexes must be non-zero and unique.
If the indexes are not valid, validateIndexes returns the first error.
*/
func (f CSVRecordFormat) validateIndexes() error {
	is := [...]uint8{f.AmountI, f.CodeI, f.CreditI, f.CurrencyI,
		f.DateI, f.DebitI, f.MemoI, f.OtherAccountI, f.ThisAccountI}

	var used [maxNFields + 1]bool

	for _, i := range is {
		switch {
		case i == 0:
			// This field is not contained in CSV records of this format.
		case f.NFields < i:
			return errIndexRange
		case used[i]:
			return errIndexUnique
		default:
			used[i] = true
		}
	}

	switch {
	case f.DateI == 0:
		return errDateI
	case f.MemoI == 0:
		return errMemoI
	default:
		return nil
	}
}

/*
ValidateOptions returns nil if the combination of optional field indexes in this CSV record format is valid.
If not, validateOptions returns the error.
*/
func (f CSVRecordFormat) validateOptions() error {
	switch {
	case f.AmountI != 0:
		return nil
	case f.CreditI != 0 && f.DebitI != 0:
		return nil
	default:
		return errAmountOption
	}
}
