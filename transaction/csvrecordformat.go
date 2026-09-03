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

	// The Go-style date layout in the record which defaults to [time.DateOnly] or "2006-01-02" for YYYY-MM-DD.
	DateLayout string
}

/*
NewCSVRecordFormat returns the CSV record format read from the named XML file.
Fields in the format default to zero except DateLayout which defaults to
Go-style layout time.DateOnly or "2006-01-02" for YYYY-MM-DD.
If it fails to read a valid format, NewCSVRecordFormat returns the first error.
*/
func NewCSVRecordFormat(name string) (f CSVRecordFormat, err error) {
	bs, err := os.ReadFile(name)
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
		return fmt.Errorf("Validate: %w %v to %v inclusive not %v", errNFieldsRange, minNFields, maxNFields, n)
	}

	err := f.validateIndexes()
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	err = f.validateOptions()
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	err = ValidateDateLayout(f.DateLayout)
	if err != nil {
		return fmt.Errorf("Validate: %w", err)
	}

	return nil
}

const (
	// The inclusive limits for the number of fields in a CSV record.
	minNFields = 3 // The minimal set of fields in a CSV record is amount, date and memo.
	maxNFields = 20
)

var (
	errAmountOption = errors.New("expect either amount field index, or credit and debit indexes to be non-zero")
	errDateI        = errors.New("expect date field index to be non-zero")
	errIndexUnique  = errors.New("expect non-zero field indexes to be unique")
	errIndexRange   = errors.New("expect field index in range")
	errMemoI        = errors.New("expect memo field index to be non-zero")
	errNFieldsRange = errors.New("expect number of fields in range")
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
			return fmt.Errorf("%w %v to %v inclusive not %v", errIndexRange, 1, f.NFields, i)
		case used[i]:
			return errIndexUnique
		default:
			used[i] = true
		}
	}

	if f.DateI == 0 {
		return errDateI
	}

	if f.MemoI == 0 {
		return errMemoI
	}

	return nil
}

/*
ValidateOptions returns nil if the combination of optional field indexes in this CSV record format is valid.
If not, validateOptions returns the error.
*/
func (f CSVRecordFormat) validateOptions() error {
	if f.AmountI == 0 && (f.CreditI == 0 || f.DebitI == 0) {
		return fmt.Errorf("%w not %v, or %v and %v", errAmountOption, f.AmountI, f.CreditI, f.DebitI)
	}

	return nil
}
