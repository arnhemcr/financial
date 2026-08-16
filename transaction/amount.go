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
	"strconv"
	"unicode"
)

var (
	errAmountZero     = errors.New("amount cannot be zero")
	errPositiveNumber = errors.New("credit and debit cannot be negative or zero")
)

/*
ParseAmount returns the value of a transaction parsed from the amount, credit or debit fields.
The value cannot be zero.
If it fails to parse a valid value, parseAmount returns the first error.
*/
func parseAmount(fields []string, f CSVRecordFormat) (v float64, err error) {
	a, c, d := fields[f.AmountI], fields[f.CreditI], fields[f.DebitI]

	switch {
	case a != "":
		v, err = parseDecimal(a)
	case c != "" && d == "":
		v, err = parsePositiveDecimal(c)
	case d != "" && c == "":
		v, err = parsePositiveDecimal(d)

		v *= -1
	default:
		return 0, fmt.Errorf("credit and debit, %q and %q, cannot both be empty or non-empty strings", c, d)
	}

	switch {
	case err != nil:
		return 0, err // This error will be wrapped by ParseCSV.
	case v == 0:
		return 0, errAmountZero
	default:
		return v, nil
	}
}

/*
ParseDecimal returns the floating-point number parsed from the string.
If the string does not have the following syntax or it fails to parse as a number, parseDecimal returns the first error.

	number = [ "-" | "+" ] ( integer_decimal | decimal )
	integer_decimal = decimal_digits [ "." [ decimal_digits ] ]
	decimal = "." decimal_digits
*/
func parseDecimal(s string) (n float64, err error) {
	var postPoint bool

	for i, r := range s {
		switch {
		case i == 0 && (r == '-' || r == '+'):
		case unicode.IsDigit(r):
		case !postPoint && r == '.':
			postPoint = true
		default:
			return n, fmt.Errorf("cannot parse %q as amount", s)
		}
	}

	n, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}

/*
ParsePositiveDecimal returns the positive floating-point number parsed from the string.
If it fails to parse a positive number, parsePositiveDecimal returns the first error.
*/
func parsePositiveDecimal(s string) (n float64, err error) {
	n, err = parseDecimal(s)

	switch {
	case err != nil:
		return n, err
	case n <= 0:
		return n, errPositiveNumber
	default:
		return n, nil
	}
}

// StringAmount returns the floating-point number as a string.
func stringAmount(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}
