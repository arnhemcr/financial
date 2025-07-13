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
	"fmt"
	"os"
	"time"
)

const ModuleCSV = "modcsv" // The name of this module's CSV record format.

// A CSVRecordFormat defines the format of CSV records representing financial transactions.
type CSVRecordFormat struct {
	NFields uint8 // The number of fields in a record.

	// The indexes of fields in the record.
	// Some fields are required; the others are optional.
	// The index for a required field is between 1 and NFields inclusive.
	// If a field is not contained in a record, its index is zero.
	//
	// Either amount, or both credit and debit are required.
	AmountI, CreditI, DebitI uint8
	CurrencyI                uint8
	CodeI                    uint8
	DateI                    uint8 // This field is required.
	MemoI                    uint8 // This field is required.
	OtherAccountI            uint8
	ThisAccountI             uint8

	// The layout of the date field in the records e.g. "02/01/2006".
	DateLayout string
}

// GetModuleCSVFormat returns this module's CSV record format.
func GetModuleCSVFormat() CSVRecordFormat {
	return CSVRecordFormat{
		NFields:       7,
		DateI:         1,
		ThisAccountI:  2,
		OtherAccountI: 3,
		CodeI:         4,
		MemoI:         5,
		AmountI:       6,
		CurrencyI:     7,
		DateLayout:    "2006-01-02",
	}
}

/*
ReadCSVFormat returns a valid CSV record format read from the named file.
If it fails to read or validate the format, ReadCSVFormat returns the first error.
*/
func ReadCSVFormat(fileName string) (CSVRecordFormat, error) {
	var crf CSVRecordFormat

	bs, err := os.ReadFile(fileName)
	if err != nil {
		return crf, fmt.Errorf("cannot read CSV record format from file: %w", err)
	}

	err = xml.Unmarshal(bs, &crf)
	if err != nil {
		return crf, fmt.Errorf("cannot parse CSV record format from file: %w", err)
	}

	err = crf.Validate()
	if err != nil {
		return crf, err
	}

	return crf, nil
}

/*
Validate returns nil if this CSV record format is valid.
If not, Validate returns the first error.
*/
func (crf CSVRecordFormat) Validate() error {
	if crf.NFields < minNFields || maxNFields < crf.NFields {
		return errNFieldsRange
	}

	err := crf.validateIndexes()
	if err != nil {
		return err
	}

	err = crf.validateOptions()
	if err != nil {
		return err
	}

	d, _ := time.Parse(crf.DateLayout, crf.DateLayout)
	if d.Format(time.DateOnly) != time.DateOnly {
		return errDateLayout
	}

	return nil
}

const (
	// The inclusive limits for the number of fields in a CSV record.
	minNFields = 3 // date, memo and amount
	maxNFields = 20
)

var (
	errAmountOption = errors.New("amount field index, " +
		"or credit and debit indexes in CSV record format cannot both be zero")
	errDateI        = errors.New("date field index in CSV record format cannot be zero")
	errDateLayout   = errors.New("date layout in CSV record format must be Go style e.g. \"02/01/2006\"")
	errIndexUnique  = errors.New("field indexes in CSV record format cannot share a non-zero value")
	errIndexRange   = errors.New("field index in CSV record format is out of range")
	errMemoI        = errors.New("memo field index in CSV record format cannot be zero")
	errNFieldsRange = errors.New("number of fields in CSV record format is out of range")
)

/*
ValidateIndexes returns nil if the field indexes in this CSV record format are valid.
It assumes the number of fields in the format is in range.
Indexes must be <= nFields.
Each non-zero index must be unique.
Required indexes must be non-zero.
If not, validateIndexes returns the first error.
*/
func (crf CSVRecordFormat) validateIndexes() error {
	is := [...]uint8{crf.AmountI, crf.CodeI, crf.CreditI, crf.DateI, crf.DebitI,
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

	switch {
	case crf.DateI == 0:
		return errDateI
	case crf.MemoI == 0:
		return errMemoI
	default:
		return nil
	}
}

/*
ValidateOptions returns nil if the combination of optional field indexes
in this CSV record format is valid.
If not, validateOptions returns the error.
*/
func (crf CSVRecordFormat) validateOptions() error {
	switch {
	case crf.AmountI != 0:
		return nil
	case crf.CreditI != 0 && crf.DebitI != 0:
		return nil
	default:
		return errAmountOption
	}
}
