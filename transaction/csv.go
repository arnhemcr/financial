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
	"slices"
	"strings"
	"time"
)

const CSV = "csv" // The name of this package's CSV format.

// A CSVFormat defines the format of CSV records representing financial transactions.
type CSVFormat struct {
	/*
		The name of the account these records belong to.
		If it is an empty string, the records must contain a this account field.
	*/
	ThisAccount string

	NFields uint8 // The number of fields in each record.

	/*
		The indexes of fields in the records.
		Index values can be:

		 * 0: this field is not contained in these records
		 * 1..NFields
	*/
	AmountI uint8
	// If the amount index is zero, both the credit and debit indexes must be non-zero.
	CreditI, DebitI uint8
	CurrencyI       uint8
	DateI           uint8
	MemoI           uint8
	OtherAccountI   uint8 // The other account index can be zero.
	// If this account is an empty string, this account index must be non-zero.
	ThisAccountI uint8

	// The layout of the date field in the records e.g. "02/01/2006".
	DateLayout string
}

/*
GetFormat returns the first CSV format from the named file.
If it fails to get a format, getFormat returns the first error.
*/
func GetFormat(fileName string) (CSVFormat, error) {
	var cf CSVFormat

	bs, err := os.ReadFile(fileName)
	if err != nil {
		return cf, fmt.Errorf("getformat: %w", err)
	}

	err = xml.Unmarshal(bs, &cf)
	if err != nil {
		return cf, fmt.Errorf("getformat: %w", err)
	}

	return cf, nil
}

// GetPkgFormat returns this package's CSV format.
func GetPkgFormat() CSVFormat {
	return CSVFormat{
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

// IsValid reports whether this CSV format is valid.
func (cf CSVFormat) IsValid() bool {
	return cf.Validate() == nil
}

/*
ParseCSV parses this transaction from fields, according to the CSV format, and returns nil.
It assumes the format is valid.
If ParseCSV fails to parse the transaction, it returns the first error.
*/
func (t *Transaction) ParseCSV(fields []string, cf CSVFormat) error {
	if len(fields) != int(cf.NFields) {
		return errNFields
	}

	/*
		Prepend fields with an empty string,
		so a field whose index is zero has a value of empty string.
	*/
	fs := slices.Insert(fields, 0, "")

	var err error

	t.Date, err = parseDate(fs[cf.DateI], cf.DateLayout)
	if err != nil {
		return err
	}

	t.Amount, err = parseValue(fs, cf)
	if err != nil {
		return err
	} else if t.Amount == 0.00 {
		return errAmount
	}

	t.Memo = fs[cf.MemoI]
	if t.Memo == "" {
		return errMemo
	}

	t.OtherAccount = fs[cf.OtherAccountI]
	if t.OtherAccount == "" {
		t.OtherAccount = DefaultOtherAccount
	}

	switch {
	case cf.ThisAccount != "":
		t.ThisAccount = cf.ThisAccount
	case fs[cf.ThisAccountI] != "":
		t.ThisAccount = fs[cf.ThisAccountI]
	default:
		return errThisAccount
	}

	return nil
}

// StringCSV returns this transaction as a CSV record in this package's format.
func (t Transaction) StringCSV() string {
	a := formatAmount(t.Amount)
	fs := []string{t.Date, t.ThisAccount, t.OtherAccount, t.Memo, a, t.Currency}

	return strings.Join(fs, ",") + "\n"
}

/*
Validate returns nil if this CSV format is valid.
If not, validate returns the first error.
*/
func (cf CSVFormat) Validate() error {
	d, _ := time.Parse(cf.DateLayout, cf.DateLayout)
	if d.Format(time.DateOnly) != time.DateOnly {
		return errDateLayout
	}

	if cf.NFields < minNFields || maxNFields < cf.NFields {
		return errNFieldsRange
	}

	err := cf.areIndexesValid()
	if err != nil {
		return err
	}

	if cf.DateI == 0 {
		return errDateI
	}

	if cf.MemoI == 0 {
		return errMemoI
	}

	err = cf.areOptionsValid()
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
	errNFieldsRange = errors.New("number of fields in CSV record is out of range")
	errThisAcctOpt  = errors.New("this account and this account index " +
		"cannot be empty string and zero respectively")
)

/*
AreIndexesValid returns nil if the field indexes in this CSV format are valid.
It assumes the number of fields in the format NFields is in range.
All indexes must be <= nFields, and all non-zero indexes must be unique.
If not, areIndexesValid returns the first error.
*/
func (cf CSVFormat) areIndexesValid() error {
	is := [...]uint8{cf.AmountI, cf.CreditI, cf.DateI, cf.DebitI,
		cf.MemoI, cf.OtherAccountI, cf.ThisAccountI}

	var inUse [maxNFields + 1]bool

	for _, i := range is {
		switch {
		case cf.NFields < i:
			return errIndexRange
		case i == 0:
			// CSV records do not contain this field
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
If not, areOptionsValid returns the first error.
*/
func (cf CSVFormat) areOptionsValid() error {
	if cf.ThisAccount == "" && cf.ThisAccountI == 0 {
		return errThisAcctOpt
	}

	if (cf.AmountI == 0) && (cf.CreditI == 0 || cf.DebitI == 0) {
		return errAmountOpt
	}

	return nil
}

/*
ParseValue parses the value of a transaction from either the amount, credit or debit fields.
If it fails to parse the value, parseValue returns the error.
*/
func parseValue(fields []string, cf CSVFormat) (float64, error) {
	a, c, d := fields[cf.AmountI], fields[cf.CreditI], fields[cf.DebitI]

	switch {
	case a != "":
		return parseAmount(a)
	case c != "" && d == "":
		return parsePositiveAmount(c)
	case d != "" && c == "":
		n, err := parsePositiveAmount(d)

		return n * -1.00, err
	default:
		return 0.00, errCreditDebit
	}
}
