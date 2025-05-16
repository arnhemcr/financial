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
	"errors"
	"slices"
	"strings"
	"time"
)

// A CSVformat defines the format of CSV records representing financial transactions.
type CSVformat struct {
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
	DateI           uint8
	MemoI           uint8
	OtherAccountI   uint8
	ThisAccountI    uint8

	// The Go-style format of the date field in the records e.g. "02/01/2006".
	DateFormat string
}

const NIndexes = 7 // The number of field indexes in a CSV format.

// GetAfFormat returns the arnhemcr/financial CSV format.
func GetAfFormat() CSVformat {
	return CSVformat{
		NFields:       5,
		DateI:         1,
		ThisAccountI:  2,
		OtherAccountI: 3,
		MemoI:         4,
		AmountI:       5,
		DateFormat:    "2006-01-02",
	}
}

// IsValid reports whether this CSV format is valid.
func (format CSVformat) IsValid() bool {
	return format.Validate() == nil
}

/*
ParseCSV parses this transaction from fields, according to the CSV format, and returns nil.
It assumes the format is valid.
If ParseCSV fails to parse the transaction, it returns the first error.
*/
func (trn *Transaction) ParseCSV(fields []string, format CSVformat) error {
	if len(fields) != int(format.NFields) {
		return errNFields
	}

	/*
		Prepend fields with an empty string,
		so a field whose index is zero has a value of empty string.
	*/
	flds := slices.Insert(fields, 0, "")

	var err error

	trn.Date, err = parseDate(flds[format.DateI], format.DateFormat)
	if err != nil {
		return err
	}

	trn.Amount, err = parseValue(flds, format)
	if err != nil {
		return err
	} else if trn.Amount == 0.00 {
		return errAmount
	}

	trn.Memo = flds[format.MemoI]
	if trn.Memo == "" {
		return errMemo
	}

	trn.OtherAccount = flds[format.OtherAccountI]
	if trn.OtherAccount == "" {
		trn.OtherAccount = DefaultOtherAccount
	}

	switch {
	case format.ThisAccount != "":
		trn.ThisAccount = format.ThisAccount
	case flds[format.ThisAccountI] != "":
		trn.ThisAccount = flds[format.ThisAccountI]
	default:
		return errThisAccount
	}

	return nil
}

// StringCSV returns this transaction as an arnhemcr/financial CSV record.
func (trn Transaction) StringCSV() string {
	amt := formatAmount(trn.Amount)
	flds := []string{trn.Date, trn.ThisAccount, trn.OtherAccount, trn.Memo, amt}

	return strings.Join(flds, ",") + "\n"
}

/*
Validate returns nil if this CSV format is valid.
If not, validate returns the first error.
*/
func (format CSVformat) Validate() error {
	val, _ := time.Parse(format.DateFormat, format.DateFormat)
	if val.Format(time.DateOnly) != time.DateOnly {
		return errDateFormat
	}

	if format.NFields < minNFields || maxNFields < format.NFields {
		return errNFieldsRange
	}

	err := format.areIndexesValid()
	if err != nil {
		return err
	}

	if format.DateI == 0 {
		return errDateI
	}

	if format.MemoI == 0 {
		return errMemoI
	}

	err = format.areOptionsValid()
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
	errDateFormat   = errors.New("date format in CSV record must be Go style e.g. \"02/01/2006\"")
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
func (format CSVformat) areIndexesValid() error {
	inxs := [NIndexes]uint8{format.AmountI, format.CreditI, format.DateI, format.DebitI,
		format.MemoI, format.OtherAccountI, format.ThisAccountI}

	var inUse [maxNFields + 1]bool

	for _, val := range inxs {
		switch {
		case format.NFields < val:
			return errIndexRange
		case val == 0:
			// CSV records do not contain this field
		case inUse[val]:
			return errIndexUnique
		default:
			inUse[val] = true
		}
	}

	return nil
}

/*
AreOptionsValid returns nil if the combination of options is valid.
If not, areOptionsValid returns the first error.
*/
func (format CSVformat) areOptionsValid() error {
	if format.ThisAccount == "" && format.ThisAccountI == 0 {
		return errThisAcctOpt
	}

	if (format.AmountI == 0) && (format.CreditI == 0 || format.DebitI == 0) {
		return errAmountOpt
	}

	return nil
}

/*
ParseValue parses the value of a transaction from either the amount, credit or debit fields.
If it fails to parse the value, parseValue returns the error.
*/
func parseValue(fields []string, format CSVformat) (float64, error) {
	amt, crt, dbt := fields[format.AmountI], fields[format.CreditI], fields[format.DebitI]

	switch {
	case amt != "":
		return parseAmount(amt)
	case crt != "" && dbt == "":
		return parsePositiveAmount(crt)
	case dbt != "" && crt == "":
		val, err := parsePositiveAmount(dbt)

		return val * -1.00, err
	default:
		return 0.00, errCreditDebit
	}
}
