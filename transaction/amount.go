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
	errAmountZero      = errors.New("expect non-zero amount")
	errCreditDebit     = errors.New("expect a string and an empty string in credit and debit")
	errDecimal         = errors.New("expect decimal number string")
	errPositiveDecimal = errors.New("expect positive decimal number as credit or debit")
)

/*
IsDecimal reports whether the string represents a decimal number with the following syntax:

	decimal = sign ( integer | fraction ) .
	sign = [ "-" | "+" ] .
	integer = digits .
	fraction = ( [ digits ] "." digits | digits "." ) .
	digits = digit { digit } .
	digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" .
*/
func isDecimal(s string) bool {
	var hasDigit, hasDot bool

	for i, r := range s {
		switch {
		case unicode.IsDigit(r):
			hasDigit = true
		case i == 0 && (r == '-' || r == '+'):
			// The decimal is signed.
		case !hasDot && r == '.':
			hasDot = true
		default:
			return false
		}
	}

	return hasDigit
}

/*
ParseAmount returns the value of a transaction as string and floating-point values.
The value is parsed from the amount, credit or debit fields; it cannot be zero.
If parseAmount fails to parse a non-zero value, it returns the first error.
*/
func parseAmount(fields []string, f CSVRecordFormat) (vText string, v float64, err error) {
	var negative bool

	a, c, d := fields[f.AmountI], fields[f.CreditI], fields[f.DebitI]

	switch {
	case a != "":
		vText, v, err = parseDecimal(a)
	case c != "" && d == "":
		vText, v, err = parsePositiveDecimal(c)
	case d != "" && c == "":
		vText, v, err = parsePositiveDecimal(d)
		negative = true
	default:
		return "", 0, fmt.Errorf("%w not %q and %q", errCreditDebit, c, d)
	}

	if err != nil {
		return "", 0, err // This error will be wrapped by ParseCSV.
	}

	if v == 0 {
		return "", 0, errAmountZero
	}

	if negative {
		vText = "-" + vText
		v *= -1
	}

	return vText, v, nil
}

/*
ParseDecimal returns the decimal number parsed from the string in string and floating-point representations.
If it fails to parse a number, parseDecimal returns the first error.
*/
func parseDecimal(s string) (nText string, n float64, err error) {
	if !isDecimal(s) {
		return "", 0, fmt.Errorf("%w not %q", errDecimal, s)
	}

	n, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return "", 0, err // This error will be wrapped by ParseCSV.
	}

	return s, n, nil
}

/*
ParsePositiveDecimal returns the positive number parsed from the string as string and floating-point values.
If it fails to parse a positive number, parsePositiveDecimal returns the first error.
*/
func parsePositiveDecimal(s string) (nText string, n float64, err error) {
	nText, n, err = parseDecimal(s)

	switch {
	case err != nil:
		return "", 0, err
	case n <= 0:
		return "", 0, fmt.Errorf("%w not %v", errPositiveDecimal, n)
	default:
		return nText, n, nil
	}
}
