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
	errAmountMismatch  = errors.New("expect amount text and number to match")
	errCreditDebit     = errors.New("expect string in credit and empty string in debit or visa versa")
	errDecimal         = errors.New("expect decimal number string")
	errDecimalZero     = errors.New("expect non-zero decimal number")
	errPositiveDecimal = errors.New("expect positive decimal number as credit or debit")
)

/*
IsDecimal reports whether the string represents a decimal number with the following syntax:

	decimal = [ sign ] ( integer | fraction ) .
	sign = "-" | "+" .
	integer = digits .
	fraction = ( [ digits ] "." digits | digits "." ) .
	digits = digit { digit } .
	digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" .
*/
func isDecimal(d string) bool {
	var hasDigit, hasDot bool

	for i, r := range d {
		switch {
		case unicode.IsDigit(r):
			hasDigit = true
		case i == 0 && (r == '-' || r == '+'):
			// This decimal is signed.
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
		vText, v, err = parseNonZeroDecimal(a)
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

	if negative {
		vText = "-" + vText
		v *= -1
	}

	return vText, v, nil
}

/*
ParseDecimal returns the non-zero decimal number parsed from the string as both string and floating-point values.
If it fails to parse a non-zero number, parseDecimal returns the first error.
*/
func parseNonZeroDecimal(d string) (nText string, n float64, err error) {
	if !isDecimal(d) {
		return "", 0, fmt.Errorf("%w not %q", errDecimal, d)
	}

	n, err = strconv.ParseFloat(d, 64)

	switch {
	case err != nil:
		return "", 0, err // This error will be wrapped by ParseCSV.
	case n == 0:
		return "", 0, fmt.Errorf("%w not %v", errDecimalZero, d)
	default:
		return d, n, nil
	}
}

/*
ParsePositiveDecimal returns the positive number parsed from the string as string and floating-point values.
If it fails to parse a positive number, parsePositiveDecimal returns the first error.
*/
func parsePositiveDecimal(d string) (nText string, n float64, err error) {
	nText, n, err = parseNonZeroDecimal(d)

	switch {
	case err != nil:
		return "", 0, err
	case n < 0:
		return "", 0, fmt.Errorf("%w not %v", errPositiveDecimal, n)
	default:
		return nText, n, nil
	}
}
