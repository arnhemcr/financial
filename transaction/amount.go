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
	errAmountSyntax   = errors.New("parseDecimal: string must be integer or decimal")
	errAmountZero     = errors.New("parseAmount: amount cannot be zero")
	errCreditDebit    = errors.New("parseAmount: credit and debit cannot both be empty string or both non-empty string")
	errPositiveNumber = errors.New("parsePositiveDecimal: number must be positive")
)

/*
ParseAmount parses the value of a transaction from either the amount, credit or debit fields.
The value cannot be zero.
If it fails to parse a non-zero value, parseAmount returns the first error.
*/
func parseAmount(fields []string, crf CSVRecordFormat) (float64, error) {
	a, c, d := fields[crf.AmountI], fields[crf.CreditI], fields[crf.DebitI]

	var (
		v   float64
		err error
	)

	switch {
	case a != "":
		v, err = parseDecimal(a)
	case c != "" && d == "":
		v, err = parsePositiveDecimal(c)
	case d != "" && c == "":
		v, err = parsePositiveDecimal(d)

		v *= -1
	default:
		return 0, errCreditDebit
	}

	switch {
	case err != nil:
		return 0, err
	case v == 0:
		return 0, errAmountZero
	default:
		return v, nil
	}
}

/*
ParseDecimal returns the floating-point number parsed from the string.
If the string does not have the following syntax or it fails to parse a number, parseDecimal returns the error.

 number = [ "-" | "+" ] ( integer_decimal | decimal )
 integer_decimal = decimal_digits [ "." [ decimal_digits ] ]
 decimal = "." decimal_digits

*/
func parseDecimal(s string) (float64, error) {
	var postPoint bool

	for i, r := range s {
		switch {
		case i == 0 && (r == '-' || r == '+'):
		case !postPoint && r == '.':
			postPoint = true
		case unicode.IsDigit(r):
		default:
			return 0, errAmountSyntax
		}
	}

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parseDecimal: %w", err)
	}

	return n, nil
}

/*
ParsePositiveDecimal returns the positive floating-point number parsed from the string.
If it fails to parse a positive number, parsePositiveDecimal returns the first error.
*/
func parsePositiveDecimal(s string) (float64, error) {
	n, err := parseDecimal(s)

	switch {
	case err != nil:
		return 0, err
	case n <= 0:
		return 0, errPositiveNumber
	default:
		return n, nil
	}
}

// StringAmount returns the floating-point number as a string.
func stringAmount(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}
