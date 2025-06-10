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
	"fmt"
	"strconv"
)

var (
	errAmount      = errors.New("t.parsecsv: amount cannot be zero")
	errCreditDebit = errors.New("t.parsecsv: credit and debit cannot both be empty string or both non-empty string")
	errNotPositive = errors.New("t.parsecsv: number must be positive")
)

/*
ParseAmount parses the value of a transaction
from either the amount, credit or debit fields.
The value cannot be zero.
If it fails to parse a non-zero value, parseAmount returns the first error.
*/
func parseAmount(fields []string, crf CSVRecordFormat) (float64, error) {
	a, c, d := fields[crf.AmountI], fields[crf.CreditI], fields[crf.DebitI]

	var v float64

	var err error

	switch {
	case a != "":
		v, err = parseFloat(a)
	case c != "" && d == "":
		v, err = parsePositiveFloat(c)
	case d != "" && c == "":
		v, err = parsePositiveFloat(d)

		v *= -1
	default:
		return 0, errCreditDebit
	}

	switch {
	case err != nil:
		return 0, err
	case v == 0:
		return 0, errAmount
	default:
		return v, nil
	}
}

/*
ParseFloat returns the floating-point number parsed from the string.
If it fails to parse a number, parseFloat returns the error.
*/
func parseFloat(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("t.parsecsv: %w", err)
	}

	return n, nil
}

/*
ParsePositiveFloat returns the positive floating-point number parsed from the string.
If it fails to parse a positive number, parsePositiveFloat returns the first error.
*/
func parsePositiveFloat(s string) (float64, error) {
	n, err := parseFloat(s)
	switch {
	case err != nil:
		return 0, err
	case n <= 0:
		return 0, errNotPositive
	default:
		return n, nil
	}
}

// StringAmount returns the floating-point number as a string.
func stringAmount(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}
